package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/majormjr/rcon"
)

type FactorioServer struct {
	Cmd      *exec.Cmd              `json:"-"`
	Savefile string                 `json:"savefile"`
	Latency  int                    `json:"latency"`
	Port     int                    `json:"port"`
	Running  bool                   `json:"running"`
	StdOut   io.ReadCloser          `json:"-"`
	StdErr   io.ReadCloser          `json:"-"`
	StdIn    io.WriteCloser         `json:"-"`
	Settings map[string]interface{} `json:"-"`
	Rcon     *rcon.RemoteConsole    `json:"-"`
	LogChan  chan []string          `json:"-"`
}

func randomPort() int {
	// Returns random port to use for rcon connection
	return rand.Intn(45000-40000) + 40000
}

func initFactorio() (f *FactorioServer, err error) {
	f = new(FactorioServer)
	f.Settings = make(map[string]interface{})

	if err = os.MkdirAll(config.FactorioConfigDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %v", err)
	}

	settingsPath := filepath.Join(config.FactorioConfigDir, config.SettingsFile)
	var settings *os.File

	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		// copy example settings to supplied settings file, if not exists
		log.Printf("Server settings at %s not found, copying example server settings.\n", settingsPath)

		examplePath := filepath.Join(config.FactorioDir, "data", "server-settings.example.json")

		example, err := os.Open(examplePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open example server settings: %v", err)
		}
		defer example.Close()

		settings, err = os.Create(settingsPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create server settings file: %v", err)
		}
		defer settings.Close()

		_, err = io.Copy(settings, example)
		if err != nil {
			return nil, fmt.Errorf("failed to copy example server settings: %v", err)
		}
	} else {
		// otherwise, open file normally
		settings, err = os.Open(settingsPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open server settings file: %v", err)
		}
		defer settings.Close()
	}

	// before reading reset offset
	if _, err = settings.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("error while seeking in settings file: %v", err)
	}

	if err = json.NewDecoder(settings).Decode(&f.Settings); err != nil {
		return nil, fmt.Errorf("error reading %s: %v", settingsPath, err)
	}

	log.Printf("Loaded Factorio settings from %s\n", settingsPath)

	return
}

func (f *FactorioServer) Run() error {
	var err error

	data, err := json.MarshalIndent(f.Settings, "", "  ")
	if err != nil {
		log.Println("Failed to marshal FactorioServerSettings: ", err)
	} else {
		ioutil.WriteFile(filepath.Join(config.FactorioConfigDir, config.SettingsFile), data, 0644)
	}

	args := []string{
		"--port", strconv.Itoa(f.Port),
		"--server-settings", filepath.Join(config.FactorioConfigDir, "server-settings.json"),
		"--rcon-port", strconv.Itoa(config.FactorioRconPort),
		"--rcon-pass", config.FactorioRconPass,
	}

	if f.Savefile == "Load Latest" {
		args = append(args, "--start-server-load-latest")
	} else {
		args = append(args, "--start-server", filepath.Join(config.FactorioSavesDir, f.Savefile))
	}

	log.Println("Starting server with command: ", config.FactorioBinary, args)

	f.Cmd = exec.Command(config.FactorioBinary, args...)

	f.StdOut, err = f.Cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error opening stdout pipe: %s", err)
		return err
	}

	f.StdIn, err = f.Cmd.StdinPipe()
	if err != nil {
		log.Printf("Error opening stdin pipe: %s", err)
		return err
	}

	f.StdErr, err = f.Cmd.StderrPipe()
	if err != nil {
		log.Printf("Error opening stderr pipe: %s", err)
		return err
	}

	go f.parseRunningCommand(f.StdOut)
	go f.parseRunningCommand(f.StdErr)

	err = f.Cmd.Start()
	if err != nil {
		log.Printf("Factorio process failed to start: %s", err)
		return err
	}
	f.Running = true

	err = f.Cmd.Wait()
	if err != nil {
		log.Printf("Factorio process exited with error: %s", err)
		f.Running = false
		return err
	}

	return nil
}

func (f *FactorioServer) parseRunningCommand(std io.ReadCloser) (err error) {
	stdScanner := bufio.NewScanner(std)
	for stdScanner.Scan() {
		log.Printf("Factorio Server: %s", stdScanner.Text())
		line := strings.Fields(stdScanner.Text())
		// Check if Factorio Server reports any errors if so handle it
		if len(line) > 0 {
			if line[1] == "Error" {
				err := f.checkLogError(line)
				if err != nil {
					log.Printf("Error checking Factorio Server Error: %s", err)
				}
			}
			// If rcon port is opened connect to rcon
			rconLog := "Starting RCON interface at port " + strconv.Itoa(config.FactorioRconPort)
			// check if slice index is greater than 2 to prevent panic
			if len(line) > 2 {
				// log line for opened rcon connection
				if strings.Join(line[3:], " ") == rconLog {
					log.Printf("Rcon running on Factorio Server")
					err = connectRC()
					if err != nil {
						log.Printf("Error: %s", err)
					}
				}
			}
		}
	}
	if err := stdScanner.Err(); err != nil {
		log.Printf("Error reading std buffer: %s", err)
		return err
	}
	return nil
}

func (f *FactorioServer) checkLogError(logline []string) error {
	// TODO Handle errors generated by running Factorio Server
	log.Println(logline)

	return nil
}

func (f *FactorioServer) Stop() error {
	// TODO: Find an alternative to os.Kill on Windows. os.Interupt
	// is not implemented. Maps will not be saved.
	if runtime.GOOS == "windows" {
		err := f.Cmd.Process.Signal(os.Kill)
		if err != nil {
			if err.Error() == "os: process already finished" {
				f.Running = false
				return err
			}
			log.Printf("Error sending SIGKILLL to Factorio process: %s", err)
			return err
		}
		f.Running = false
		log.Println("Sent SIGKILL to Factorio process. Factorio forced to exit.")

		return nil
	}

	err := f.Cmd.Process.Signal(os.Interrupt)
	if err != nil {
		if err.Error() == "os: process already finished" {
			f.Running = false
			return err
		}
		log.Printf("Error sending SIGINT to Factorio process: %s", err)
		return err
	}
	f.Running = false
	log.Printf("Sent SIGINT to Factorio process. Factorio shutting down...")

	err = f.Rcon.Close()
	if err != nil {
		log.Printf("Error close rcon connection: %s", err)
	}

	return nil
}
