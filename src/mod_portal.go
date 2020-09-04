package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
func modPortalModDetails(modId string) (ModPortalStruct, error, int) {
	var jsonVal ModPortalStruct

	req, err := http.NewRequest(http.MethodGet, "https://mods.factorio.com/api/mods/"+modId, nil)
	if err != nil {
		return jsonVal, err, 500
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return jsonVal, err, 500
	}

	text, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return jsonVal, err, 500
	}

	json.Unmarshal(text, &jsonVal)

	return jsonVal, nil, resp.StatusCode
}

//Log the user into factorio, so mods can be downloaded
func factorioLogin(username string, password string) (string, error, int) {
	var err error

	resp, err := http.PostForm("https://auth.factorio.com/api-login",
		url.Values{"require_game_ownership": {"true"}, "username": {username}, "password": {password}})

	if err != nil {
		log.Printf("error on logging in: %s", err)
		return "", err, resp.StatusCode
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error on reading resp.Body: %s", err)
		return "", err, http.StatusInternalServerError
	}

	bodyString := string(bodyBytes)

	if resp.StatusCode != http.StatusOK {
		log.Println("error Statuscode not 200")
		return bodyString, errors.New(bodyString), resp.StatusCode
	}

	var successResponse []string
	err = json.Unmarshal(bodyBytes, &successResponse)
	if err != nil {
		log.Printf("error on unmarshal body: %s", err)
		return err.Error(), err, http.StatusInternalServerError
	}

	credentials := FactorioCredentials{
		Username: username,
		Userkey:  successResponse[0],
	}

	err = credentials.save()
	if err != nil {
		log.Printf("error saving the credentials. %s", err)
		return err.Error(), err, http.StatusInternalServerError
	}

	return "", nil, http.StatusOK
}
