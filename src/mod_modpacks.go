package main

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type ModPackMap map[string]*ModPack
type ModPack struct {
	Mods Mods
}

type ModPackResult struct {
	Name string         `json:"name"`
	Mods ModsResultList `json:"mods"`
}
type ModPackResultList struct {
	ModPacks []ModPackResult `json:"mod_packs"`
}

func newModPackMap() (ModPackMap, error) {
	var err error
	modPackMap := make(ModPackMap)

	err = modPackMap.reload()
	if err != nil {
		log.Printf("error on loading the modpacks: %s", err)
		return modPackMap, err
	}

	return modPackMap, nil
}

func newModPack(modPackFolder string) (*ModPack, error) {
	var err error
	var modPack ModPack

	modPack.Mods, err = newMods(modPackFolder)
	if err != nil {
		log.Printf("error on loading mods in mod_pack_dir: %s", err)
		return &modPack, err
	}

	return &modPack, err
}

func (modPackMap *ModPackMap) reload() error {
	var err error
	newModPackMap := make(ModPackMap)

	err = filepath.Walk(config.FactorioModPackDir, func(path string, info os.FileInfo, err error) error {
		if path == config.FactorioModPackDir || !info.IsDir() {
			return nil
		}

		modPackName := filepath.Base(path)

		newModPackMap[modPackName], err = newModPack(path)
		if err != nil {
			log.Printf("error on creating newModPack: %s", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Printf("error on walking over the ModDir: %s", err)
		return err
	}

	*modPackMap = newModPackMap

	return nil
}

func (modPackMap *ModPackMap) listInstalledModPacks() ModPackResultList {
	var modPackResultList ModPackResultList

	for modPackName, modPack := range *modPackMap {
		var modPackResult ModPackResult
		modPackResult.Name = modPackName
		modPackResult.Mods = modPack.Mods.listInstalledMods()

		modPackResultList.ModPacks = append(modPackResultList.ModPacks, modPackResult)
	}

	return modPackResultList
}

func (modPackMap *ModPackMap) createModPack(modPackName string) error {
	var err error

	modPackFolder := filepath.Join(config.FactorioModPackDir, modPackName)

	if modPackMap.checkModPackExists(modPackName) == true {
		log.Printf("ModPack %s already existis", modPackName)
		return errors.New("ModPack " + modPackName + " already exists, please choose a different name")
	}

	sourceFileInfo, err := os.Stat(config.FactorioModsDir)
	if err != nil {
		log.Printf("error when reading factorioModsDir. %s", err)
		return err
	}

	//Create the modPack-folder
	err = os.MkdirAll(modPackFolder, sourceFileInfo.Mode())
	if err != nil {
		log.Printf("error on creating the new ModPack directory: %s", err)
		return err
	}

	files, err := ioutil.ReadDir(config.FactorioModsDir)
	if err != nil {
		log.Printf("error on reading the dactorio mods dir: %s", err)
		return err
	}

	for _, file := range files {
		if file.IsDir() == false {
			sourceFilepath := filepath.Join(config.FactorioModsDir, file.Name())
			destinationFilepath := filepath.Join(modPackFolder, file.Name())

			sourceFile, err := os.Open(sourceFilepath)
			if err != nil {
				log.Printf("error on opening sourceFilepath: %s", err)
				return err
			}
			defer sourceFile.Close()

			destinationFile, err := os.Create(destinationFilepath)
			if err != nil {
				log.Printf("error on creating destinationFilepath: %s", err)
				return err
			}
			defer destinationFile.Close()

			_, err = io.Copy(destinationFile, sourceFile)
			if err != nil {
				log.Printf("error on copying data from source to destination: %s", err)
				return err
			}

			sourceFile.Close()
			destinationFile.Close()
		}
	}

	//reload the ModPackList
	err = modPackMap.reload()
	if err != nil {
		log.Printf("error on reloading ModPack: %s", err)
		return err
	}

	return nil
}

func (modPackMap *ModPackMap) checkModPackExists(modPackName string) bool {
	for modPackId := range *modPackMap {
		if modPackId == modPackName {
			return true
		}
	}

	return false
}

func (modPackMap *ModPackMap) deleteModPack(modPackName string) error {
	var err error

	modPackDir := filepath.Join(config.FactorioModPackDir, modPackName)

	err = os.RemoveAll(modPackDir)
	if err != nil {
		log.Printf("error on removing the ModPack: %s", err)
		return err
	}

	err = modPackMap.reload()
	if err != nil {
		log.Printf("error on reloading the ModPackList: %s", err)
		return err
	}

	return nil
}

func (modPack *ModPack) loadModPack() error {
	var err error

	//get filemode, so it can be restored
	fileInfo, err := os.Stat(config.FactorioModsDir)
	if err != nil {
		log.Printf("error on trying to save folder infos: %s", err)
		return err
	}
	folderMode := fileInfo.Mode()

	//clean factorio mod directory
	err = os.RemoveAll(config.FactorioModsDir)
	if err != nil {
		log.Printf("error on removing the factorio mods dir: %s", err)
		return err
	}

	err = os.Mkdir(config.FactorioModsDir, folderMode)
	if err != nil {
		log.Printf("error on recreating mod dir: %s", err)
		return err
	}

	//copy the modpack folder to the normal mods directory
	err = filepath.Walk(modPack.Mods.ModInfoList.Destination, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		newFile, err := os.Create(filepath.Join(config.FactorioModsDir, info.Name()))
		if err != nil {
			log.Printf("error on creting mod file: %s", err)
			return err
		}
		defer newFile.Close()

		oldFile, err := os.Open(path)
		if err != nil {
			log.Printf("error on opening modFile: %s", err)
			return err
		}
		defer oldFile.Close()

		_, err = io.Copy(newFile, oldFile)
		if err != nil {
			log.Printf("error on copying data to the new file: %s", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Printf("error on copying the mod pack: %s", err)
		return err
	}

	return nil
}
