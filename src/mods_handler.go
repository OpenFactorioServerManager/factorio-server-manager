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

func CreateNewMods(w http.ResponseWriter, resp *interface{}) (mods Mods, err error) {
	mods, err = newMods(config.FactorioModsDir)
	if err != nil {
		*resp = fmt.Sprintf("Error creating mods object: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
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

// Returns JSON response of all mods installed in factorio/mods
func listInstalledModsHandler(w http.ResponseWriter, r *http.Request) {
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

// LoginFactorioModPortal returns JSON response with success or error-message
func LoginFactorioModPortal(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	username := r.FormValue("username")
	password := r.FormValue("password")

	loginStatus, err, statusCode := factorioLogin(username, password)
	if err != nil {
		resp = fmt.Sprintf("Error trying to login into Factorio: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if loginStatus == "" {
		resp = true
	}

	w.WriteHeader(statusCode)
}

func LoginstatusFactorioModPortal(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	var credentials FactorioCredentials
	resp, err = credentials.load()

	if err != nil {
		resp = fmt.Sprintf("Error getting the factorio credentials: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func LogoutFactorioModPortalHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	var credentials FactorioCredentials
	err = credentials.del()

	if err != nil {
		resp = fmt.Sprintf("Error on logging out of factorio: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp = false
}

//ModPortalSearchHandler returns JSON response with the found mods
func ModPortalSearchHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	// Get Data out of the request
	searchKeyword := r.FormValue("search")

	var statusCode int
	resp, err, statusCode = searchModPortal(searchKeyword)

	if err != nil {
		resp = fmt.Sprintf("Error in searching in mod portal: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
}

//ModPortalDetailsHandler returns JSON response with the mod details
func ModPortalDetailsHandler(w http.ResponseWriter, r *http.Request) {
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

	var statusCode int
	resp, err, statusCode = getModDetails(mod.Name)

	if err != nil {
		resp = fmt.Sprintf("Error getting mod details: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
}

func ModPortalInstallHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//Get Data out of the request
	downloadUrl := r.FormValue("link")
	filename := r.FormValue("filename")
	modName := r.FormValue("modName")

	mods, err := CreateNewMods(w, &resp)
	if err != nil {
		return
	}

	err = mods.downloadMod(downloadUrl, filename, modName)
	if err != nil {
		resp = fmt.Sprintf("Error downloading a mod: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp = mods.listInstalledMods()
}

func ModPortalInstallMultipleHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

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
					resp = fmt.Sprintf("Error loading version from uploaded form: %s", err)
					log.Println(resp)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				versionsList = append(versionsList, v)
			}
		}
	}

	mods, err := CreateNewMods(w, &resp)
	if err != nil {
		return
	}

	for modIndex, mod := range modsList {
		var err error

		//get details of mod
		modDetails, err, statusCode := getModDetails(mod)
		if err != nil {
			resp = fmt.Sprintf("Error getting mod details of mod {%s}: %s", mod, err)
			log.Println(resp)
			w.WriteHeader(statusCode)
			return
		}

		modDetailsArray := []byte(modDetails)
		var modDetailsStruct ModPortalStruct

		//read mod-data into Struct
		err = json.Unmarshal(modDetailsArray, &modDetailsStruct)
		if err != nil {
			resp = fmt.Sprintf("error unmarshalling mod details: %s", err)
			log.Println(resp)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//find correct mod-version
		for _, release := range modDetailsStruct.Releases {
			if release.Version.Equals(versionsList[modIndex]) {
				err = mods.downloadMod(release.DownloadURL, release.FileName, modDetailsStruct.Name)
				if err != nil {
					resp = fmt.Sprintf("Error downloading mod {%s}, error: %s", modDetailsStruct.Name, err)
					log.Println(resp)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				break
			}
		}
	}

	resp = mods.listInstalledMods()
}

func ToggleModHandler(w http.ResponseWriter, r *http.Request) {
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

func DeleteModHandler(w http.ResponseWriter, r *http.Request) {
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

func DeleteAllModsHandler(w http.ResponseWriter, r *http.Request) {
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

func UpdateModHandler(w http.ResponseWriter, r *http.Request) {
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

func ListModPacksHandler(w http.ResponseWriter, r *http.Request) {
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

func CreateModPackHandler(w http.ResponseWriter, r *http.Request) {
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

	err = modPackMap.deleteModPack(modPackStruct.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = fmt.Sprintf("Error deleting modpack file: %s", err)
		log.Println(resp)
		return
	}

	resp = modPackStruct.Name
}

func LoadModPackHandler(w http.ResponseWriter, r *http.Request) {
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

	modPackMap[modPackStruct.Name].loadModPack()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = fmt.Sprintf("Error loading modpack file: %s", err)
		log.Println(resp)
		return
	}

	resp = modPackMap[modPackStruct.Name].Mods.listInstalledMods()
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
