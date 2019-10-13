package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lockfile"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
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

// Returns JSON response of all mods installed in factorio/mods
func listInstalledModsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	mods, err := newMods(config.FactorioModsDir)

	if err != nil {
		w.WriteHeader(500)
		resp.Data = fmt.Sprintf("Error in ListInstalledMods handler: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in list mods: %s", err)
		}
		return
	}

	resp.Data = mods.listInstalledMods().ModsResult
	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in list mods: %s", err)
	}
}

// LoginFactorioModPortal returns JSON response with success or error-message
func LoginFactorioModPortal(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	username := r.FormValue("username")
	password := r.FormValue("password")

	loginStatus, err, statusCode := factorioLogin(username, password)
	if loginStatus == "" && err == nil {
		resp.Data = true
	}

	w.WriteHeader(statusCode)

	if err != nil {
		resp.Data = fmt.Sprintf("Error trying to login into Factorio: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in Factorio-Login: %s", err)
		}
		return
	}

	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in Factorio-Login: %s", err)
	}
}

func LoginstatusFactorioModPortal(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	var credentials FactorioCredentials
	resp.Data, err = credentials.load()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error getting the factorio credentials: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in Factorio-Login: %s", err)
		}
		return
	}

	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in Factorio-Login: %s", err)
	}
}

func LogoutFactorioModPortalHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	var credentials FactorioCredentials
	err = credentials.del()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error on logging out of factorio: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in Factorio-Login: %s", err)
		}
		return
	}

	resp.Data = false
	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in Factorio-Login: %s", err)
	}
}

//ModPortalSearchHandler returns JSON response with the found mods
func ModPortalSearchHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//Get Data out of the request
	searchKeyword := r.FormValue("search")

	var statusCode int
	resp.Data, err, statusCode = searchModPortal(searchKeyword)

	w.WriteHeader(statusCode)

	if err != nil {
		w.WriteHeader(500)
		resp.Data = fmt.Sprintf("Error in searchModPortal: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in searchModPortal: %s", err)
		}
		return
	}

	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in ModPortalSearch: %s", err)
	}
}

//ModPortalDetailsHandler returns JSON response with the mod details
func ModPortalDetailsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//Get Data out of the request
	modId := r.FormValue("modId")

	var statusCode int
	resp.Data, err, statusCode = getModDetails(modId)

	w.WriteHeader(statusCode)

	if err != nil {
		resp.Data = fmt.Sprintf("Error in searchModPortal: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in searchModPortal: %s", err)
		}
		return
	}

	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in ModPortalSearch: %s", err)
	}
}

func ModPortalInstallHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//Get Data out of the request
	downloadUrl := r.FormValue("link")
	filename := r.FormValue("filename")
	modName := r.FormValue("modName")

	mods, err := newMods(config.FactorioModsDir)
	if err == nil {
		err = mods.downloadMod(downloadUrl, filename, modName)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error in installMod: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in installMod: %s", err)
		}
		return
	}

	resp.Data = mods.listInstalledMods()
	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in ModPortalInstallHandler: %s", err)
	}
}

func ModPortalInstallMultipleHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	r.ParseForm()

	var modsList []string
	var versionsList []Version

	//Parse incoming data
	for key, values := range r.PostForm {
		if key == "mod_name" {
			for _, v := range values {
				modsList = append(modsList, v)
			}
		} else if key == "mod_version" {
			for _, value := range values {
				var v Version
				if err := v.UnmarshalText([]byte(value)); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					resp.Data = fmt.Sprintf("Error in searchModPortal: %s", err)
					if err := json.NewEncoder(w).Encode(resp); err != nil {
						log.Printf("Error in searchModPortal: %s", err)
					}
				}
				versionsList = append(versionsList, v)
			}
		}
	}

	mods, err := newMods(config.FactorioModsDir)
	if err != nil {
		log.Printf("error creating mods: %s", err)

		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error in searchModPortal: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in searchModPortal: %s", err)
		}
		return
	}

	for modIndex, mod := range modsList {
		var err error

		//get details of mod
		modDetails, err, statusCode := getModDetails(mod)
		if err != nil {
			w.WriteHeader(statusCode)
			resp.Data = fmt.Sprintf("Error in searchModPortal: %s", err)
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error in searchModPortal: %s", err)
			}
			return
		}

		modDetailsArray := []byte(modDetails)
		var modDetailsStruct ModPortalStruct

		//read mod-data into Struct
		err = json.Unmarshal(modDetailsArray, &modDetailsStruct)
		if err != nil {
			log.Printf("error reading modPortalDetails: %s", err)

			w.WriteHeader(http.StatusInternalServerError)
			resp.Data = fmt.Sprintf("Error in searchModPortal: %s", err)
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error in searchModPortal: %s", err)
			}
			return
		}

		//find correct mod-version
		for _, release := range modDetailsStruct.Releases {
			if release.Version.Equals(versionsList[modIndex]) {
				mods.downloadMod(release.DownloadURL, release.FileName, modDetailsStruct.Name)
				break
			}
		}
	}

	resp.Data = mods.listInstalledMods()

	resp.Success = true
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in ToggleModHandler: %s", err)
	}
}

func ToggleModHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//Get Data out of the request
	modName := r.FormValue("modName")

	mods, err := newMods(config.FactorioModsDir)
	if err == nil {
		err, resp.Data = mods.ModSimpleList.toggleMod(modName)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error in listInstalledModsByFolder: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in listInstalledModsByFolder: %s", err)
		}
		return
	}

	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in ToggleModHandler: %s", err)
	}
}

func DeleteModHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//Get Data out of the request
	modName := r.FormValue("modName")

	mods, err := newMods(config.FactorioModsDir)
	if err == nil {
		mods.deleteMod(modName)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error in deleteMod: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in DeleteModHandler: %s", err)
		}
		return
	}

	resp.Data = modName
	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in DeleteModHandler: %s", err)
	}
}

func DeleteAllModsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//delete mods folder
	err = deleteAllMods()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error in deleteMod: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in DeleteModHandler: %s", err)
		}
		return
	}

	resp.Data = nil
	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in DeleteModHandler: %s", err)
	}
}

func UpdateModHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//Get Data out of the request
	modName := r.FormValue("modName")
	downloadUrl := r.FormValue("downloadUrl")
	fileName := r.FormValue("filename")

	log.Println("--------------------------------------------------------------")

	mods, err := newMods(config.FactorioModsDir)
	if err == nil {
		err = mods.updateMod(modName, downloadUrl, fileName)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error in deleteMod: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in DeleteModHandler: %s", err)
		}
		return
	}

	installedMods := mods.listInstalledMods().ModsResult
	for _, mod := range installedMods {
		if mod.Name == modName {
			resp.Data = mod
			break
		}
	}
	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in DeleteModHandler: %s", err)
	}
}

func UploadModHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("Error in UploadModHandler: %s", err)
	}
}

func DownloadModsHandler(w http.ResponseWriter, r *http.Request) {
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
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//Get Data out of the request
	SaveFile := r.FormValue("saveFile")

	path := filepath.Join(config.FactorioSavesDir, SaveFile)
	f, err := OpenArchiveFile(path, "level.dat")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("cannot open save level file: %v", err)
		resp.Data = "Error opening save file"
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in loadModsFromSave: %s", err)
		}
		return
	}
	defer f.Close()

	var header SaveHeader
	err = header.ReadFrom(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("cannot read save header: %v", err)
		resp.Data = "Error reading save file"
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in loadModsFromSave: %s", err)
		}
		return
	}

	resp.Data = header
	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in LoadModsFromSave: %s", err)
	}
}

func ListModPacksHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	modPackMap, err := newModPackMap()

	if err != nil {
		w.WriteHeader(500)
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error listing modpack files: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing modpacks: %s", err)
		}
		return
	}

	resp.Data = modPackMap.listInstalledModPacks()
	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error listing saves: %s", err)
	}
}

func CreateModPackHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	name := r.FormValue("name")

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	modPackMap, err := newModPackMap()
	if err == nil {
		err = modPackMap.createModPack(name)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error creating modpack file: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error creating modpack: %s", err)
		}
		return
	}

	resp.Data = modPackMap.listInstalledModPacks()
	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error creating modpack response: %s", err)
	}
}

func DownloadModPackHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	vars := mux.Vars(r)
	modpack := vars["modpack"]

	modPackMap, err := newModPackMap()
	if err != nil {
		log.Printf("error on loading modPacks: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if modPackMap.checkModPackExists(modpack) {
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
			log.Printf("error on walking over the modpack: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("requested modPack doesnt exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	writerHeader := w.Header()
	writerHeader.Set("Content-Type", "application/zip;charset=UTF-8")
	writerHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", modpack+".zip"))
}

func DeleteModPackHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	name := r.FormValue("name")

	modPackMap, err := newModPackMap()
	if err == nil {
		err = modPackMap.deleteModPack(name)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error deleting modpack file: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error deleting modpack: %s", err)
		}
		return
	}

	resp.Data = name
	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error creating delete modpack response: %s", err)
	}
}

func LoadModPackHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	name := r.FormValue("name")

	modPackMap, err := newModPackMap()
	if err == nil {
		modPackMap[name].loadModPack()
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error loading modpack file: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error loading modpack: %s", err)
		}
		return
	}

	resp.Data = modPackMap[name].Mods.listInstalledMods()
	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error creating loading modpack response: %s", err)
	}
}

func ModPackToggleModHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	modName := r.FormValue("modName")
	modPackName := r.FormValue("modPack")

	modPackMap, err := newModPackMap()
	if err == nil {
		err, resp.Data = modPackMap[modPackName].Mods.ModSimpleList.toggleMod(modName)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error loading modpack file: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error loading modpack: %s", err)
		}
		return
	}

	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error creating loading modpack response: %s", err)
	}
}

func ModPackDeleteModHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	modName := r.FormValue("modName")
	modPackName := r.FormValue("modPackName")

	modPackMap, err := newModPackMap()
	if err == nil {
		if modPackMap.checkModPackExists(modPackName) {
			err = modPackMap[modPackName].Mods.deleteMod(modName)
		} else {
			err = errors.New("ModPack " + modPackName + " does not exist")
		}
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error loading modpack file: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error loading modpack: %s", err)
		}
		return
	}

	resp.Data = true
	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error creating loading modpack response: %s", err)
	}

	return
}

func ModPackUpdateModHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//Get Data out of the request
	modName := r.FormValue("modName")
	downloadUrl := r.FormValue("downloadUrl")
	fileName := r.FormValue("filename")
	modPackName := r.FormValue("modPackName")

	modPackMap, err := newModPackMap()
	if err == nil {
		if modPackMap.checkModPackExists(modPackName) {
			err = modPackMap[modPackName].Mods.updateMod(modName, downloadUrl, fileName)
		} else {
			err = errors.New("ModPack " + modPackName + "does not exist")
		}
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error in deleteMod: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in DeleteModHandler: %s", err)
		}
		return
	}

	installedMods := modPackMap[modPackName].Mods.listInstalledMods().ModsResult
	for _, mod := range installedMods {
		if mod.Name == modName {
			resp.Data = mod
			break
		}
	}
	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in DeleteModHandler: %s", err)
	}
}
