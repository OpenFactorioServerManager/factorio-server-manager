package main

import (
	"io/ioutil"
	"log"
)

// Lists save files in factorio/saves
func listSaves() []string {
	saveDir := config.FactorioDir + "/saves"
	result := []string{}

	files, err := ioutil.ReadDir(saveDir)
	if err != nil {
		log.Printf("Error listing save directory: %s", err)
		return result
	}

	for _, f := range files {
		result = append(result, f.Name())
	}

	return result
}
