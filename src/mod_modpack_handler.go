package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func CheckModPackExists(modPackMap ModPackMap, modPackName string, w http.ResponseWriter, resp interface{}) error {
	exists := modPackMap.checkModPackExists(modPackName)
	if !exists {
		resp = fmt.Sprintf("requested modPack {%s} does not exist", modPackName)
		log.Println(resp)
		w.WriteHeader(http.StatusNotFound)
		return errors.New("requested modPack does not exist")
	}
	return nil
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

func ReadModPackRequest(w http.ResponseWriter, r *http.Request, resp *interface{}) (err error, packMap ModPackMap, modPackName string) {
	vars := mux.Vars(r)
	modPackName = vars["modpack"]

	packMap, err = CreateNewModPackMap(w, resp)
	if err != nil {
		return
	}

	if err = CheckModPackExists(packMap, modPackName, w, resp); err != nil {
		return
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

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	var modPackStruct struct {
		Name string `json:"name"`
	}
	ReadFromRequestBody(w, r, &resp, &modPackStruct)
	if err != nil {
		return
	}

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

	err, modPackMap, modPackName := ReadModPackRequest(w, r, &resp)
	if err != nil {
		return
	}

	err = modPackMap.deleteModPack(modPackName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = fmt.Sprintf("Error deleting modpack file: %s", err)
		log.Println(resp)
		return
	}

	resp = modPackName
}

func ModPackDownloadHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	err, _, modPackName := ReadModPackRequest(w, r, &resp)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		WriteResponse(w, resp)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", modPackName))

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	//iterate over folder and create everything in the zip
	err = filepath.Walk(filepath.Join(config.FactorioModPackDir, modPackName), func(path string, info os.FileInfo, err error) error {
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

	w.Header().Set("Content-Type", "application/zip;charset=UTF-8")
}

func ModPackLoadHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	err, modPackMap, modPackName := ReadModPackRequest(w, r, &resp)
	if err != nil {
		return
	}

	err = modPackMap[modPackName].loadModPack()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = fmt.Sprintf("Error loading modpack file: %s", err)
		log.Println(resp)
		return
	}

	resp = modPackMap[modPackName].Mods.listInstalledMods()
}

//////////////////////////////////
// Mods inside Mod Pack Handler //
//////////////////////////////////
func ModPackModListHandler(w http.ResponseWriter, r *http.Request) {
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	err, modPackMap, modPackName := ReadModPackRequest(w, r, &resp)
	if err != nil {
		return
	}

	resp = modPackMap[modPackName].Mods.listInstalledMods()
}

func ModPackModToggleHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	err, packMap, packName := ReadModPackRequest(w, r, &resp)
	if err != nil {
		return
	}

	var modPackStruct struct {
		ModName string `json:"name"`
	}
	ReadFromRequestBody(w, r, &resp, &modPackStruct)
	if err != nil {
		return
	}

	err, resp = packMap[packName].Mods.ModSimpleList.toggleMod(modPackStruct.ModName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = fmt.Sprintf("Error toggling mod inside modPack: %s", err)
		log.Println(resp)
		return
	}
}

func ModPackModDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	err, packMap, packName := ReadModPackRequest(w, r, &resp)
	if err != nil {
		return
	}

	var modPackStruct struct {
		Name string `json:"name"`
	}
	err = ReadFromRequestBody(w, r, &resp, &modPackStruct)
	if err != nil {
		return
	}

	err = packMap[packName].Mods.deleteMod(modPackStruct.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = fmt.Sprintf("Error deleting mod {%s} in modpack {%s}: %s", modPackStruct.Name, packName, err)
		log.Println(resp)
		return
	}

	resp = true
}

func ModPackModUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//Get Data out of the request
	var modPackStruct struct {
		modName     string `json:"modName"`
		downloadUrl string `json:"downloadUrl"`
		filename    string `json:"filename"`
	}
	err = ReadFromRequestBody(w, r, &resp, &modPackStruct)
	if err != nil {
		return
	}

	err, packMap, packName := ReadModPackRequest(w, r, &resp)
	if err != nil {
		return
	}

	err = packMap[packName].Mods.updateMod(modPackStruct.modName, modPackStruct.downloadUrl, modPackStruct.filename)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = fmt.Sprintf("Error updating mod {%s} in modpack {%s}: %s", modPackStruct.modName, packName, err)
		log.Println(resp)
		return
	}

	installedMods := packMap[packName].Mods.listInstalledMods().ModsResult
	var found = false
	for _, mod := range installedMods {
		if mod.Name == modPackStruct.modName {
			resp = mod
			found = true
			return
		}
	}

	if !found {
		resp = fmt.Sprintf(`Could not find mod %s`, modPackStruct.modName)
		log.Println(resp)
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func ModPackModDeleteAllHandler(w http.ResponseWriter, r *http.Request) {
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	err, packMap, packName := ReadModPackRequest(w, r, &resp)
	if err != nil {
		return
	}

	// Delete Modpack
	err = packMap.deleteModPack(packName)
	if err != nil {
		resp = fmt.Sprintf("Error deleting modPackDir: %s", err)
		log.Println(resp)
		return
	}

	// recreate modPack without mods
	err = packMap.createEmptyModPack(packName)
	if err != nil {
		resp = fmt.Sprintf("Error recreating modPackDir: %s", err)
		log.Println(resp)
		return
	}

	resp = true
}

func ModPackModUploadHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponseFileInput{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	r.ParseMultipartForm(32 << 20)

	err, packMap, packName := ReadModPackRequest(w, r, &resp.Data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error in uploadMod, listing mods wasn't successful: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in uploadMod, listing mods wasn't successful: %s", err)
		}
		return
	}

	mods := packMap[packName].Mods
	if err == nil {
		for fileKey, modFile := range r.MultipartForm.File["mod_file"] {
			err = mods.uploadMod(modFile)
			if err != nil {
				resp.ErrorKeys = append(resp.ErrorKeys, fileKey)
				resp.Error = "An error occurred during upload or saving, pls check manually, if all went well and delete invalid files. (This program also could be crashed)"
			}
		}
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error in uploadMod, listing mods wasn't successful: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in uploadMod, listing mods wasn't successful: %s", err)
		}
		return
	}

	resp.Data = mods.listInstalledMods()
	resp.Success = true

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in ModUploadHandler: %s", err)
	}
}
