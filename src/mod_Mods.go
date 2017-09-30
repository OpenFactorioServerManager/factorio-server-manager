package main

import (
    "log"
    "net/http"
    "io"
    "mime/multipart"
    "bytes"
    "archive/zip"
    "errors"
    "fmt"
    "io/ioutil"
    "path/filepath"
)

type Mods struct {
    ModSimpleList ModSimpleList `json:"mod_simple_list"`
    ModInfoList ModInfoList `json:"mod_info_list"`
}
type ModsResult struct {
    ModInfo
    Enabled bool `json:"enabled"`
}
type ModsResultList struct {
    ModsResult []ModsResult `json:"mods"`
}

func newMods(destination string) (Mods, error) {
    var err error
    var mods Mods

    mods.ModSimpleList, err = newModSimpleList(destination)
    if err != nil {
        log.Printf("error on creating newModSimpleList: %s", err)
        return mods, err
    }

    mods.ModInfoList, err = newModInfoList(destination)
    if err != nil {
        log.Printf("error on creating newModInfoList: %s", err)
        return mods, err
    }

    return mods, nil
}

func (mods *Mods) listInstalledMods() ModsResultList {
    var result ModsResultList

    for _, mod_info := range mods.ModInfoList.Mods {
        var mods_result ModsResult
        mods_result.Name = mod_info.Name
        mods_result.FileName = mod_info.FileName
        mods_result.Author = mod_info.Author
        mods_result.Title = mod_info.Title
        mods_result.Version = mod_info.Version

        for _, simple_mod := range mods.ModSimpleList.Mods {
            if simple_mod.Name == mods_result.Name {
                mods_result.Enabled = simple_mod.Enabled
                break
            }
        }

        result.ModsResult = append(result.ModsResult, mods_result)
    }

    return result
}

func (mods *Mods) deleteMod(mod_name string) (error) {
    var err error

    err = mods.ModInfoList.deleteMod(mod_name)
    if err != nil {
        log.Printf("error when deleting mod in ModInfoList: %s", err)
        return err
    }

    err = mods.ModSimpleList.deleteMod(mod_name)
    if err != nil {
        log.Printf("error when deleting mod in ModSimpleList: %s", err)
        return err
    }

    return nil
}

func (mods *Mods) createMod(mod_name string, file_name string, file_rc io.Reader) (error) {
    var err error

    //check if mod already exists and delete it
    if mods.ModSimpleList.checkModExists(mod_name) {
        err = mods.ModSimpleList.deleteMod(mod_name)
        if err != nil {
            log.Printf("error when deleting mod: %s", err)
            return err
        }

        err = mods.ModInfoList.deleteMod(mod_name)
        if err != nil {
            log.Printf("error when deleting mod: %s", err)
            return err
        }
    }


    //create new mod
    err = mods.ModInfoList.createMod(mod_name, file_name, file_rc)
    if err != nil {
        log.Printf("error on creating mod-file: %s", err)
        return err
    }
    err = mods.ModSimpleList.createMod(mod_name)
    if err != nil {
        log.Printf("error on adding mod to the mod-list.json: %s", err)
        return err
    }

    return nil
}

func (mods *Mods) downloadMod(url string, filename string, mod_id string) (error) {
    var err error

    var credentials FactorioCredentials
    status, err := credentials.load()
    if err != nil {
        log.Printf("error loading credentials: %s", err)
        return err
    }
    if status == false {
        log.Printf("error: credentials are invalid")
        return errors.New("error: credentials are invalid")
    }

    //download the mod from the mod portal api
    complete_url := "https://mods.factorio.com" + url + "?username=" + credentials.Username + "&token=" + credentials.Userkey

    response, err := http.Get(complete_url)
    if err != nil {
        log.Printf("error on downloading mod: %s", err)
        return err
    }

    defer response.Body.Close()

    if response.StatusCode != 200 {
        log.Printf("StatusCode: %d", response.StatusCode)

        return errors.New("Statuscode not 200: " + fmt.Sprint(response.StatusCode))
    }

    err = mods.createMod(mod_id, filename, response.Body)
    if err != nil {
        log.Printf("error when creating Mod: %s", err)
        return err
    }

    //done everything is made inside the createMod

    return nil
}

func (mods *Mods) uploadMod(header *multipart.FileHeader) (error) {
    var err error

    if filepath.Ext(header.Filename) != ".zip" {
        log.Print("The uploaded file wasn't a zip-file -> ignore it")
        return nil //simply do nothing xD
    }

    file, err := header.Open()
    if err != nil {
        log.Printf("error on open file via fileHeader. %s", err)
        return err
    }
    defer file.Close()

    file_byte_array, err := ioutil.ReadAll(file)
    if err != nil {
        log.Printf("error reading file: %s", err)
        return err
    }

    zip_reader, err := zip.NewReader(bytes.NewReader(file_byte_array), int64(len(file_byte_array)))
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

    err = mods.createMod(mod_info.Name, header.Filename, bytes.NewReader(file_byte_array))
    if err != nil {
        log.Printf("error on creating Mod: %s", err)
        return err
    }

    return nil
}

func (mods *Mods) updateMod(mod_name string, url string, filename string) error {
    var err error

    //err = mods.deleteMod(mod_name)
    //if err != nil {
    //    log.Printf("updateMod ... error when deleting mod: %s", err)
    //    return err
    //}

    err = mods.downloadMod(url, filename, mod_name)
    if err != nil {
        log.Printf("updateMod ... error when downloading the new Mod: %s", err)
        return err
    }

    return nil
}