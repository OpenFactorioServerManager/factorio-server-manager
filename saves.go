package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type Save struct {
	Name    string    `json:"name"`
	LastMod time.Time `json:"last_mod"`
	Size    int64     `json:"size"`
}

func (s Save) String() string {
	return fmt.Sprintf("%s", s.Name)
}

// Lists save files in factorio/saves
func listSaves(saveDir string) ([]Save, error) {
	result := []Save{}

	files, err := ioutil.ReadDir(saveDir)
	if err != nil {
		log.Printf("Error listing save directory: %s", err)
		return result, err
	}

	for _, f := range files {
		save := Save{f.Name(), f.ModTime(), f.Size()}
		result = append(result, save)
	}

	return result, nil
}

func rmSave(saveName string) error {
	removed := false
	if saveName == "" {
		return errors.New("No save name provided")
	}

	saves, err := listSaves(config.FactorioSavesDir)
	if err != nil {
		log.Printf("Error in remove save: %s", err)
		return err
	}

	for _, save := range saves {
		log.Printf("Checking if %s in %s", save, saveName)
		if strings.Contains(save.Name, saveName) {
			err := os.Remove(config.FactorioSavesDir + "/" + save.Name)
			if err != nil {
				log.Printf("Error removing save %s: %s", saveName, err)
				return err
			}
			log.Printf("Deleted save: %s", save)
			removed = true
		}
	}

	if !removed {
		log.Printf("Did not remove save: %s", saveName)
		return errors.New(fmt.Sprintf("Did not remove save: %s", saveName))
	}

	return nil
}
