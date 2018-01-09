package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"lockfile"
	"log"
	"os"
	"path/filepath"
)

type ModInfoList struct {
	Mods        []ModInfo `json:"mods"`
	Destination string    `json:"-"`
}
type ModInfo struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	FileName        string `json:"file_name"`
	FactorioVersion string `json:"factorio_version"`
}

func newModInfoList(destination string) (ModInfoList, error) {
	var err error
	modInfoList := ModInfoList{
		Destination: destination,
	}

	err = modInfoList.listInstalledMods()
	if err != nil {
		log.Printf("ModInfoList ... error listing installed Mods: %s", err)
		return modInfoList, err
	}

	return modInfoList, nil
}

func (modInfoList *ModInfoList) listInstalledMods() error {
	var err error
	modInfoList.Mods = nil

	err = filepath.Walk(modInfoList.Destination, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".zip" {
			err = fileLock.RLock(path)
			if err != nil && err == lockfile.ErrorAlreadyLocked {
				log.Println(err)
				return nil
			} else if err != nil {
				log.Printf("error locking file: %s", err)
				return err
			}
			defer fileLock.RUnlock(path)

			zipFile, err := zip.OpenReader(path)
			if err != nil {
				log.Fatalln(err)
				return err
			}

			var modInfo ModInfo
			err = modInfo.getModInfo(&zipFile.Reader)
			if err != nil {
				log.Fatalf("Error in getModInfo: %s", err)
			}

			modInfo.FileName = info.Name()
			modInfoList.Mods = append(modInfoList.Mods, modInfo)
		}

		return nil
	})

	if err != nil {
		log.Printf("error while walking over the given dir: %s", err)
		return err
	}

	return nil
}

func (modInfoList *ModInfoList) deleteMod(modName string) error {
	var err error

	//search for mod, that should be deleted
	for _, mod := range modInfoList.Mods {
		if mod.Name == modName {
			filePath := modInfoList.Destination + "/" + mod.FileName

			fileLock.LockW(filePath)
			//delete mod
			err = os.Remove(filePath)
			fileLock.Unlock(filePath)
			if err != nil {
				log.Printf("ModInfoList ... error when deleting mod: %s", err)
				return err
			}

			//reload mod-list
			err = modInfoList.listInstalledMods()
			if err != nil {
				log.Printf("ModInfoList ... error while refreshing installedModList: %s", err)
				return err
			}

			return nil
		}
	}

	log.Printf("the mod-file for mod %s doesntt exists!", modName)
	return nil
}

func (modInfo *ModInfo) getModInfo(reader *zip.Reader) error {
	for _, singleFile := range reader.File {
		if singleFile.FileInfo().Name() == "info.json" {
			//interpret info.json
			rc, err := singleFile.Open()

			if err != nil {
				log.Fatal(err)
				return err
			}

			byteArray, err := ioutil.ReadAll(rc)
			rc.Close()
			if err != nil {
				log.Fatal(err)
				return err
			}

			err = json.Unmarshal(byteArray, modInfo)
			if err != nil {
				log.Fatalln(err)
				return err
			}

			return nil
		}
	}

	return errors.New("info.json not found in zip-file")
}

func (modInfoList *ModInfoList) createMod(modName string, fileName string, modFile io.Reader) error {
	var err error

	//save uploaded file
	filePath := modInfoList.Destination + "/" + fileName
	newFile, err := os.Create(filePath)
	if err != nil {
		log.Printf("error on creating new file - %s: %s", fileName, err)
		return err
	}
	defer newFile.Close()

	fileLock.LockW(filePath)

	_, err = io.Copy(newFile, modFile)
	if err != nil {
		log.Printf("error on copying file to disk: %s", err)
		return err
	}

	err = newFile.Close()
	if err != nil {
		log.Printf("error on closing new created zip-file: %s", err)
		return err
	}

	fileLock.Unlock(filePath)

	//reload the list
	err = modInfoList.listInstalledMods()
	if err != nil {
		log.Printf("error on listing mod-infos: %s", err)
		return err
	}

	return nil
}
