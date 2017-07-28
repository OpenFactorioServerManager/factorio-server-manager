package main

import (
    "io/ioutil"
    "log"
    "encoding/json"
    "os"
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
        os.Exit(1)
    }

    log.Print(file)

    var result ModsList
    err_json := json.Unmarshal(file, &result)

    if err_json != nil {
        log.Println(err_json.Error())
        os.Exit(1)
    }

    log.Printf("%v", result)

	return result.Mods, nil
}
