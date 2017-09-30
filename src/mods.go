package main

import (
    "net/http"
    "net/url"
    "log"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
    "errors"
    "archive/zip"
    "encoding/json"
    "bytes"
)

type LoginErrorResponse struct {
    Message string  `json:"message"`
    Status  int     `json:"status"`
}
type LoginSuccessResponse struct {
    UserKey []string  `json:""`
}
type FactorioCredentials struct {
    Username    string  `json:"username"`
    Userkey     string  `json:"userkey"`
}

func (credentials *FactorioCredentials) save() error {
    var err error

    credentials_json, err := json.Marshal(credentials)
    if err != nil {
        log.Printf("error mashalling the credentials: %s", err)
        return err
    }

    err = ioutil.WriteFile(config.FactorioCredentialsFile, credentials_json, 0664)
    if err != nil {
        log.Printf("error on saving the credentials. %s", err)
        return err
    }

    return nil
}

func (credentials *FactorioCredentials) load() (bool, error) {
    var err error

    if _, err := os.Stat(config.FactorioCredentialsFile); os.IsNotExist(err) {
        return false, nil
    }

    file_bytes, err := ioutil.ReadFile(config.FactorioCredentialsFile)
    if err != nil {
        credentials.del()
        log.Printf("error reading CredentialsFile: %s", err)
        return false, err
    }

    err = json.Unmarshal(file_bytes, credentials)
    if err != nil {
        credentials.del()
        log.Printf("error on unmarshal credentials_file: %s", err)
        return false, err
    }

    if credentials.Userkey != "" && credentials.Username != "" {
        return true, nil
    } else {
        credentials.del()
        return false, errors.New("incredients incomplete")
    }
}

func (credentials *FactorioCredentials) del() error {
    var err error

    err = os.Remove(config.FactorioCredentialsFile)
    if err != nil {
        log.Printf("error delete the credentialfile: %s", err)
        return err
    }

    return nil
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

    body_bytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("error on reading resp.Body: %s", err)
        return "", err, http.StatusInternalServerError
    }

    body_string := string(body_bytes)

    if resp.StatusCode != http.StatusOK {
        log.Println("error Statuscode not 200")
        return body_string, errors.New(body_string), resp.StatusCode
    }

    var success_response []string
    err = json.Unmarshal(body_bytes, &success_response)
    if err != nil {
        log.Printf("error on unmarshal body: %s", err)
        return err.Error(), err, http.StatusInternalServerError
    }

    credentials := FactorioCredentials{
        Username: username,
        Userkey: success_response[0],
    }

    err = credentials.save()
    if err != nil {
        log.Printf("error saving the credentials. %s", err)
        return err.Error(), err, http.StatusInternalServerError
    }

    return "", nil, http.StatusOK
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
        return "error", err, http.StatusInternalServerError
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

func deleteAllMods() (error) {
    var err error

    mods_dir_info, err := os.Stat(config.FactorioModsDir)
    if err != nil {
        log.Printf("error getting stats of FactorioModsDir: %s", err)
        return err
    }

    mods_dir_perm := mods_dir_info.Mode().Perm()

    err = os.RemoveAll(config.FactorioModsDir)
    if err != nil {
        log.Printf("removing FactorioModsDir failed: %s", err)
        return err
    }

    err = os.Mkdir(config.FactorioModsDir, mods_dir_perm)
    if err != nil {
        log.Printf("error recreating modPackDir: %s", err)
        return err
    }

    return nil
}

func modStartUp() {
    var err error

    //get main-folder info
    factorioDir_info, err := os.Stat(config.FactorioDir)
    if err != nil {
        log.Printf("error getting stats from FactorioDir: %s", err)
        return
    }
    factorioDir_perm := factorioDir_info.Mode().Perm()

    //create mods dir
    if _, err = os.Stat(config.FactorioModsDir); os.IsNotExist(err) {
        log.Println("no mods dir found ... creating one ...")
        os.Mkdir(config.FactorioModsDir, factorioDir_perm)
    }

    //crate mod_pack dir
    if _, err = os.Stat(config.FactorioModPackDir); os.IsNotExist(err) {
        log.Println("no ModPackDir found ... creating one ...")
        os.Mkdir(config.FactorioModPackDir, factorioDir_perm)
    }

    old_modpack_dir := filepath.Join(config.FactorioDir, "modpacks")
    if _, err := os.Stat(filepath.Join(old_modpack_dir)); !os.IsNotExist(err) {
        log.Printf("found old modpack files, rebuild into new system...")

        err = filepath.Walk(old_modpack_dir, func(path string, info os.FileInfo, err error) error {
            if info.IsDir() {
                return nil
            }
            if filepath.Ext(info.Name()) != ".zip" {
                log.Printf("file is not a zip or a directory -> skip")
                return nil
            }

            filename := info.Name()
            n := strings.LastIndexByte(info.Name(), '.')
            mod_pack_name := filename[:n]

            log.Printf("loading modPack %s into new system ...", mod_pack_name)

            mod_pack_dir := filepath.Join(config.FactorioModPackDir, mod_pack_name)

            if _, err := os.Stat(mod_pack_dir); !os.IsNotExist(err) {
                log.Printf("modPack already exists")
                return errors.New("modPack already exists")
            }

            err = os.Mkdir(mod_pack_dir, factorioDir_perm)
            if err != nil {
                log.Printf("error creating newModPackDir: %s", err)
                return err
            }

            //create mod-info.json
            mod_simple_list := ModSimpleList{
                Destination: mod_pack_dir,
                Mods: []ModSimple{
                    ModSimple{
                        Name: "base",
                        Enabled: true,
                    },
                },
            }
            new_json, _ := json.Marshal(mod_simple_list)

            err = ioutil.WriteFile(mod_simple_list.Destination + "/mod-list.json", new_json, 0664)
            if err != nil {
                log.Printf("error when writing new mod-list: %s", err)
                return err
            }

            mod_pack_file, err := zip.OpenReader(path)
            if err != nil {
                return err
            }
            defer mod_pack_file.Close()

            mods, err := newMods(mod_pack_dir)
            if err != nil {
                log.Printf("error reading mods: %s", err)
                return err
            }

            for _, mod_file := range mod_pack_file.File {
                mod_file_rc, err := mod_file.Open()
                if err != nil {
                    log.Printf("error opening mod_file: %s", err)
                    return err
                }
                defer mod_file_rc.Close()

                mod_file_buffer, err := ioutil.ReadAll(mod_file_rc)
                if err != nil {
                    log.Printf("error reading mod_file_rc: %s", err)
                    return err
                }
                mod_file_rc.Close()

                mod_file_byte_reader := bytes.NewReader(mod_file_buffer)
                mod_file_zip_reader, err := zip.NewReader(mod_file_byte_reader, int64(len(mod_file_buffer)))
                if err != nil {
                    log.Printf("error creating Reader on byte_array: %s", err)
                    return err
                }

                var mod_info ModInfo
                err = mod_info.getModInfo(mod_file_zip_reader)
                if err != nil {
                    log.Printf("error loading the ModInfo: %s", err)
                    return err
                }

                err = mods.createMod(mod_info.Name, mod_file.Name, bytes.NewReader(mod_file_buffer))
                if err != nil {
                    log.Printf("error on creating mod: %s", err)
                    return err
                }
            }

            log.Printf("loading modPack %s successful", mod_pack_name)

            return nil
        })

        if err != nil {
            log.Printf("error on loading old modpacks into the new system: %s\n please check if empty modPacks are creating and delete them", err)
        } else {
            log.Printf("all modPacks are loaded into the new system successfully")
            log.Printf("deleting old modPackDir")
            os.RemoveAll(old_modpack_dir)
        }
    }
}
