package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ModList struct {
	Mods []Mod `json:"mods"`
}

type Mod struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled,string"`
}

// List mods installed in the factorio/mods directory
func listInstalledMods(modDir string) ([]string, error) {
	result := []string{}

	files, err := ioutil.ReadDir(modDir)
	if err != nil {
		log.Printf("Error listing installed mods: %s", err)
		return result, err
	}
	for _, f := range files {
		if f.Name() == "mod-list.json" {
			continue
		}
		result = append(result, f.Name())
	}

	return result, nil
}

// Delete mod by provided filename
func rmMod(modName string) error {
	removed := false
	if modName == "" {
		return errors.New("No mod name provided.")
	}
	// Get list of installed mods
	installedMods, err := listInstalledMods(config.FactorioModsDir)
	if err != nil {
		log.Printf("Error in remove mod list: %s", err)
		return err
	}

	// Check if provided mod matches one thats installed else return err
	for _, mod := range installedMods {
		if strings.Contains(mod, modName) {
			log.Printf("Removing mod: %s", mod)
			err := os.Remove(filepath.Join(config.FactorioModsDir, mod))
			if err != nil {
				log.Printf("Error removing mod %s: %s", mod, err)
				return err
			}
			removed = true
			log.Printf("Removed mod: %s", mod)
		}
	}

	if !removed {
		log.Printf("Did not remove mod: %s", modName)
		return errors.New(fmt.Sprintf("Did not remove mod: %s", modName))
	}

	return nil
}

func rmModPack(modpack string) error {
	removed := false
	if modpack == "" {
		return errors.New("No mod pack name provided.")
	}
	// Get list of modpacks
	modpacks, err := listModPacks(filepath.Join(config.FactorioDir, "modpacks"))
	if err != nil {
		log.Printf("Error listing modpacks in rmModPack: %s", err)
		return err
	}

	for _, m := range modpacks {
		if strings.Contains(m, modpack) {
			log.Printf("Removing modpack: %s", m)
			err := os.Remove(filepath.Join(config.FactorioDir, "modpacks", m))
			if err != nil {
				log.Printf("Error trying to remove modpack: %s: %s", m, err)
				return err
			}
			removed = true
			log.Printf("Removed modpack: %s", m)
		}
	}

	if !removed {
		log.Printf("Did not remove modpack: %s", modpack)
		return errors.New(fmt.Sprintf("Did not remove modpack: %s", modpack))
	}

	return nil
}

func createModPackDir() error {
	err := os.Mkdir(filepath.Join(config.FactorioDir, "modpacks"), 0775)
	if err != nil {
		log.Printf("Could not create modpacks directory: %s", err)
		return err
	}

	return nil
}

// Create's modpack zip file from provided title, mods parameter is a string of mod filenames
func createModPack(title string, mods ...string) error {
	zipfile, err := os.Create(filepath.Join(config.FactorioDir, "modpacks", title+".zip"))
	if err != nil {
		log.Printf("Error creating zipfile: %s, error: %s", title, err)
	}
	defer zipfile.Close()
	// Create Zip writer
	z := zip.NewWriter(zipfile)
	defer z.Close()

	for _, mod := range mods {
		// Process mod file, add to zipfile
		f, err := os.Open(filepath.Join(config.FactorioDir, "mods", mod))
		if err != nil {
			log.Printf("Error creating modpack file %s for archival: ", mod, err)
			return err
		}
		// Read contents of mod to be compressed
		modfile, err := ioutil.ReadAll(f)
		if err != nil {
			log.Printf("Error reading modfile contents: %s", err)
			continue
		}
		// Add file to zip archive
		fmt.Println(mod)
		zip, err := z.Create(mod)
		if err != nil {
			log.Printf("Error adding file: %s to zip: %s", f.Name, err)
			continue
		}
		// Write file contents to zip archive
		_, err = zip.Write(modfile)
		if err != nil {
			log.Printf("Error writing to zipfile: %s", err)
			continue
		}
	}

	err = z.Close()
	if err != nil {
		log.Printf("Error trying to zip: %s, error: %s", title, err)
	}

	return nil
}

func listModPacks(modDir string) ([]string, error) {
	result := []string{}

	files, err := ioutil.ReadDir(modDir)
	if err != nil {
		log.Printf("Error listing modpacks: %s", err)
		return result, err
	}
	for _, f := range files {
		result = append(result, f.Name())
	}

	return result, nil
}

// Parses mod-list.json file in factorio/mods
// returns ModList struct
func parseModList() (ModList, error) {
	var mods ModList
	modListFile := filepath.Join(config.FactorioModsDir, "mod-list.json")

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

// Toggles Enabled boolean for mod specified in name parameter in mod-list.json file
func (m *ModList) toggleMod(name string) error {
	found := false
	status := false

	for i := range m.Mods {
		if m.Mods[i].Name == name {
			found = true
			if m.Mods[i].Enabled == true {
				m.Mods[i].Enabled = false
			} else {
				m.Mods[i].Enabled = true
				status = true
			}
		}
	}

	if found {
		err := m.save()
		if err != nil {
			log.Printf("Error saving changes to mod-list-.json file: %s", err)
			return err
		}
		log.Printf("Mod: %s was toggled to %v", name, status)
	}

	return nil
}

// Saves ModList object to mod-list.json file
// Overwrites old file
func (m ModList) save() error {
	modListFile := filepath.Join(config.FactorioModsDir, "mod-list.json")
	b, _ := json.MarshalIndent(m, "", "    ")

	err := ioutil.WriteFile(modListFile, b, 0644)
	if err != nil {
		log.Printf("Error writing to mod-list.json file: %s", err)
		return err
	}

	return nil
}

//TODO Add method to allow downloading all installed mods in zip file
