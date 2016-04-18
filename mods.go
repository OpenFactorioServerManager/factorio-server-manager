package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type ModList struct {
	Mods []Mod `json:"mods"`
}

type Mod struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled,string"`
}

func listInstalledMods() []string {
	modDir := config.FactorioDir + "/mods"
	result := []string{}

	files, err := ioutil.ReadDir(modDir)
	if err != nil {
		log.Printf("Error listing installed mods")
		return result
	}
	for _, f := range files {
		if f.Name() == "mod-list.json" {
			continue
		}
		result = append(result, f.Name())
	}

	return result
}

func parseModList() (ModList, error) {
	var mods ModList
	modListFile := config.FactorioDir + "/mods/mod-list.json"

	modList, err := ioutil.ReadFile(modListFile)
	if err != nil {
		log.Printf("Error reading mod-list.json file: %s", err)
		return mods, err
	}

	err = json.Unmarshal(modList, &mods)
	if err != nil {
		log.Printf("Error parsing mod-list.json JSON: %s", err)
		return mods, err
	}

	return mods, nil

}

func (m *ModList) toggleMod(name string) error {
	found := false

	for i := range m.Mods {
		if m.Mods[i].Name == name {
			found = true
			if m.Mods[i].Enabled == true {
				m.Mods[i].Enabled = false
			} else {
				m.Mods[i].Enabled = true
			}
		}
	}

	if found == false {
		err := fmt.Errorf("Mod with name %s not found", name)
		return err
	}

	m.save()

	return nil
}

func (m ModList) save() error {
	modListFile := config.FactorioDir + "/mods/mod-list.json"
	b, _ := json.MarshalIndent(m, "", "    ")

	err := ioutil.WriteFile(modListFile, b, 0644)
	if err != nil {
		log.Printf("Error writing to mod-list.json file: %s", err)
	}

	return nil
}
