package main

import (
	"encoding/json"
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
	Cmd      *exec.Cmd
	Savefile string
	Latency  int  `json:"latency"`
	Port     int  `json:"port"`
	Running  bool `json:"running"`
	StdOut   io.ReadCloser
	StdErr   io.ReadCloser
	StdIn    io.WriteCloser
	Settings FactorioServerSettings
}

type FactorioServerSettings struct {
	Name                                 string   `json:"name"`
	Description                          string   `json:"description"`
	Tags                                 []string `json:"tags"`
	MaxPlayers                           int      `json:"max_players"`
	Visibility                           string   `json:"visibility"`
	Username                             string   `json:"username"`
	Password                             string   `json:"password"`
	Token                                string   `json:"token"`
	GamePassword                         string   `json:"game_password"`
	RequireUserVerification              bool     `json:"require_user_verification"`
	MaxUploadInKilobytesPerSecond        int      `json:"max_upload_in_kilobytes_per_second"`
	IgnorePlayerLimitForReturningPlayers bool     `json:"ignore_player_limit_for_returning_players"`
	AllowCommands                        string   `json:"allow_commands"`
	AutosaveInterval                     int      `json:"autosave_interval"`
	AutosaveSlots                        int      `json:"autosave_slots"`
	AfkAutoKickInterval                  int      `json:"afk_autokick_interval"`
	AutoPause                            bool     `json:"auto_pause"`
	OnlyAdminsCanPauseThegame            bool     `json:"only_admins_can_pause_the_game"`
	Admins                               []string `json:"admins"`
	AutosaveOnlyOnServer                 bool     `json:"autosave_only_on_server"`
}

func initFactorio() *FactorioServer {
	f := FactorioServer{}
	settingsFile := filepath.Join(config.FactorioConfigDir, config.SettingsFile)

	if _, err := os.Stat(settingsFile); err == nil {
		// server-settings.json file exists
		settings, err := os.Open(settingsFile)
		if err != nil {
			log.Printf("Error in reading %s: %s", settingsFile, err)
		}
		defer settings.Close()

		settingsParser := json.NewDecoder(settings)
		if err = settingsParser.Decode(&f.Settings); err != nil {
			log.Printf("Error in reading %s: %s", settingsFile, err)
		}
		log.Printf("Loaded Factorio settings from %s, settings: %+v", settingsFile, &f.Settings)

	} else {
		// default settings from server-settings.example.json
		f.Settings = FactorioServerSettings{
			Name:                                 "Factorio",
			Description:                          "Created by Factorio Server Manager",
			Tags:                                 []string{},
			MaxPlayers:                           0,
			Visibility:                           "public",
			Username:                             "",
			Password:                             "",
			Token:                                "",
			GamePassword:                         "",
			RequireUserVerification:              true,
			MaxUploadInKilobytesPerSecond:        0,
			IgnorePlayerLimitForReturningPlayers: false,
			AllowCommands:                        "admins-only",
			AutosaveInterval:                     10,
			AutosaveSlots:                        5,
			AfkAutoKickInterval:                  0,
			AutosaveOnlyOnServer:                 true,
			AutoPause:                            true,
			OnlyAdminsCanPauseThegame:            true,
			Admins: []string{},
		}
		log.Printf("Loaded Default Factorio settings settings: %+v", &f.Settings)
	}

	return &f
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
