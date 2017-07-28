package main

import (
    "io/ioutil"
    "log"
    "encoding/json"
)

type Mod struct {
    Name    string `json:"name"`
    Enabled bool   `json:"enabled"`
}

type ModsList struct {
    Mods    []Mod   `json:"mods"`
}

// List mods installed in the factorio/mods directory
func listInstalledMods(modDir string) ([]Mod, error) {
    file, err := ioutil.ReadFile(modDir + "/mod-list.json")

    if err != nil {
        log.Println(err.Error())
        return nil, err
    }

    var result ModsList
    err_json := json.Unmarshal(file, &result)

    if err_json != nil {
        log.Println(err_json.Error())
        return result.Mods, err_json
    }

	return result.Mods, nil
}
