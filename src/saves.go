package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Save struct {
	Name    string    `json:"name"`
	LastMod time.Time `json:"last_mod"`
	Size    int64     `json:"size"`
}

func (s Save) String() string {
	return s.Name
}

// Lists save files in factorio/saves
func listSaves(saveDir string) (saves []Save, err error) {
	err = filepath.Walk(saveDir, func(path string, info os.FileInfo, err error) error {
		saves = append(saves, Save{
			info.Name(),
			info.ModTime(),
			info.Size(),
		})
		return nil
	})
	return
}

func findSave(name string) (*Save, error) {
	saves, err := listSaves(config.FactorioSavesDir)
	if err != nil {
		return nil, fmt.Errorf("error listing saves: %v", err)
	}

	for _, save := range saves {
		if save.Name == name {
			return &save, nil
		}
	}

	return nil, errors.New("save not found")
}

func (s *Save) remove() error {
	if s.Name == "" {
		return errors.New("save name cannot be blank")
	}

	return os.Remove(filepath.Join(config.FactorioSavesDir, s.Name))
}
