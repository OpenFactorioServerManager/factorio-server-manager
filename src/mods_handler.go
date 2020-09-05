package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mroote/factorio-server-manager/lockfile"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ModPortalStruct struct {
	DownloadsCount int    `json:"downloads_count"`
	Name           string `json:"name"`
	Owner          string `json:"owner"`
	Releases       []struct {
		DownloadURL string `json:"download_url"`
		FileName    string `json:"file_name"`
		InfoJSON    struct {
			FactorioVersion string `json:"factorio_version"`
		} `json:"info_json"`
		ReleasedAt time.Time `json:"released_at"`
		Sha1       string    `json:"sha1"`
		Version    Version   `json:"version"`
	} `json:"releases"`
	Summary string `json:"summary"`
	Title   string `json:"title"`
}

func CreateNewMods(w http.ResponseWriter, resp *interface{}) (mods Mods, err error) {
	mods, err = newMods(config.FactorioModsDir)
	if err != nil {
		*resp = fmt.Sprintf("Error creating mods object: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

// Returns JSON response of all mods installed in factorio/mods
func ListInstalledModsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	mods, err := CreateNewMods(w, &resp)
	if err != nil {
		return
	}

	resp = mods.listInstalledMods().ModsResult
}

func ModToggleHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//Get Data out of the request
	body, err := ReadRequestBody(w, r, &resp)
	if err != nil {
		return
	}

	var mod struct {
		Name string `json:"modName"`
	}
	err = json.Unmarshal(body, &mod)
	if err != nil {
		resp = fmt.Sprintf("Error unmarshalling modName JSON: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mods, err := CreateNewMods(w, &resp)
	if err != nil {
		return
	}

	err, resp = mods.ModSimpleList.toggleMod(mod.Name)
	if err != nil {
		resp = fmt.Sprintf("Error in toggling mod in simple list: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ModDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	// Get Data out of the request
	body, err := ReadRequestBody(w, r, &resp)
	if err != nil {
		return
	}

	var mod struct {
		Name string `json:"modName"`
	}
	err = json.Unmarshal(body, &mod)
	if err != nil {
		resp = fmt.Sprintf("Error unmarshalling modName JSON: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mods, err := CreateNewMods(w, &resp)
	if err != nil {
		return
	}

	mods.deleteMod(mod.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = fmt.Sprintf("Error in deleting mod {%s}: %s", mod.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp = mod.Name
}

func ModDeleteAllHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//delete mods folder
	err = deleteAllMods()
	if err != nil {
		resp = fmt.Sprintf("Error deleting all mods: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp = nil
}

func ModUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//Get Data out of the request
	body, err := ReadRequestBody(w, r, &resp)
	if err != nil {
		return
	}

	var modData struct {
		Name        string `json:"modName"`
		DownloadUrl string `json:"downloadUrl"`
		Filename    string `json:"fileName"`
	}
	err = json.Unmarshal(body, &modData)
	if err != nil {
		resp = fmt.Sprintf("Error unmarshalling modName JSON: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mods, err := CreateNewMods(w, &resp)
	if err != nil {
		return
	}

	err = mods.updateMod(modData.Name, modData.DownloadUrl, modData.Filename)
	if err != nil {
		resp = fmt.Sprintf("Error updating mod {%s}: %s", modData.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	installedMods := mods.listInstalledMods().ModsResult
	for _, mod := range installedMods {
		if mod.Name == modData.Name {
			resp = mod
			return
		}
	}
}

func ModUploadHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponseFileInput{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	r.ParseMultipartForm(32 << 20)

	mods, err := newMods(config.FactorioModsDir)
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
		w.WriteHeader(500)
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

func ModDownloadHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	//iterate over folder and create everything in the zip
	err = filepath.Walk(config.FactorioModsDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() == false {
			//Lock the file, that we are want to read
			err := fileLock.RLock(path)
			if err != nil {
				log.Printf("error locking file for reading, something else has locked it")
				return err
			}
			defer fileLock.RUnlock(path)

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
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				log.Printf("error on copying file into zip: %s", err)
				return err
			}

			err = file.Close()
			if err != nil {
				log.Printf("error closing file: %s", err)
				return err
			}
		}

		return nil
	})
	if err == lockfile.ErrorAlreadyLocked {
		w.WriteHeader(http.StatusLocked)
		return
	}
	if err != nil {
		log.Printf("error on walking over the mods: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writerHeader := w.Header()
	writerHeader.Set("Content-Type", "application/zip;charset=UTF-8")
	writerHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", "all_installed_mods.zip"))
}

//LoadModsFromSaveHandler returns JSON response with the found mods
func LoadModsFromSaveHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//Get Data out of the request
	body, err := ReadRequestBody(w, r, &resp)
	if err != nil {
		return
	}

	var saveFileStruct struct {
		Name string `json:"saveFile"`
	}
	err = json.Unmarshal(body, &saveFileStruct)
	if err != nil {
		resp = fmt.Sprintf("Error unmarshalling saveFile JSON: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	path := filepath.Join(config.FactorioSavesDir, saveFileStruct.Name)
	f, err := OpenArchiveFile(path, "level.dat")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = fmt.Sprintf("cannot open save level file: %v", err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	var header SaveHeader
	err = header.ReadFrom(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = fmt.Sprintf("cannot read save header: %v", err)
		log.Println(resp)
		return
	}

	resp = header
}

func ModPackToggleModHandler(w http.ResponseWriter, r *http.Request) {
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
		modName string `json:"modName"`
		modPack string `json:"modPack"`
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

	err, resp = modPackMap[modPackStruct.modPack].Mods.ModSimpleList.toggleMod(modPackStruct.modName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = fmt.Sprintf("Error toggling mod inside modPack file: %s", err)
		log.Println(resp)
		return
	}
}

func ModPackDeleteModHandler(w http.ResponseWriter, r *http.Request) {
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
		modName string `json:"modName"`
		modPack string `json:"modPack"`
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

	if modPackMap.checkModPackExists(modPackStruct.modPack) {
		err = modPackMap[modPackStruct.modPack].Mods.deleteMod(modPackStruct.modName)
	} else {
		err = errors.New("ModPack " + modPackStruct.modPack + " does not exist")
		resp = fmt.Sprintf("Error loading modpack file: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = fmt.Sprintf("Error deleting mod {%s} in modpack {%s}: %s", modPackStruct.modName, modPackStruct.modPack, err)
		log.Println(resp)
		return
	}

	resp = true
}

func ModPackUpdateModHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//Get Data out of the request
	body, err := ReadRequestBody(w, r, &resp)
	if err != nil {
		return
	}

	var modPackStruct struct {
		modName     string `json:"modName"`
		downloadUrl string `json:"downloadUrl"`
		filename    string `json:"filename"`
		modPack     string `json:"modPack"`
	}
	err = json.Unmarshal(body, &modPackStruct)
	if err != nil {
		resp = fmt.Sprintf("Error unmarshalling modPack JSON: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	modPackMap, err := CreateNewModPackMap(w, &resp)
	if err != nil {
		return
	}

	if modPackMap.checkModPackExists(modPackStruct.modPack) {
		err = modPackMap[modPackStruct.modPack].Mods.updateMod(modPackStruct.modName, modPackStruct.downloadUrl, modPackStruct.filename)
	} else {
		err = errors.New("ModPack " + modPackStruct.modPack + " does not exist")
		resp = fmt.Sprintf("Error loading modpack file: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = fmt.Sprintf("Error updating mod {%s} in modpack {%s}: %s", modPackStruct.modName, modPackStruct.modPack, err)
		log.Println(resp)
		return
	}

	installedMods := modPackMap[modPackStruct.modPack].Mods.listInstalledMods().ModsResult
	for _, mod := range installedMods {
		if mod.Name == modPackStruct.modName {
			resp = mod
			break
		}
	}
}
