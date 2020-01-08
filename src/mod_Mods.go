package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"github.com/mroote/factorio-server-manager/lockfile"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type Mods struct {
	ModSimpleList ModSimpleList `json:"mod_simple_list"`
	ModInfoList   ModInfoList   `json:"mod_info_list"`
}
type ModsResult struct {
	ModInfo
	Enabled bool `json:"enabled"`
}
type ModsResultList struct {
	ModsResult []ModsResult `json:"mods"`
}

var fileLock lockfile.FileLock = lockfile.NewLock()

func newMods(destination string) (Mods, error) {
	var err error
	var mods Mods

	mods.ModSimpleList, err = newModSimpleList(destination)
	if err != nil {
		log.Printf("error on creating newModSimpleList: %s", err)
		return mods, err
	}

	mods.ModInfoList, err = newModInfoList(destination)
	if err != nil {
		log.Printf("error on creating newModInfoList: %s", err)
		return mods, err
	}

	return mods, nil
}

func (mods *Mods) listInstalledMods() ModsResultList {
	var result ModsResultList

	for _, modInfo := range mods.ModInfoList.Mods {
		var modsResult ModsResult
		modsResult.Name = modInfo.Name
		modsResult.FileName = modInfo.FileName
		modsResult.Author = modInfo.Author
		modsResult.Title = modInfo.Title
		modsResult.Version = modInfo.Version
		modsResult.FactorioVersion = modInfo.FactorioVersion
		modsResult.Compatibility = modInfo.Compatibility

		for _, simpleMod := range mods.ModSimpleList.Mods {
			if simpleMod.Name == modsResult.Name {
				modsResult.Enabled = simpleMod.Enabled
				break
			}
		}

		result.ModsResult = append(result.ModsResult, modsResult)
	}

	return result
}

func (mods *Mods) deleteMod(modName string) error {
	var err error

	err = mods.ModInfoList.deleteMod(modName)
	if err != nil {
		log.Printf("error when deleting mod in ModInfoList: %s", err)
		return err
	}

	err = mods.ModSimpleList.deleteMod(modName)
	if err != nil {
		log.Printf("error when deleting mod in ModSimpleList: %s", err)
		return err
	}

	return nil
}

func (mods *Mods) createMod(modName string, fileName string, fileRc io.Reader) error {
	var err error

	//check if mod already exists and delete it
	if mods.ModSimpleList.checkModExists(modName) {
		err = mods.ModSimpleList.deleteMod(modName)
		if err != nil {
			log.Printf("error when deleting mod: %s", err)
			return err
		}

		err = mods.ModInfoList.deleteMod(modName)
		if err != nil {
			log.Printf("error when deleting mod: %s", err)
			return err
		}
	}

	//create new mod
	err = mods.ModInfoList.createMod(modName, fileName, fileRc)
	if err != nil {
		log.Printf("error on creating mod-file: %s", err)
		return err
	}
	err = mods.ModSimpleList.createMod(modName)
	if err != nil {
		log.Printf("error on adding mod to the mod-list.json: %s", err)
		return err
	}

	return nil
}

func (mods *Mods) downloadMod(url string, filename string, modId string) error {
	var err error

	var credentials FactorioCredentials
	status, err := credentials.load()
	if err != nil {
		log.Printf("error loading credentials: %s", err)
		return err
	}
	if status == false {
		log.Printf("error: credentials are invalid")
		return errors.New("error: credentials are invalid")
	}

	//download the mod from the mod portal api
	completeUrl := "https://mods.factorio.com" + url + "?username=" + credentials.Username + "&token=" + credentials.Userkey

	response, err := http.Get(completeUrl)
	if err != nil {
		log.Printf("error on downloading mod: %s", err)
		return err
	}

	log.Printf("download complete\n StatusCode: %d\n Status: %s", response.StatusCode, response.Status)

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("StatusCode: %d", response.StatusCode)

		return errors.New("Statuscode not 200: " + fmt.Sprint(response.StatusCode))
	}

	err = mods.createMod(modId, filename, response.Body)
	if err != nil {
		log.Printf("error when creating Mod: %s", err)
		return err
	}

	log.Printf("completed copying the response.Body")

	//done everything is made inside the createMod

	return nil
}

func (mods *Mods) uploadMod(header *multipart.FileHeader) error {
	var err error

	if filepath.Ext(header.Filename) != ".zip" {
		log.Print("The uploaded file wasn't a zip-file -> ignore it")
		return nil //simply do nothing xD
	}

	file, err := header.Open()
	if err != nil {
		log.Printf("error on open file via fileHeader. %s", err)
		return err
	}
	defer file.Close()

	fileByteArray, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("error reading file: %s", err)
		return err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(fileByteArray), int64(len(fileByteArray)))
	if err != nil {
		log.Printf("Uploaded file could not put into zip.Reader: %s", err)
		return err
	}

	var modInfo ModInfo
	err = modInfo.getModInfo(zipReader)
	if err != nil {
		log.Printf("Error in getModInfo: %s", err)
		return err
	}

	err = mods.createMod(modInfo.Name, header.Filename, bytes.NewReader(fileByteArray))
	if err != nil {
		log.Printf("error on creating Mod: %s", err)
		return err
	}

	return nil
}

func (mods *Mods) updateMod(modName string, url string, filename string) error {
	var err error

	err = mods.downloadMod(url, filename, modName)
	if err != nil {
		log.Printf("updateMod ... error when downloading the new Mod: %s", err)
		return err
	}

	return nil
}
