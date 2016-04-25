package main

import (
	"log"
	"os/exec"
)

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
