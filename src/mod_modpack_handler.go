package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func CheckModPackExists(modPackMap ModPackMap, modPackName string, w http.ResponseWriter, resp interface{}) bool {
	exists := modPackMap.checkModPackExists(modPackName)
	if !exists {
		resp = fmt.Sprintf("requested modPack {%s} doesnt exist", modPackName)
		log.Println(resp)
		w.WriteHeader(http.StatusNotFound)
	}
	return exists
}

func CreateNewModPackMap(w http.ResponseWriter, resp *interface{}) (modPackMap ModPackMap, err error) {
	modPackMap, err = newModPackMap()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		*resp = fmt.Sprintf("Error creating modpackmap aka. list of all modpacks files : %s", err)
		log.Println(resp)
	}
	return
}

//////////////////////
// Mod Pack Handler //
//////////////////////

func ModPackListHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	modPackMap, err := CreateNewModPackMap(w, &resp)
	if err != nil {
		return
	}

	resp = modPackMap.listInstalledModPacks()
}

func ModPackCreateHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	body, err := ReadRequestBody(w, r, &resp)
	if err != nil {
		return
	}

	var modPackStruct struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal(body, &modPackStruct)
	if err != nil {
		resp = fmt.Sprintf("Error unmarshalling saveFile JSON: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	modPackMap, err := CreateNewModPackMap(w, &resp)
	if err != nil {
		return
	}

	err = modPackMap.createModPack(modPackStruct.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = fmt.Sprintf("Error creating modpack file: %s", err)
		log.Println(resp)
		return
	}

	resp = modPackMap.listInstalledModPacks()
}

func ModPackDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	body, err := ReadRequestBody(w, r, &resp)
	if err != nil {
		return
	}

	var modPackStruct struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal(body, &modPackStruct)
	if err != nil {
		resp = fmt.Sprintf("Error unmarshalling saveFile JSON: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	modPackMap, err := CreateNewModPackMap(w, &resp)
	if err != nil {
		return
	}

	if !CheckModPackExists(modPackMap, modPackStruct.Name, w, &resp) {
		return
	}

	err = modPackMap.deleteModPack(modPackStruct.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = fmt.Sprintf("Error deleting modpack file: %s", err)
		log.Println(resp)
		return
	}

	resp = modPackStruct.Name
}

func ModPackDownloadHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	vars := mux.Vars(r)
	modpack := vars["modpack"]

	packMap, err := CreateNewModPackMap(w, &resp)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		WriteResponse(w, resp)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", modpack))

	if CheckModPackExists(packMap, modpack, w, &resp) {
		zipWriter := zip.NewWriter(w)
		defer zipWriter.Close()

		//iterate over folder and create everything in the zip
		err = filepath.Walk(filepath.Join(config.FactorioModPackDir, modpack), func(path string, info os.FileInfo, err error) error {
			if info.IsDir() == false {
				writer, err := zipWriter.Create(info.Name())
				if err != nil {
					log.Printf("error on creating new file inside zip: %s", err)
					return err
				}

				file, err := os.Open(path)
				if err != nil {
					log.Printf("error on opening modfile: %s", err)
					return err
				}
				// Close file, when function returns
				defer func() {
					err2 := file.Close()
					if err == nil && err2 != nil {
						log.Printf("Error closing file: %s", err2)
						err = err2
					}
				}()

				_, err = io.Copy(writer, file)
				if err != nil {
					log.Printf("error on copying file into zip: %s", err)
					return err
				}
			}

			return nil
		})
		if err != nil {
			resp = fmt.Sprintf("error on walking over the modpack: %s", err)
			log.Println(resp)
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json;charset=UTF-8")
			WriteResponse(w, resp)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		WriteResponse(w, resp)
		return
	}

	w.Header().Set("Content-Type", "application/zip;charset=UTF-8")
}

func ModPackLoadHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	body, err := ReadRequestBody(w, r, &resp)
	if err != nil {
		return
	}

	var modPackStruct struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal(body, &modPackStruct)
	if err != nil {
		resp = fmt.Sprintf("Error unmarshalling saveFile JSON: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	modPackMap, err := CreateNewModPackMap(w, &resp)
	if err != nil {
		return
	}

	if !CheckModPackExists(modPackMap, modPackStruct.Name, w, &resp) {
		return
	}

	err = modPackMap[modPackStruct.Name].loadModPack()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = fmt.Sprintf("Error loading modpack file: %s", err)
		log.Println(resp)
		return
	}

	resp = modPackMap[modPackStruct.Name].Mods.listInstalledMods()
}

//////////////////////////////////
// Mods inside Mod Pack Handler //
//////////////////////////////////
