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
	"mime/multipart"
    "errors"
    "bytes"
)

type Mod struct {
    Name    string `json:"name"`
    Enabled bool   `json:"enabled"`
}

type ModsList struct {
    Mods    []Mod   `json:"mods"`
}

func (mod_list *ModsList) check_mod_exists(mod_name string) bool {
    for _, single_mod := range mod_list.Mods {
        if single_mod.Name == mod_name {
            return true
        }
    }

    return false
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

func (mod_info *ModInfo) getModInfo(reader *zip.Reader) error {
    for _, single_file := range reader.File {
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

            //var mod_info ModInfo
            err = json.Unmarshal(byte_array, mod_info)
            if err != nil {
                log.Fatalln(err)
                return err
            }

            return nil
        }
    }

    return errors.New("info.json not found in zip-file!")
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
                return err
			}

            var mod_info ModInfo
            err = mod_info.getModInfo(&zip_file.Reader)
            if err != nil {
                log.Fatalf("Error in getModInfo: %s", err)
            }

            mod_info.FileName = info.Name()
            result.Mods = append(result.Mods, mod_info)
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

	if err != nil {
		return nil, err
	}

	for index, mod := range mod_list.Mods {
		if mod.Name == mod_name {
			mod_list.Mods[index].Enabled = !mod_list.Mods[index].Enabled
			break
		}
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

func deleteMod(mod_name string) ([]ModInfo, error) {
	var err error
	mod_list, err := listInstalledMods()

	if err != nil {
		return nil, err
	}

	for index, mod := range mod_list.Mods {
		if mod.Name == mod_name {
			slice1 := mod_list.Mods[:index]
			slice2 := mod_list.Mods[index + 1:]
			var new_mod_list []Mod
			new_mod_list = append(new_mod_list, slice1...)
			new_mod_list = append(new_mod_list, slice2...)
			mod_list.Mods = new_mod_list
			break
		}
	}

	//build new json
	new_json, _ := json.Marshal(mod_list)

	ioutil.WriteFile(config.FactorioModsDir + "/mod-list.json", new_json, 0664)

	mod_info_list, err := listInstalledModsByFolder()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var delete_file_name string
	//search for mod in own setup
	for index, mod := range mod_info_list.Mods {
		if mod.Name == mod_name {
			delete_file_name = mod.FileName

			//remove mod from list (faster than scanning path new)
			slice1 := mod_info_list.Mods[:index]
			slice2 := mod_info_list.Mods[index+1:]
			var new_mod_list []ModInfo
			new_mod_list = append(new_mod_list, slice1...)
			new_mod_list = append(new_mod_list, slice2...)
			mod_info_list.Mods = new_mod_list
			break
		}
	}

	os.Remove(config.FactorioModsDir + "/" + delete_file_name)

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

    mod_list, err := listInstalledMods()
    if err != nil {
        return nil, err, 500
    }

    if mod_list.check_mod_exists(mod_id) {
        log.Printf("delete old mod %s.", mod_id)
        _, err = deleteMod(mod_id)
        if err != nil {
            log.Printf("error on deleting mod: %s", err)
            return nil, err, 500
        }
    }

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

    mod_list, err = listInstalledMods()
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

func uploadMod(header *multipart.FileHeader) (error) {
    var err error
    if header.Header.Get("Content-Type") != "application/zip" {
        log.Print("The uploaded file wasn't a zip-file -> ignore it")
        return nil //simply do nothing xD
    }

    if _,err_file := os.Stat(config.FactorioModsDir + "/" + header.Filename); !os.IsNotExist(err_file) {
        log.Print("The uploaded file already exists -> ignore it")
        return nil //simply do nothing xD
    }

    file, err := header.Open()
    if err != nil {
        log.Printf("error on open file via fileHeader. %s", err)
        return err
    }

    var buff bytes.Buffer
    file_length, err := buff.ReadFrom(file)
    if err != nil {
        log.Printf("Error occured while reading bytes.Buffer.ReadFrom: %s", err)
        return err
    }

    zip_reader, err := zip.NewReader(file, file_length)
    if err != nil {
        log.Printf("Uploaded file could not put into zip.Reader: %s", err)
        return err
    }

    var mod_info ModInfo
    err = mod_info.getModInfo(zip_reader)
    if err != nil {
        log.Printf("Error in getModInfo: %s", err)
        return err
    }

    //check if mod already exists in mod_list.json
    mods_list, err := listInstalledMods()
    if err != nil {
        log.Printf("Error in listInstalledMods: %s", err)
        return err
    }

    mod_already_exists := mods_list.check_mod_exists(mod_info.Name)

    if mod_already_exists {
        _, err = deleteMod(mod_info.Name)
        if err != nil {
            log.Printf("error when trying to delete mod: %s", err)
            return err
        }
    }

    //save uploaded file
    new_file, err := os.Create(config.FactorioModsDir + "/" + header.Filename)
    if err != nil {
        log.Printf("error on creating new file - %s: %s", header.Filename, err)
        return err
    }
    defer new_file.Close()

    file.Seek(0,0) //reset file-cursor to 0,0
    _, err = io.Copy(new_file, file)
    if err != nil {
        log.Printf("error on copying file to disk: %s", err)
        return err
    }

    //build new json
    mods_list, err = listInstalledMods()
    if err != nil {
        log.Printf("Error in listInstalledMods: %s", err)
        return err
    }

    //add this mod
    mods_list.Mods = append(mods_list.Mods, Mod{
        Name: mod_info.Name,
        Enabled: true,
    })

    //save mod-list.json with the new mod
    new_json, _ := json.Marshal(mods_list)
    ioutil.WriteFile(config.FactorioModsDir + "/mod-list.json", new_json, 0664)

	return nil
}
