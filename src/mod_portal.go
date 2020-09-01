package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// get all mods uploaded to the factorio modPortal
func modPortalList() (interface{}, error, int) {
	req, err := http.NewRequest(http.MethodGet, "https://mods.factorio.com/api/mods?page_size=max", nil)
	if err != nil {
		return "error", err, 500
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "error", err, 500
	}

	text, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "error", err, 500
	}

	var jsonVal interface{}
	json.Unmarshal(text, &jsonVal)

	return jsonVal, nil, resp.StatusCode
}

// get the details (mod-info, releases, etc.) from a specific mod from the modPortal
func modPortalModDetails(modId string) (interface{}, error, int) {
	req, err := http.NewRequest(http.MethodGet, "https://mods.factorio.com/api/mods/"+modId+"/full", nil)
	if err != nil {
		return "error", err, 500
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "error", err, 500
	}

	text, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "error", err, 500
	}

	var jsonVal interface{}
	json.Unmarshal(text, &jsonVal)

	return jsonVal, nil, resp.StatusCode
}
