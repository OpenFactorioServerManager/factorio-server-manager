package main

import (
	"encoding/json"
	"flag"
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

type Config struct {
	FactorioDir string
}

var config Config

func listInstalledMods() []string {
	modDir := config.FactorioDir + "/mods"
	result := []string{}

	files, err := ioutil.ReadDir(modDir)
	if err != nil {
		log.Printf("Error listing installed mods")
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
	fmt.Println(m)
	for i := range m.Mods {
		if m.Mods[i].Name == name {
			if m.Mods[i].Enabled == true {
				m.Mods[i].Enabled = false
			} else {
				m.Mods[i].Enabled = true
			}
		}
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

func main() {
	factorioDir := flag.String("dir", "./", "Specify location of Factorio config directory.")
	flag.Parse()

	config.FactorioDir = *factorioDir

	fmt.Println(listInstalledMods())

	//m, err := parseModList()
	//if err != nil {
	//	log.Printf("Error")
	//}
}
