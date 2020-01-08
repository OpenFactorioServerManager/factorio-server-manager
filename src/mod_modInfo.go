package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"github.com/mroote/factorio-server-manager/lockfile"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ModInfoList struct {
	Mods        []ModInfo `json:"mods"`
	Destination string    `json:"-"`
}
type ModInfo struct {
	Name            string   `json:"name"`
	Version         string   `json:"version"`
	Title           string   `json:"title"`
	Author          string   `json:"author"`
	FileName        string   `json:"file_name"`
	FactorioVersion Version  `json:"factorio_version"`
	Dependencies    []string `json:"dependencies"`
	Compatibility   bool     `json:"compatibility"`
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
			defer zipFile.Close()

			var modInfo ModInfo
			err = modInfo.getModInfo(&zipFile.Reader)
			if err != nil {
				log.Fatalf("Error in getModInfo: %s", err)
			}

			modInfo.FileName = info.Name()

			var base Version
			var op string
			for _, dep := range modInfo.Dependencies {
				dep = strings.TrimSpace(dep)
				if dep == "" {
					continue
				}

				parts := strings.Split(dep, " ")
				if len(parts) > 3 {
					log.Printf("skipping dependency '%s' in '%s': invalid format\n", dep, modInfo.Name)
					continue
				}
				if parts[0] != "base" {
					continue
				}
				if len(parts) == 1 {
					base = modInfo.FactorioVersion
					op = ">="
					continue
				}

				op = parts[1]

				if err := base.UnmarshalText([]byte(parts[2])); err != nil {
					log.Printf("skipping dependency '%s' in '%s': %v\n", dep, modInfo.Name, err)
					continue
				}

				break
			}

			if !base.Equals(NilVersion) {
				modInfo.Compatibility = FactorioServ.Version.Compare(base, op)
			} else {
				log.Println("error finding basemodDependency. Using FactorioVersion...")
				modInfo.Compatibility = !FactorioServ.Version.Less(modInfo.FactorioVersion)
			}

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
			filePath := filepath.Join(modInfoList.Destination, mod.FileName)

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

	log.Printf("the mod-file for mod %s doesn't exists!", modName)
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
			if err != nil {
				log.Fatal(err)
				return err
			}
			err = rc.Close()
			if err != nil {
				log.Printf("Error closing singleFile: %s", err)
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
	filePath := filepath.Join(modInfoList.Destination, fileName)
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
