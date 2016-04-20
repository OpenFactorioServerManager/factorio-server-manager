package main

import (
	"io/ioutil"
	"log"
	"time"
)

type Save struct {
	Name    string    `json:"name"`
	LastMod time.Time `json:"last_mod"`
	Size    int64     `json:"size"`
}

// Lists save files in factorio/saves
func listSaves() []Save {
	saveDir := config.FactorioDir + "/saves"
	result := []Save{}

	files, err := ioutil.ReadDir(saveDir)
	if err != nil {
		log.Printf("Error listing save directory: %s", err)
		return result
	}

	for _, f := range files {
		save := Save{f.Name(), f.ModTime(), f.Size()}
		result = append(result, save)
	}

	return result
}
