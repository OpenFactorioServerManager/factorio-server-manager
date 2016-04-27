package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type JSONResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,string"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

// Returns JSON response of all mods installed in factorio/mods
func ListInstalledMods(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	modDir := config.FactorioDir + "/mods"

	resp.Data, err = listInstalledMods(modDir)
	if err != nil {
		resp.Data = fmt.Sprintf("Error in ListInstalledMods handler: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in list mods: %s", err)
		}
		return
	}

	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in list mods: %s", err)
	}
}

// Toggles mod passed in through mod variable
// Updates mod-list.json file to toggle the enabled status of mods
func ToggleMod(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	vars := mux.Vars(r)
	modName := vars["mod"]

	m, err := parseModList()
	if err != nil {
		resp.Data = fmt.Sprintf("Could not parse mod list: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in list mods: %s", err)
		}
		return
	}

	err = m.toggleMod(modName)
	if err != nil {
		resp.Data = fmt.Sprintf("Could not toggle mod: %s error: %s", modName, err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in list mods: %s", err)
		}
		return
	}

	resp.Success = true
	resp.Data = m

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in toggle mod: %s", err)
	}
}

// Returns JSON response of all mods in the mod-list.json file
func ListMods(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	resp.Data, err = parseModList()
	if err != nil {
		resp.Data = fmt.Sprintf("Could not parse mod list: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in list mods: %s", err)
		}
		return
	}

	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error listing mods: %s", err)
	}
}

// Uploads mod to the mods directory
func UploadMod(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	switch r.Method {
	case "GET":
		resp.Data = "Unsupported method"
		if err = json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing mods: %s", err)
		}
	case "POST":
		log.Println("Uploading file")
		r.ParseMultipartForm(32 << 20)
		file, header, err := r.FormFile("modfile")
		if err != nil {
			log.Printf("No mod filename provided for upload: %s", err)
			json.NewEncoder(w).Encode("No mod file provided.")
			return
		}
		defer file.Close()

		out, err := os.Create(config.FactorioModsDir + "/" + header.Filename)
		if err != nil {
			resp.Data = err.Error()
			json.NewEncoder(w).Encode(resp)
			log.Printf("Error in out")
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			resp.Data = err.Error()
			json.NewEncoder(w).Encode(resp)
			log.Printf("Error in io copy")
			return
		}
		log.Printf("Uploaded mod file: %s", header.Filename)
		resp.Data = "File '" + header.Filename + "' submitted successfully"
		resp.Success = true
		json.NewEncoder(w).Encode(resp)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func RemoveMod(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	vars := mux.Vars(r)
	modName := vars["mod"]

	err = rmMod(modName)
	if err == nil {
		// No error returned means mod was removed
		resp.Data = fmt.Sprintf("Removed mod: %s", modName)
		resp.Success = true
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error removing mod: %s", err)
		}
	} else {
		log.Printf("Error in remove mod handler: %s", err)
		resp.Data = fmt.Sprintf("Error in remove mod handler: %s", err)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error removing mod: %s", err)
		}
		return
	}
}

func DownloadMod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	vars := mux.Vars(r)
	mod := vars["mod"]
	modFile := config.FactorioModsDir + "/" + mod

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", mod))
	log.Printf("%s downloading: %s", r.Host, modFile)

	http.ServeFile(w, r, modFile)
}

// Lists all save files in the factorio/saves directory
func ListSaves(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	saveDir := config.FactorioDir + "/saves"

	resp.Data, err = listSaves(saveDir)
	if err != nil {
		resp.Data = fmt.Sprintf("Error listing save files: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing saves: %s", err)
		}
		return
	}

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
	saveName := vars["save"]

	err = rmSave(saveName)
	if err == nil {
		// save was removed
		resp.Data = fmt.Sprintf("Removed save: %s", saveName)
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

	cmdOut, err := createSave(saveName)
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
		log.Printf("Error tailing logfile", err)
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
			log.Printf("Error tailing logfile", err)
		}
		return
	}

	resp.Data = configContents
	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding config file JSON reponse: ", err)
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
			log.Printf("Error encoding JSON response: ", err)
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

		// TODO get form parameters for starting server

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

		go func() {
			err = FactorioServ.Run()
			if err != nil {
				log.Printf("Error starting Factorio server: %s", err)
				resp.Data = fmt.Sprintf("Error starting Factorio server: %s", err)
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					log.Printf("Error encoding config file JSON reponse: ", err)
				}
				return
			}
		}()

		if FactorioServ.Running {
			log.Printf("Factorio server started on port: %d", FactorioServ.Port)
		}

		resp.Data = fmt.Sprintf("Factorio server started on port: %s", FactorioServ.Port)
		resp.Success = true

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding config file JSON reponse: ", err)
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
				log.Printf("Error encoding config file JSON reponse: ", err)
			}
			return
		}
		log.Printf("Stopped Factorio server.")
		resp.Success = true
		resp.Data = fmt.Sprintf("Factorio server stopped")
	} else {
		resp.Data = "Factorio server is not running"
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding config file JSON reponse: ", err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding config file JSON reponse: ", err)
	}
}

func RunningServer(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if FactorioServ.Running {
		log.Printf("Creating server status response")
		resp.Success = true
		status := map[string]string{}
		status["status"] = "running"
		status["port"] = strconv.Itoa(FactorioServ.Port)
		status["savefile"] = FactorioServ.Savefile
		resp.Data = status
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding config file JSON reponse: ", err)
		}
		log.Printf("Server status sent with data: %+v", resp.Data)
	} else {
		log.Printf("Server not running, creating status response")
		resp.Success = true
		status := map[string]string{}
		status["status"] = "stopped"
		resp.Data = status
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding config file JSON reponse: ", err)
		}
	}
}
