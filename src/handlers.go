package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"github.com/gorilla/mux"
	"archive/zip"
	"errors"
)

type JSONResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,string"`
}
type JSONResponseFileInput struct {
	Success	bool		`json:"success"`
	Data	interface{}	`json:"data,string"`
	Error	string		`json:"error"`
	ErrorKeys	[]int	`json:"errorkeys"`
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

// Returns JSON response with the userKey or the error-message
func LoginFactorioModPortal(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	username := r.FormValue("username")
	password := r.FormValue("password")

    var statusCode int
	resp.Data, err, statusCode = getUserToken(username, password)

    w.WriteHeader(statusCode)

	if err != nil {
		w.WriteHeader(500)
		resp.Data = fmt.Sprintf("Error in getUserToken or LoginFactorioModPortal handler: %s", err)
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

//Returns JSON response with the found mods
func ModPortalSearchHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    resp := JSONResponse{
        Success: false,
    }

    w.Header().Set("Content-Type", "application/json;charset=UTF-8")

    //Get Data out of the request
    search_keyword := r.FormValue("search")

    var statusCode int
    resp.Data, err, statusCode = searchModPortal(search_keyword)

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

//Returns JSON response with the mod details
func ModPortalDetailsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//Get Data out of the request
	mod_id := r.FormValue("mod_id")

	var statusCode int
	resp.Data, err, statusCode = getModDetails(mod_id)

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

func ModPortalInstallHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//Get Data out of the request
	username := r.FormValue("username")
	userKey := r.FormValue("userKey")
	downloadUrl := r.FormValue("link")
	filename := r.FormValue("filename")
	mod_name := r.FormValue("modName")

	mods, err := newMods(config.FactorioModsDir)
	if err == nil {
		err = mods.downloadMod(username, userKey, downloadUrl, filename, mod_name)
	}

	if err != nil {
		w.WriteHeader(500)
		resp.Data = fmt.Sprintf("Error in installMod: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in installMod: %s", err)
		}
		return
	}

	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in ModPortalInstallHandler: %s", err)
	}
}

func ToggleModHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//Get Data out of the request
	mod_name := r.FormValue("mod_name")

	mods, err := newMods(config.FactorioModsDir)
	if err == nil {
		err = mods.ModSimpleList.toggleMod(mod_name)
	}

	if err != nil {
		w.WriteHeader(500)
		resp.Data = fmt.Sprintf("Error in listInstalledModsByFolder: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in listInstalledModsByFolder: %s", err)
		}
		return
	}

	resp.Data = mods.listInstalledMods().ModsResult
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
	mod_name := r.FormValue("mod_name")

	mods, err := newMods(config.FactorioModsDir)
	if err == nil {
		mods.deleteMod(mod_name)
	}

	if err != nil {
		w.WriteHeader(500)
		resp.Data = fmt.Sprintf("Error in deleteMod: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in DeleteModHandler: %s", err)
		}
		return
	}

	resp.Data = mods.listInstalledMods().ModsResult
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
	mod_name := r.FormValue("mod_name")
	username := r.FormValue("username")
	user_key := r.FormValue("userKey")
	download_url := r.FormValue("downloadUrl")
	file_name := r.FormValue("filename")

	mods, err := newMods(config.FactorioModsDir)
	if err == nil {
		err = mods.updateMod(mod_name, username, user_key, download_url, file_name)
	}

	if err != nil {
		w.WriteHeader(500)
		resp.Data = fmt.Sprintf("Error in deleteMod: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in DeleteModHandler: %s", err)
		}
		return
	}

	resp.Data = mods.listInstalledMods()
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
		for file_key, mod_file := range r.MultipartForm.File["mod_file"] {
			err = mods.uploadMod(mod_file)
			if err != nil {
				resp.ErrorKeys = append(resp.ErrorKeys, file_key)
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

func ListModPacksHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	mod_pack_map, err := newModPackMap()

	if err != nil {
		w.WriteHeader(500)
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error listing modpack files: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing modpacks: %s", err)
		}
		return
	}

	resp.Data = mod_pack_map.listInstalledModPacks()
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

	mod_pack_map, err := newModPackMap()
	if err == nil {
		err = mod_pack_map.createModPack(name)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error creating modpack file: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error creating modpack: %s", err)
		}
		return
	}

	resp.Data = mod_pack_map.listInstalledModPacks()
	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error creating modpack response: %s", err)
	}
}

func DownloadModPackHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	vars := mux.Vars(r)
	modpack := vars["modpack"]

	mod_pack_map, err := newModPackMap()
	if err != nil {
		log.Printf("error on loading modPacks: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if mod_pack_map.checkModPackExists(modpack) {
		zip_writer := zip.NewWriter(w)
		defer zip_writer.Close()

		//iterate over folder and create everything in the zip
		err = filepath.Walk(filepath.Join(config.FactorioModPackDir, modpack), func(path string, info os.FileInfo, err error) error {
			if info.IsDir() == false {
				writer, err := zip_writer.Create(info.Name())
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

	writer_header := w.Header()
	writer_header.Set("Content-Type", "application/zip;charset=UTF-8")
	writer_header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", modpack + ".zip"))
}

func DeleteModPackHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	name := r.FormValue("name")

	mod_pack_map, err := newModPackMap()
	if err == nil {
		err = mod_pack_map.deleteModPack(name)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error deleting modpack file: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error deleting modpack: %s", err)
		}
		return
	}

	resp.Data = mod_pack_map.listInstalledModPacks()
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

	mod_pack_map, err := newModPackMap()
	if err == nil {
		mod_pack_map[name].loadModPack()
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Data = fmt.Sprintf("Error loading modpack file: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error loading modpack: %s", err)
		}
		return
	}

	resp.Data = mod_pack_map[name].Mods.listInstalledMods()
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

    mod_name := r.FormValue("mod_name")
    mod_pack_name := r.FormValue("mod_pack")

    mod_pack_map, err := newModPackMap()
	if err == nil {
		err = mod_pack_map[mod_pack_name].Mods.ModSimpleList.toggleMod(mod_name)
	}
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        resp.Data = fmt.Sprintf("Error loading modpack file: %s", err)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error loading modpack: %s", err)
        }
        return
    }

    resp.Data = mod_pack_map.listInstalledModPacks()
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

    mod_name := r.FormValue("mod_name")
    mod_pack_name := r.FormValue("mod_pack_name")

    mod_pack_map, err := newModPackMap()
    if err == nil {
    	if mod_pack_map.checkModPackExists(mod_pack_name) {
			err = mod_pack_map[mod_pack_name].Mods.deleteMod(mod_name)
		} else {
			err = errors.New("ModPack " + mod_pack_name + " does not exist")
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

    resp.Data = mod_pack_map.listInstalledModPacks()
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
	mod_name := r.FormValue("mod_name")
	username := r.FormValue("username")
	user_key := r.FormValue("userKey")
	download_url := r.FormValue("downloadUrl")
	file_name := r.FormValue("filename")
	mod_pack_name := r.FormValue("mod_pack_name")

	mod_pack_map, err := newModPackMap()
	if err == nil {
		if mod_pack_map.checkModPackExists(mod_pack_name){
			err = mod_pack_map[mod_pack_name].Mods.updateMod(mod_name, username, user_key, download_url, file_name)
		} else {
			err = errors.New("ModPack " + mod_pack_name + "does not exist")
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

	resp.Data = mod_pack_map.listInstalledModPacks()
	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in DeleteModHandler: %s", err)
	}
}

// Lists all save files in the factorio/saves directory
func ListSaves(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	savesList, err := listSaves(config.FactorioSavesDir)
	if err != nil {
		resp.Success = false
		resp.Data = fmt.Sprintf("Error listing save files: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing saves: %s", err)
		}
		return
	}

	loadLatest := Save{Name: "Load Latest"}
	savesList = append(savesList, loadLatest)

	resp.Data = savesList

	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error listing saves: %s", err)
	}
}

func DLSave(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")

	vars := mux.Vars(r)
	save := vars["save"]
	saveName := config.FactorioSavesDir + "/" + save

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", save))
	log.Printf("%s downloading: %s", r.Host, saveName)

	http.ServeFile(w, r, saveName)
}

func UploadSave(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	switch r.Method {
	case "GET":
		resp.Data = "Unsupported method"
		resp.Success = false
		if err = json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing mods: %s", err)
		}
	case "POST":
		log.Println("Uploading save file")
		r.ParseMultipartForm(32 << 20)
		file, header, err := r.FormFile("savefile")
		if err != nil {
			resp.Success = false
			resp.Data = err.Error()
			json.NewEncoder(w).Encode(resp)
			log.Printf("Error in upload save formfile: %s", err.Error())
			return
		}
		defer file.Close()

		out, err := os.Create(config.FactorioSavesDir + "/" + header.Filename)
		if err != nil {
			resp.Success = false
			resp.Data = err.Error()
			json.NewEncoder(w).Encode(resp)
			log.Printf("Error in out: %s", err)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			resp.Success = false
			resp.Data = err.Error()
			json.NewEncoder(w).Encode(resp)
			log.Printf("Error in io copy: %s", err)
			return
		}
		log.Printf("Uploaded save file: %s", header.Filename)
		resp.Data = "File '" + header.Filename + "' uploaded successfully"
		resp.Success = true
		json.NewEncoder(w).Encode(resp)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Deletes provided save
func RemoveSave(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	vars := mux.Vars(r)
	name := vars["save"]

	save, err := findSave(name)
	if err != nil {
		resp.Data = fmt.Sprintf("Error removing save: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error removing save %s", err)
		}
	}

	err = save.remove()
	if err == nil {
		// save was removed
		resp.Data = fmt.Sprintf("Removed save: %s", save.Name)
		resp.Success = true
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error removing save %s", err)
		}
	} else {
		log.Printf("Error in remove save handler: %s", err)
		resp.Data = fmt.Sprintf("Error in remove save handler: %s", err)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error removing save: %s", err)
		}
	}
}

// Launches Factorio server binary with --create flag to create save
// Url must include save name for creation of savefile
func CreateSaveHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	vars := mux.Vars(r)
	saveName := vars["save"]

	if saveName == "" {
		log.Printf("Error creating save, no name provided: %s", err)
		resp.Data = "No save name provided."
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding save handler response: %s", err)
		}
		return
	}

	saveFile := filepath.Join(config.FactorioSavesDir, saveName)
	cmdOut, err := createSave(saveFile)
	if err != nil {
		log.Printf("Error creating save: %s", err)
		resp.Data = "Error creating savefile."
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding save handler response: %s", err)
		}
		return
	}

	resp.Success = true
	resp.Data = fmt.Sprintf("Save %s created successfully. Command output: \n%s", saveName, cmdOut)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding save response: %s", err)
	}
}

// Returns last lines of the factorio-current.log file
func LogTail(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	resp.Data, err = tailLog(config.FactorioLog)
	if err != nil {
		resp.Data = fmt.Sprintf("Could not tail %s: %s", config.FactorioLog, err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Could not tail %s: %s", config.FactorioLog, err)
		}
		return
	}

	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error tailing logfile: %s", err)
	}
}

// Return JSON response of config.ini file
func LoadConfig(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	configContents, err := loadConfig(config.FactorioConfigFile)
	if err != nil {
		log.Printf("Could not retrieve config.ini: %s", err)
		resp.Data = "Error getting config.ini"
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error tailing logfile: %s", err)
		}
	} else {
		resp.Data = configContents
		resp.Success = true
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding config file JSON reponse: %s", err)
	}

	log.Printf("Sent config.ini response")
}

func StartServer(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if FactorioServ.Running {
		resp.Data = "Factorio server is already running"
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding JSON response: %s", err)
		}
		return
	}

	switch r.Method {
	case "GET":
		log.Printf("GET not supported for startserver handler")
		resp.Data = "Unsupported method"
		resp.Success = false
		if err = json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing mods: %s", err)
		}
	case "POST":
		log.Printf("Starting Factorio server.")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error in starting factorio server handler body: %s", err)
			return
		}

		log.Printf("Starting Factorio server with settings: %v", string(body))

		err = json.Unmarshal(body, &FactorioServ)
		if err != nil {
			log.Printf("Error unmarshaling server settings JSON: %s", err)
			return
		}

		// Check if savefile was submitted with request to start server.
		if FactorioServ.Savefile == "" {
			log.Printf("Error starting Factorio server: no save file provided")
			resp.Success = false
			resp.Data = fmt.Sprintf("Error starting Factorio server: %s", "No save file provided")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error encoding config file JSON reponse: %s", err)
			}
			return
		}

		go func() {
			err = FactorioServ.Run()
			if err != nil {
				log.Printf("Error starting Factorio server: %+v", err)
				return
			}
		}()

		timeout := 0
		for timeout <= 3 {
			time.Sleep(1 * time.Second)
			if FactorioServ.Running {
				resp.Data = fmt.Sprintf("Factorio server with save: %s started on port: %d", FactorioServ.Savefile, FactorioServ.Port)
				resp.Success = true
				log.Printf("Factorio server started on port: %v", FactorioServ.Port)
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					log.Printf("Error encoding config file JSON reponse: %s", err)
				}
				break
			} else {
				log.Printf("Did not detect running Factorio server attempt: %+v", timeout)
			}

			timeout++
		}
		if FactorioServ.Running == false {
			log.Printf("Error starting Factorio server: %s", err)
			resp.Data = fmt.Sprintf("Error starting Factorio server: %s", err)
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error encoding start server JSON response: %s", err)
			}
		}
	}
}

func StopServer(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if FactorioServ.Running {
		err := FactorioServ.Stop()
		if err != nil {
			log.Printf("Error in stop server handler: %s", err)
			resp.Data = fmt.Sprintf("Error in stop server handler: %s", err)
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error encoding config file JSON reponse: %s", err)
			}
			return
		}

		log.Printf("Stopped Factorio server.")
		resp.Success = true
		resp.Data = fmt.Sprintf("Factorio server stopped")
	} else {
		resp.Data = "Factorio server is not running"
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding config file JSON reponse: %s", err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding config file JSON reponse: %s", err)
	}
}


func KillServer(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if FactorioServ.Running {
		err := FactorioServ.Kill()
		if err != nil {
			log.Printf("Error in kill server handler: %s", err)
			resp.Data = fmt.Sprintf("Error in kill server handler: %s", err)
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error encoding config file JSON reponse: %s", err)
			}
			return
		}

		log.Printf("Killed Factorio server.")
		resp.Success = true
		resp.Data = fmt.Sprintf("Factorio server killed")
	} else {
		resp.Data = "Factorio server is not running"
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding config file JSON reponse: %s", err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding config file JSON reponse: %s", err)
	}
}

func CheckServer(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if FactorioServ.Running {
		resp.Success = true
		status := map[string]string{}
		status["status"] = "running"
		status["port"] = strconv.Itoa(FactorioServ.Port)
		status["savefile"] = FactorioServ.Savefile
		resp.Data = status
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding config file JSON reponse: %s", err)
		}
	} else {
		resp.Success = true
		status := map[string]string{}
		status["status"] = "stopped"
		resp.Data = status
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding config file JSON reponse: %s", err)
		}
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	switch r.Method {
	case "GET":
		log.Printf("GET not supported for login handler")
		resp.Data = "Unsupported method"
		resp.Success = false
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing mods: %s", err)
		}
	case "POST":
		var user User
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error in starting factorio server handler body: %s", err)
			return
		}

		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Printf("Error unmarshaling server settings JSON: %s", err)
			return
		}

		log.Printf("Logging in user: %s", user.Username)

		err = Auth.aaa.Login(w, r, user.Username, user.Password, "/")
		if err != nil {
			log.Printf("Error logging in user: %s, error: %s", user.Username, err)
			resp.Data = fmt.Sprintf("Error logging in user: %s", user.Username)
			resp.Success = false
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error listing mods: %s", err)
			}
			return
		}

		log.Printf("User: %s, logged in successfully", user.Username)
		resp.Data = fmt.Sprintf("User: %s, logged in successfully", user.Username)
		resp.Success = true
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing mods: %s", err)
		}
	}
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if err := Auth.aaa.Logout(w, r); err != nil {
		log.Printf("Error logging out current user")
		return
	}

	resp.Success = true
	resp.Data = "User logged out successfully."
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error logging out: %s", err)
	}
}

func GetCurrentLogin(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	user, err := Auth.aaa.CurrentUser(w, r)
	if err != nil {
		log.Printf("Error getting current user status: %s", err)
		resp.Data = fmt.Sprintf("Error getting user status: %s", user.Username)
		resp.Success = false
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing mods: %s", err)
		}
		return
	}

	resp.Success = true
	resp.Data = user

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error getting user status: %s", err)
	}

}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	users, err := Auth.listUsers()
	if err != nil {
		log.Printf("Error in ListUsers handler: %s", err)
		resp.Data = fmt.Sprint("Error listing users")
		resp.Success = false
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing mods: %s", err)
		}
		return
	}

	resp.Success = true
	resp.Data = users

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error getting user status: %s", err)
	}
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	switch r.Method {
	case "GET":
		log.Printf("GET not supported for add user handler")
		resp.Data = "Unsupported method"
		resp.Success = false
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error adding user: %s", err)
		}
	case "POST":
		user := User{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error in reading add user POST: %s", err)
			resp.Data = fmt.Sprintf("Error in adding user: %s", err)
			resp.Success = false
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error adding user: %s", err)
			}
			return
		}

		log.Printf("Adding user: %v", string(body))

		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Printf("Error unmarshaling user add JSON: %s", err)
			resp.Data = fmt.Sprintf("Error in adding user: %s", err)
			resp.Success = false
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error adding user: %s", err)
			}
			return
		}

		err = Auth.addUser(user.Username, user.Password, user.Email, user.Role)
		if err != nil {
			log.Printf("Error in adding user: %s", err)
			resp.Data = fmt.Sprintf("Error in adding user: %s", err)
			resp.Success = false
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error adding user: %s", err)
			}
			return
		}

		resp.Success = true
		resp.Data = fmt.Sprintf("User: %s successfully added.", user.Username)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in returning added user response: %s", err)
		}
	}
}

func RemoveUser(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	switch r.Method {
	case "GET":
		log.Printf("GET not supported for add user handler")
		resp.Data = "Unsupported method"
		resp.Success = false
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error adding user: %s", err)
		}
	case "POST":
		user := User{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error in reading remove user POST: %s", err)
			resp.Data = fmt.Sprintf("Error in removing user: %s", err)
			resp.Success = false
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error adding user: %s", err)
			}
			return
		}
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Printf("Error unmarshaling user remove JSON: %s", err)
			resp.Data = fmt.Sprintf("Error in removing user: %s", err)
			resp.Success = false
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error removing user: %s", err)
			}
			return
		}

		err = Auth.removeUser(user.Username)
		if err != nil {
			log.Printf("Error in remove user handler: %s", err)
			resp.Data = fmt.Sprintf("Error in removing user: %s", err)
			resp.Success = false
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error adding user: %s", err)
			}
			return
		}

		resp.Success = true
		resp.Data = fmt.Sprintf("User: %s successfully removed.", user.Username)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in returning remove user response: %s", err)
		}
	}
}

// Return JSON response of server-settings.json file
func GetServerSettings(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	resp.Data = FactorioServ.Settings
	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding server settings JSON reponse: %s", err)
	}

	log.Printf("Sent server settings response")
}

func UpdateServerSettings(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	switch r.Method {
	case "GET":
		log.Printf("GET not supported for add user handler")
		resp.Data = "Unsupported method"
		resp.Success = false
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error adding user: %s", err)
		}
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error in reading server settings POST: %s", err)
			resp.Data = fmt.Sprintf("Error in updating settings: %s", err)
			resp.Success = false
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error updating settings: %s", err)
			}
			return
		}
		log.Printf("Received settings JSON: %s", body)

		err = json.Unmarshal(body, &FactorioServ.Settings)
		if err != nil {
			log.Printf("Error unmarshaling server settings JSON: %s", err)
			resp.Data = fmt.Sprintf("Error in updating settings: %s", err)
			resp.Success = false
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error encoding server settings response: %s", err)
			}
			return
		}

		settings, err := json.MarshalIndent(&FactorioServ.Settings, "", "  ")
		if err != nil {
			log.Printf("Failed to marshal server settings: %s", err)
			return
		} else {
			if err = ioutil.WriteFile(filepath.Join(config.FactorioConfigDir, config.SettingsFile), settings, 0644); err != nil {
				log.Printf("Failed to save server settings: %v\n", err)
				return
			}
			log.Printf("Saved Factorio server settings in server-settings.json")
		}

		resp.Success = true
		resp.Data = fmt.Sprintf("Settings successfully saved: %s", &FactorioServ.Settings)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in sending server settings response: %s", err)
		}
	}
}
