package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type LoginErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}
type LoginSuccessResponse struct {
	UserKey []string `json:""`
}
type FactorioCredentials struct {
	Username string `json:"username"`
	Userkey  string `json:"userkey"`
}

func (credentials *FactorioCredentials) save() error {
	var err error

	credentialsJson, err := json.Marshal(credentials)
	if err != nil {
		log.Printf("error mashalling the credentials: %s", err)
		return err
	}

	err = ioutil.WriteFile(config.FactorioCredentialsFile, credentialsJson, 0664)
	if err != nil {
		log.Printf("error on saving the credentials. %s", err)
		return err
	}

	return nil
}

func (credentials *FactorioCredentials) load() (bool, error) {
	var err error

	if _, err := os.Stat(config.FactorioCredentialsFile); os.IsNotExist(err) {
		return false, nil
	}

	fileBytes, err := ioutil.ReadFile(config.FactorioCredentialsFile)
	if err != nil {
		credentials.del()
		log.Printf("error reading CredentialsFile: %s", err)
		return false, err
	}

	err = json.Unmarshal(fileBytes, credentials)
	if err != nil {
		credentials.del()
		log.Printf("error on unmarshal credentials_file: %s", err)
		return false, err
	}

	if credentials.Userkey != "" && credentials.Username != "" {
		return true, nil
	} else {
		credentials.del()
		return false, errors.New("incredients incomplete")
	}
}

func (credentials *FactorioCredentials) del() error {
	var err error

	err = os.Remove(config.FactorioCredentialsFile)
	if err != nil {
		log.Printf("error delete the credentialfile: %s", err)
		return err
	}

	return nil
}

//Log the user into factorio, so mods can be downloaded
func factorioLogin(username string, password string) (string, error, int) {
	var err error

	resp, err := http.PostForm("https://auth.factorio.com/api-login",
		url.Values{"require_game_ownership": {"true"}, "username": {username}, "password": {password}})

	if err != nil {
		log.Printf("error on logging in: %s", err)
		return "", err, resp.StatusCode
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error on reading resp.Body: %s", err)
		return "", err, http.StatusInternalServerError
	}

	bodyString := string(bodyBytes)

	if resp.StatusCode != http.StatusOK {
		log.Println("error Statuscode not 200")
		return bodyString, errors.New(bodyString), resp.StatusCode
	}

	var successResponse []string
	err = json.Unmarshal(bodyBytes, &successResponse)
	if err != nil {
		log.Printf("error on unmarshal body: %s", err)
		return err.Error(), err, http.StatusInternalServerError
	}

	credentials := FactorioCredentials{
		Username: username,
		Userkey:  successResponse[0],
	}

	err = credentials.save()
	if err != nil {
		log.Printf("error saving the credentials. %s", err)
		return err.Error(), err, http.StatusInternalServerError
	}

	return "", nil, http.StatusOK
}

//Search inside the factorio mod portal
func searchModPortal(keyword string) (string, error, int) {
	req, err := http.NewRequest(http.MethodGet, "https://mods.factorio.com/api/mods", nil)
	if err != nil {
		return "error", err, 500
	}

	query := req.URL.Query()
	query.Add("q", keyword)
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "error", err, 500
	}

	text, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "error", err, 500
	}

	textString := string(text)

	return textString, nil, resp.StatusCode
}

func getModDetails(modId string) (string, error, int) {
	var err error
	newLink := "https://mods.factorio.com/api/mods/" + modId
	resp, err := http.Get(newLink)

	if err != nil {
		return "error", err, http.StatusInternalServerError
	}

	//get the response-text
	text, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	textString := string(text)

	if err != nil {
		log.Fatal(err)
		return "error", err, resp.StatusCode
	}

	return textString, nil, resp.StatusCode
}

func deleteAllMods() error {
	var err error

	modsDirInfo, err := os.Stat(config.FactorioModsDir)
	if err != nil {
		log.Printf("error getting stats of FactorioModsDir: %s", err)
		return err
	}

	modsDirPerm := modsDirInfo.Mode().Perm()

	err = os.RemoveAll(config.FactorioModsDir)
	if err != nil {
		log.Printf("removing FactorioModsDir failed: %s", err)
		return err
	}

	err = os.Mkdir(config.FactorioModsDir, modsDirPerm)
	if err != nil {
		log.Printf("error recreating modPackDir: %s", err)
		return err
	}

	return nil
}

func modStartUp() {
	var err error

	//get main-folder info
	factorioDirInfo, err := os.Stat(config.FactorioDir)
	if err != nil {
		log.Printf("error getting stats from FactorioDir: %s", err)
		return
	}
	factorioDirPerm := factorioDirInfo.Mode().Perm()

	//create mods dir
	if _, err = os.Stat(config.FactorioModsDir); os.IsNotExist(err) {
		log.Println("no mods dir found ... creating one ...")
		os.Mkdir(config.FactorioModsDir, factorioDirPerm)
	}

	//crate mod_pack dir
	if _, err = os.Stat(config.FactorioModPackDir); os.IsNotExist(err) {
		log.Println("no ModPackDir found ... creating one ...")
		os.Mkdir(config.FactorioModPackDir, factorioDirPerm)
	}

	oldModpackDir := filepath.Join(config.FactorioDir, "modpacks")
	if _, err := os.Stat(filepath.Join(oldModpackDir)); !os.IsNotExist(err) {
		log.Printf("found old modpack files, rebuild into new system...")

		err = filepath.Walk(oldModpackDir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(info.Name()) != ".zip" {
				log.Printf("file is not a zip or a directory -> skip")
				return nil
			}

			filename := info.Name()
			n := strings.LastIndexByte(info.Name(), '.')
			modPackName := filename[:n]

			log.Printf("loading modPack %s into new system ...", modPackName)

			modPackDir := filepath.Join(config.FactorioModPackDir, modPackName)

			if _, err := os.Stat(modPackDir); !os.IsNotExist(err) {
				log.Printf("modPack already exists")
				return errors.New("modPack already exists")
			}

			err = os.Mkdir(modPackDir, factorioDirPerm)
			if err != nil {
				log.Printf("error creating newModPackDir: %s", err)
				return err
			}

			//create mod-info.json
			modSimpleList := ModSimpleList{
				Destination: modPackDir,
				Mods: []ModSimple{
					ModSimple{
						Name:    "base",
						Enabled: true,
					},
				},
			}
			newJson, _ := json.Marshal(modSimpleList)

			err = ioutil.WriteFile(modSimpleList.Destination+"/mod-list.json", newJson, 0664)
			if err != nil {
				log.Printf("error when writing new mod-list: %s", err)
				return err
			}

			modPackFile, err := zip.OpenReader(path)
			if err != nil {
				return err
			}
			defer modPackFile.Close()

			mods, err := newMods(modPackDir)
			if err != nil {
				log.Printf("error reading mods: %s", err)
				return err
			}

			for _, modFile := range modPackFile.File {
				modFileRc, err := modFile.Open()
				if err != nil {
					log.Printf("error opening mod_file: %s", err)
					return err
				}
				defer modFileRc.Close()

				modFileBuffer, err := ioutil.ReadAll(modFileRc)
				if err != nil {
					log.Printf("error reading mod_file_rc: %s", err)
					return err
				}
				modFileRc.Close()

				modFileByteReader := bytes.NewReader(modFileBuffer)
				modFileZipReader, err := zip.NewReader(modFileByteReader, int64(len(modFileBuffer)))
				if err != nil {
					log.Printf("error creating Reader on byte_array: %s", err)
					return err
				}

				var modInfo ModInfo
				err = modInfo.getModInfo(modFileZipReader)
				if err != nil {
					log.Printf("error loading the ModInfo: %s", err)
					return err
				}

				err = mods.createMod(modInfo.Name, modFile.Name, bytes.NewReader(modFileBuffer))
				if err != nil {
					log.Printf("error on creating mod: %s", err)
					return err
				}
			}

			log.Printf("loading modPack %s successful", modPackName)

			return nil
		})

		if err != nil {
			log.Printf("error on loading old modpacks into the new system: %s\n please check if empty modPacks are creating and delete them", err)
		} else {
			log.Printf("all modPacks are loaded into the new system successfully")
			log.Printf("deleting old modPackDir")
			os.RemoveAll(oldModpackDir)
		}
	}
}
