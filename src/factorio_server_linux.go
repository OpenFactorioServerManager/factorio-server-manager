// use this file only when compiling not windows (all unix systems)
// +build !windows

package main

import (
	"log"
	"os"
)

// Stubs for windows-only functions

func (f *FactorioServer) Kill() error {
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

func (f *FactorioServer) Stop() error {
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
