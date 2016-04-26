package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

type FactorioServer struct {
	Cmd              *exec.Cmd
	Savefile         string
	Latency          int  `json:"latency"`
	AutosaveInterval int  `json:"autosave_interval"`
	AutosaveSlots    int  `json:"autosave_slots"`
	Port             int  `json:"port"`
	DisallowCmd      bool `json:"disallow_cmd"`
	Running          bool
	PeerToPeer       bool `json:"peer2peer"`
	AutoPause        bool `json:"auto_pause"`
	StdOut           io.ReadCloser
	StdErr           io.ReadCloser
	StdIn            io.WriteCloser
}

func createSave(saveName string) (string, error) {
	args := []string{"--create", saveName}

	cmdOutput, err := exec.Command(config.FactorioBinary, args...).Output()
	if err != nil {
		log.Printf("Error in creating Factorio save: %s", err)
		return "", err
	}

	result := string(cmdOutput)

	return result, nil
}

func initFactorio() *FactorioServer {
	// TODO move values to config struct
	f := FactorioServer{
		Latency:          100,
		AutosaveInterval: 5,
		AutosaveSlots:    10,
	}

	return &f
}

func (f *FactorioServer) Run() error {
	var err error

	args := []string{"--start-server", f.Savefile,
		"--latency-ms", strconv.Itoa(f.Latency),
		"--autosave-interval", strconv.Itoa(f.AutosaveInterval),
		"--autosave-slots", strconv.Itoa(f.AutosaveSlots),
		"--port", strconv.Itoa(f.Port)}
	if f.DisallowCmd {
		args = append(args, "--disallow-commands")
	}
	if f.PeerToPeer {
		args = append(args, "--peer-to-peer")
	}
	if f.AutoPause {
		args = append(args, "--no-auto-pause")
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
	err := f.Cmd.Process.Signal(syscall.SIGINT)
	if err != nil {
		log.Printf("Error sending SIGINT to Factorio process: %s", err)
	}

	f.Running = false

	log.Printf("Sent SIGINT to Factorio process")

	return nil
}
