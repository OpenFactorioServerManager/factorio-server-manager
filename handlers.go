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

func ListMods(w http.ResponseWriter, r *http.Request) {
	mods := listInstalledMods()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)

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
