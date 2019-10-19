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
	"time"

	"regexp"

	"github.com/majormjr/rcon"
)

type FactorioServer struct {
	Cmd            *exec.Cmd              `json:"-"`
	Savefile       string                 `json:"savefile"`
	Latency        int                    `json:"latency"`
	BindIP         string                 `json:"bindip"`
	Port           int                    `json:"port"`
	Running        bool                   `json:"running"`
	Version        Version                `json:"fac_version"`
	BaseModVersion string                 `json:"base_mod_version"`
	StdOut         io.ReadCloser          `json:"-"`
	StdErr         io.ReadCloser          `json:"-"`
	StdIn          io.WriteCloser         `json:"-"`
	Settings       map[string]interface{} `json:"-"`
	Rcon           *rcon.RemoteConsole    `json:"-"`
	LogChan        chan []string          `json:"-"`
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

		err = example.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close example server settings: %s", err)
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

	out := []byte{}
	//Load factorio version
	if config.glibcCustom == "true" {
		out, err = exec.Command(config.glibcLocation, "--library-path", config.glibcLibLoc, config.FactorioBinary, "--version").Output()
	} else {
		out, err = exec.Command(config.FactorioBinary, "--version").Output()
	}

	if err != nil {
		log.Printf("error on loading factorio version: %s", err)
		return
	}

	reg := regexp.MustCompile("Version.*?((\\d+\\.)?(\\d+\\.)?(\\*|\\d+)+)")
	found := reg.FindStringSubmatch(string(out))
	err = f.Version.UnmarshalText([]byte(found[1]))
	if err != nil {
		log.Printf("could not parse version: %v", err)
		return
	}

	//Load baseMod version
	baseModInfoFile := filepath.Join(config.FactorioDir, "data", "base", "info.json")
	bmifBa, err := ioutil.ReadFile(baseModInfoFile)
	if err != nil {
		log.Printf("couldn't open baseMods info.json: %s", err)
		return
	}
	var modInfo ModInfo
	err = json.Unmarshal(bmifBa, &modInfo)
	if err != nil {
		log.Printf("error unmarshalling baseMods info.json to a modInfo: %s", err)
		return
	}

	f.BaseModVersion = modInfo.Version

	// load admins from additional file
	if(f.Version.Greater(Version{0,17,0})) {
		if _, err := os.Stat(filepath.Join(config.FactorioConfigDir, config.FactorioAdminFile)); os.IsNotExist(err) {
			//save empty admins-file
			ioutil.WriteFile(filepath.Join(config.FactorioConfigDir, config.FactorioAdminFile), []byte("[]"), 0664)
		} else {
			data, err := ioutil.ReadFile(filepath.Join(config.FactorioConfigDir, config.FactorioAdminFile))
			if err != nil {
				log.Printf("Error loading FactorioAdminFile: %s", err)
				return f, err
			}

			var jsonData interface{}
			err = json.Unmarshal(data, &jsonData)
			if err != nil {
				log.Printf("Error unmarshalling FactorioAdminFile: %s", err)
				return f, err
			}

			f.Settings["admins"] = jsonData
		}
	}

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

	args := []string{}

	//The factorio server refenences its executable-path, since we execute the ld.so file and pass the factorio binary as a parameter
	//the game would use the path to the ld.so file as it's executable path and crash, to prevent this the parameter "--executable-path" is added
	if config.glibcCustom == "true" {
		log.Println("Custom glibc selected, glibc.so location:", config.glibcLocation, " lib location:", config.glibcLibLoc)
		args = append(args, "--library-path", config.glibcLibLoc, config.FactorioBinary, "--executable-path", config.FactorioBinary)
	}

	args = append(args,
		"--bind", (f.BindIP),
		"--port", strconv.Itoa(f.Port),
		"--server-settings", filepath.Join(config.FactorioConfigDir, config.SettingsFile),
		"--rcon-port", strconv.Itoa(config.FactorioRconPort),
		"--rcon-password", config.FactorioRconPass)

	if(f.Version.Greater(Version{0,17,0})) {
		args = append(args, "--server-adminlist", filepath.Join(config.FactorioConfigDir, config.FactorioAdminFile))
	}

	if f.Savefile == "Load Latest" {
		args = append(args, "--start-server-load-latest")
	} else {
		args = append(args, "--start-server", filepath.Join(config.FactorioSavesDir, f.Savefile))
	}

	if config.glibcCustom == "true" {
		log.Println("Starting server with command: ", config.glibcLocation, args)
		f.Cmd = exec.Command(config.glibcLocation, args...)
	} else {
		log.Println("Starting server with command: ", config.FactorioBinary, args)
		f.Cmd = exec.Command(config.FactorioBinary, args...)
	}

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
		if err := f.writeLog(stdScanner.Text()); err != nil {
			log.Printf("Error: %s", err)
		}

		line := strings.Fields(stdScanner.Text())
		// Ensure logline slice is in bounds
		if len(line) > 1 {
			// Check if Factorio Server reports any errors if so handle it
			if line[1] == "Error" {
				err := f.checkLogError(line)
				if err != nil {
					log.Printf("Error checking Factorio Server Error: %s", err)
				}
			}
			// If rcon port opens indicated in log connect to rcon
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

func (f *FactorioServer) writeLog(logline string) error {
	logfileName := filepath.Join(config.FactorioDir, "factorio-server-console.log")
	file, err := os.OpenFile(logfileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Cannot open logfile for appending Factorio Server output: %s", err)
		return err
	}
	defer file.Close()

	logline = logline + "\n"

	if _, err = file.WriteString(logline); err != nil {
		log.Printf("Error appending to factorio-server-console.log: %s", err)
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
	if runtime.GOOS == "windows" {

		// Disable our own handling of CTRL+C, so we don't close when we send it to the console.
		setCtrlHandlingIsDisabledForThisProcess(true)

		// Send CTRL+C to all processes attached to the console (ourself, and the factorio server instance)
		sendCtrlCToPid(0)
		log.Println("Sent SIGINT to Factorio process. Factorio shutting down...")

		// Somehow, the Factorio devs managed to code the game to react appropriately to CTRL+C, including
		// saving the game, but not actually exit. So, we still have to manually kill the process, and
		// for extra fun, there's no way to know when the server save has actually completed (unless we want
		// to inject filesystem logic into what should be a process-level Stop() routine), so our best option
		// is to just wait an arbitrary amount of time and hope that the save is successful in that time.
		time.Sleep(2 * time.Second)
		f.Cmd.Process.Signal(os.Kill)

		// Re-enable handling of CTRL+C after we're sure that the factrio server is shut down.
		setCtrlHandlingIsDisabledForThisProcess(false)

		f.Running = false
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

func (f *FactorioServer) Kill() error {
	if runtime.GOOS == "windows" {

		err := f.Cmd.Process.Signal(os.Kill)
		if err != nil {
			if err.Error() == "os: process already finished" {
				f.Running = false
				return err
			}
			log.Printf("Error sending SIGKILL to Factorio process: %s", err)
			return err
		}
		f.Running = false
		log.Println("Sent SIGKILL to Factorio process. Factorio forced to exit.")

		return nil
	}

	err := f.Cmd.Process.Signal(os.Kill)
	if err != nil {
		if err.Error() == "os: process already finished" {
			f.Running = false
			return err
		}
		log.Printf("Error sending SIGKILL to Factorio process: %s", err)
		return err
	}
	f.Running = false
	log.Printf("Sent SIGKILL to Factorio process. Factorio forced to exit.")

	err = f.Rcon.Close()
	if err != nil {
		log.Printf("Error close rcon connection: %s", err)
	}

	return nil
}
