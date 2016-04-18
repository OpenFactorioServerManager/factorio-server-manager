package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func ListInstalledMods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	mods := listInstalledMods()

	if err := json.NewEncoder(w).Encode(mods); err != nil {
		log.Printf("Error in list mods: %s", err)
	}
}

func ToggleMod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	vars := mux.Vars(r)
	modName := vars["mod"]

	m, err := parseModList()
	if err != nil {
		log.Printf("Could parse mod list: %s", err)
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

func ListSaves(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	saves := listSaves()

	if err := json.NewEncoder(w).Encode(saves); err != nil {
		log.Printf("Error listing saves: %s", err)
	}
}

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
