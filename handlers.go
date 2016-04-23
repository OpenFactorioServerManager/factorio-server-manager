package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

// Returns JSON response of all mods installed in factorio/mods
func ListInstalledMods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	modDir := config.FactorioDir + "/mods"

	mods, err := listInstalledMods(modDir)
	if err != nil {
		log.Printf("Error in ListInstalledMods handler: %s", err)
		return
	}

	if err := json.NewEncoder(w).Encode(mods); err != nil {
		log.Printf("Error in list mods: %s", err)
	}
}

// Toggles mod passed in through mod variable
// Updates mod-list.json file to toggle the enabled status of mods
func ToggleMod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	vars := mux.Vars(r)
	modName := vars["mod"]

	m, err := parseModList()
	if err != nil {
		log.Printf("Could not parse mod list: %s", err)
		return
	}

	err = m.toggleMod(modName)
	if err != nil {
		log.Printf("Could not toggle mod: %s error: %s", modName, err)
		return
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("Error in toggle mod: %s", err)
	}
}

// Returns JSON response of all mods in the mod-list.json file
func ListMods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	m, err := parseModList()
	if err != nil {
		log.Printf("Could not parse mod list: %s", err)
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("Error listing mods: %s", err)
	}
}

// Uploads mod to the mods directory
func UploadMod(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		resp := "Unsupported method"
		if err := json.NewEncoder(w).Encode(resp); err != nil {
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
			json.NewEncoder(w).Encode(err.Error())
			log.Printf("Error in out")
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
			log.Printf("Error in io copy")
			return
		}
		log.Printf("Uploaded mod file: %s", header.Filename)
		resp := "File '" + header.Filename + "' submitted successfully"
		json.NewEncoder(w).Encode(resp)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func RemoveMod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	vars := mux.Vars(r)
	modName := vars["mod"]

	err := rmMod(modName)
	if err == nil {
		// No error returned means mod was removed
		resp := fmt.Sprintf("Removed mod: %s", modName)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error removing mod: %s", err)
		}
	} else {
		log.Printf("Error in remove mod handler: %s", err)
		resp := fmt.Sprintf("Error in remove mod handler: %s", err)

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
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	saveDir := config.FactorioDir + "/saves"

	saves, err := listSaves(saveDir)
	if err != nil {
		log.Printf("Error in ListSaves handler: %s", err)
		return
	}

	if err := json.NewEncoder(w).Encode(saves); err != nil {
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
	switch r.Method {
	case "GET":
		resp := "Unsupported method"
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing mods: %s", err)
		}
	case "POST":
		log.Println("Uploading save file")
		r.ParseMultipartForm(32 << 20)
		file, header, err := r.FormFile("savefile")
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
			log.Printf("%+v", file)
			log.Printf("%+v", header)
			log.Printf("Error in formfile")
			return
		}
		defer file.Close()

		out, err := os.Create(config.FactorioSavesDir + "/" + header.Filename)
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
			log.Printf("Error in out: %s", err)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
			log.Printf("Error in io copy: %s", err)
			return
		}
		log.Printf("Uploaded save file: %s", header.Filename)
		resp := "File '" + header.Filename + "' uploaded successfully"
		json.NewEncoder(w).Encode(resp)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Deletes provided save
//TODO sanitize
func RemoveSave(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	vars := mux.Vars(r)
	saveName := vars["save"]

	err := rmSave(saveName)
	if err == nil {
		// save was removed
		resp := fmt.Sprintf("Removed save: %s", saveName)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error removing save %s", err)
		}
	} else {
		log.Printf("Error in remove save handler: %s", err)
		resp := fmt.Sprintf("Error in remove save handler: %s", err)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error removing save: %s", err)
		}
	}
}

// Returns last lines of the factorio-current.log file
func LogTail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	logLines, err := tailLog(config.FactorioLog)
	if err != nil {
		log.Printf("Could not tail %s: %s", config.FactorioLog, err)
		return
	}

	if err := json.NewEncoder(w).Encode(logLines); err != nil {
		log.Printf("Error tailing logfile", err)
	}
}

// Return JSON response of config.ini file
func LoadConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	configContents, err := loadConfig(config.FactorioConfigFile)
	if err != nil {
		log.Printf("Could not retrieve config.ini: %s", err)
		resp := "Error getting config.ini"
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error tailing logfile", err)
		}
		return
	}
	if _, err := w.Write(configContents); err != nil {
		log.Printf("Error encoding config.ini response: %s", err)
	}
	log.Printf("Sent config.ini response")
}
