package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func FactorioModPortalListModsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	var statusCode int
	resp, err, statusCode = modPortalList()

	if err != nil {
		resp = fmt.Sprintf("Error in listing mods from mod portal: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
}

// FactorioModPortalModInfoHandler returns JSON response with the mod details
func FactorioModPortalModInfoHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	vars := mux.Vars(r)
	modId := vars["mod"]

	var statusCode int
	resp, err, statusCode = modPortalModDetails(modId)

	if err != nil {
		resp = fmt.Sprintf("Error in getting mod details from mod portal: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
}

func FactorioModPortalInstallHandler(w http.ResponseWriter, r *http.Request) {
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

	// Get Data out of the request
	var data struct {
		DownloadURL string `json:"link"`
		Filename    string `json:"filename"`
		ModName     string `json:"modName"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		resp = fmt.Sprintf("Error reading data from request {%s}: %s", r.RequestURI, err)
		log.Println(resp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mods, err := CreateNewMods(w, &resp)
	if err != nil {
		return
	}

	err = mods.downloadMod(data.DownloadURL, data.Filename, data.ModName)
	if err != nil {
		resp = fmt.Sprintf("Error downloading a mod: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp = mods.listInstalledMods()
}

func FactorioModPortalLoginHandler(w http.ResponseWriter, r *http.Request) {
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

	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		resp = fmt.Sprintf("Error reading data from request {%s}: %s", r.RequestURI, err)
		log.Println(resp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	loginStatus, err, statusCode := factorioLogin(data.Username, data.Password)
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
