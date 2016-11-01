package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
)

type FactorioServer struct {
	Cmd      *exec.Cmd              `json:"-"`
	Savefile string                 `json:"-"`
	Latency  int                    `json:"latency"`
	Port     int                    `json:"port"`
	Running  bool                   `json:"running"`
	StdOut   io.ReadCloser          `json:"-"`
	StdErr   io.ReadCloser          `json:"-"`
	StdIn    io.WriteCloser         `json:"-"`
	Settings map[string]interface{} `json:"-"`
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
		"--start-server", filepath.Join(config.FactorioSavesDir, f.Savefile),
		"--port", strconv.Itoa(f.Port),
		"--server-settings", filepath.Join(config.FactorioDir, "server-settings.json"),
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

	go io.Copy(os.Stdout, f.StdOut)
	go io.Copy(os.Stderr, f.StdErr)

	err = f.Cmd.Start()
	if err != nil {
		log.Printf("Error starting server process: %s", err)
		return err
	}

	f.Running = true

	err = f.Cmd.Wait()
	if err != nil {
		log.Printf("Command exited with error: %s", err)
		return err
	}

	return nil
}

func (f *FactorioServer) Stop() error {
	// TODO: Find an alternative to os.Kill on Windows. os.Interupt
	// is not implemented. Maps will not be saved.
	if runtime.GOOS == "windows" {
		err := f.Cmd.Process.Signal(os.Kill)
		if err != nil {
			log.Printf("Error sending SIGKILLL to Factorio process: %s", err)
			return err
		} else {
			f.Running = false
			log.Println("Sent SIGKILL to Factorio process. Factorio forced to exit.")
		}
	} else {
		err := f.Cmd.Process.Signal(os.Interrupt)
		if err != nil {
			log.Printf("Error sending SIGINT to Factorio process: %s", err)
			return err
		} else {
			f.Running = false
			log.Printf("Sent SIGINT to Factorio process. Factorio shutting down...")
		}
	}

	return nil
}
