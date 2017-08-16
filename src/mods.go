package main

import (
    "io/ioutil"
    "log"
    "encoding/json"
    "net/http"
    "net/url"
	"os"
	"io"
	"path/filepath"
	"archive/zip"
)

type Mod struct {
    Name    string `json:"name"`
    Enabled bool   `json:"enabled"`
}

type ModsList struct {
    Mods    []Mod   `json:"mods"`
}

// List mods installed in the factorio/mods directory
func listInstalledMods() (ModsList, error) {
    file, err := ioutil.ReadFile(config.FactorioModsDir + "/mod-list.json")

    if err != nil {
        log.Println(err.Error())
    }

    var result ModsList
    err_json := json.Unmarshal(file, &result)

    if err_json != nil {
        log.Println(err_json.Error())
        return result, err_json
    }

	return result, nil
}

type ModInfoList struct {
	Mods []ModInfo `json:"mods"`
}
type ModInfo struct {
	Name string `json:"name"`
	Version string `json:"version"`
	Title string `json:"title"`
	Author string `json:"author"`
	FileName string `json:"file_name"`
	Enabled bool `json:"enabled"`
}
func listInstalledModsByFolder() (ModInfoList, error) {
	//scan ModFolder
	var result ModInfoList
	var err_o error
	err_o = filepath.Walk(config.FactorioModsDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".zip" {
			zip_file, err := zip.OpenReader(path)
			if err != nil {
				log.Fatalln(err)
			}

			//iterate through all files inside the zip (search for info.json)
			for _, single_file := range zip_file.File {
				if single_file.FileInfo().Name() == "info.json" {
					//interpret info.json
					rc, err := single_file.Open()

					if err != nil {
						log.Fatal(err)
						return err
					}

					byte_array, err := ioutil.ReadAll(rc)
					rc.Close()
					if err != nil {
						log.Fatal(err)
						return err
					}

					var mod_info ModInfo
					err = json.Unmarshal(byte_array, &mod_info)
					if err != nil {
						log.Fatalln(err)
						return err
					}

					mod_info.FileName = info.Name()
					mod_info.FileName = info.Name()
					result.Mods = append(result.Mods, mod_info)

					break
				}
			}
		}

		return nil
	})

	if err_o != nil {
		return ModInfoList{}, err_o
	}

	mod_list_by_json, err_o := listInstalledMods()
	if err_o != nil {
		return ModInfoList{}, err_o
	}

	for _, json_mod := range mod_list_by_json.Mods {
		for result_index, result_mod := range result.Mods {
			if result_mod.Name == json_mod.Name {
				result.Mods[result_index].Enabled = json_mod.Enabled
				break
			}
		}
	}

	return result, nil
}

func toggleMod(mod_name string)([]ModInfo, error) {
	var err error

	mod_list, err := listInstalledMods()

	for index, mod := range mod_list.Mods {
		if mod.Name == mod_name {
			mod_list.Mods[index].Enabled = !mod_list.Mods[index].Enabled
			break
		}
	}

	if err != nil {
		return nil, err
	}

	//build new json
	new_json, _ := json.Marshal(mod_list)

	ioutil.WriteFile(config.FactorioModsDir + "/mod-list.json", new_json, 0664)

	mod_info_list, err := listInstalledModsByFolder()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return mod_info_list.Mods, nil
}


type LoginErrorResponse struct {
    Message string  `json:"message"`
    Status  int     `json:"status"`
}
type LoginSuccessResponse struct {
    UserKey []string  `json:""`
}
//Log the user into factorio, so mods can be downloaded
func getUserToken(username string, password string) (string, error, int) {
    resp, get_err := http.PostForm("https://auth.factorio.com/api-login",
        url.Values{"require_game_ownership": {"true"}, "username": {username}, "password": {password}})
    if get_err != nil {
        log.Fatal(get_err)
        return "error", get_err, 500
    }

    //get the response-text
    text, err_io := ioutil.ReadAll(resp.Body)
    resp.Body.Close()

    text_string := string(text)

    if err_io != nil {
        log.Fatal(err_io)
        return "error", err_io, resp.StatusCode
    }

    return text_string, nil, resp.StatusCode
}


//Search inside the factorio mod portal
func searchModPortal(keyword string) (string, error, int) {
    //resp, get_err := http.Get
    req, err := http.NewRequest(http.MethodGet, "https://mods.factorio.com/api/mods", nil)
    if err != nil {
        return "error", err, 500
    }

    query := req.URL.Query()
    query.Add("q", keyword)
    req.URL.RawQuery = query.Encode()

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "error", err, 500
    }

    text, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    if err != nil {
        return "error", err, 500
    }

    text_string := string(text)

    return text_string, nil, resp.StatusCode
}

func getModDetails(modId string) (string, error, int) {
    var err error
    new_link := "https://mods.factorio.com/api/mods/" + modId
    resp, err := http.Get(new_link)

    if err != nil {
        return "error", err, 500
    }

    //get the response-text
    text, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()

    text_string := string(text)

    if err != nil {
        log.Fatal(err)
        return "error", err, resp.StatusCode
    }

    return text_string, nil, resp.StatusCode
}

func installMod(username string, userKey string, url string, filename string, mod_id string) ([]ModInfo, error, int) {
	var err error
	//download the mod from the mod portal api
	complete_url := "https://mods.factorio.com" + url + "?username=" + username + "&token=" + userKey

	// don't worry about errors
	response, err := http.Get(complete_url)
	if err != nil {
		log.Fatal(err)
		return nil, err, 500
	}

	if response.StatusCode != 200 {
		text, _ := ioutil.ReadAll(response.Body)
		log.Printf("StatusCode: %d \n ResponseBody: %s", response.StatusCode, text)

		defer response.Body.Close()
		return nil, err, response.StatusCode
	}

	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create(config.FactorioModsDir + "/" + filename)
	if err != nil {
		log.Fatal(err)
		return nil,  err, 500
	}
	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err, 500
	}
	file.Close()

	mod_list, err := listInstalledMods()

	if err != nil {
		return nil, err, 500
	}

	//add new mod
	new_mod_entry := Mod{
		Name: mod_id,
		Enabled:true,
	}
	mod_list.Mods = append(mod_list.Mods, new_mod_entry)

	//build new json
	new_json, _ := json.Marshal(mod_list)

	ioutil.WriteFile(config.FactorioModsDir + "/mod-list.json", new_json, 0664)

	mod_info_list, err := listInstalledModsByFolder()
	if err != nil {
		log.Fatal(err)
		return nil, err, 500
	}

	return mod_info_list.Mods, nil, response.StatusCode
}
