package main

import (
    "io/ioutil"
    "log"
    "encoding/json"
    "os"
)

type Mod struct {
    Name    string `json:"name"`
    Enabled bool   `json:"enabled,string"`
}

// List mods installed in the factorio/mods directory
func listInstalledMods(modDir string) ([]Mod, error) {
	result := []Mod{}

    file, err := ioutil.ReadFile(modDir + "mod-list.json")

    if err != nil {
        log.Println(err.Error())
        os.Exit(1)
    }

    err_json := json.Unmarshal(file, &result)

    if err_json != nil {
        log.Println(err.Error())
        os.Exit(1)
    }

    log.Printf("%v", result)


	//files, err := ioutil.ReadDir(modDir)
	//if err != nil {
	//	log.Printf("Error listing installed mods: %s", err)
	//	return result, err
	//}
	//for _, f := range files {
	//	if f.Name() == "mod-list.json" {
	//		continue
	//	}
	//	result = append(result, f.Name())
	//}

	return result, nil
}
