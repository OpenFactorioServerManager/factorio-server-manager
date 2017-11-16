package main

import (
    "log"
    "encoding/json"
    "io/ioutil"
    "os"
    "path/filepath"
    "archive/zip"
    "errors"
    "io"
    "lockfile"
)

type ModInfoList struct {
    Mods []ModInfo `json:"mods"`
    Destination string `json:"-"`
}
type ModInfo struct {
    Name string `json:"name"`
    Version string `json:"version"`
    Title string `json:"title"`
    Author string `json:"author"`
    FileName string `json:"file_name"`
    FactorioVersion string `json:"factorio_version"`
}

func newModInfoList(destination string) (ModInfoList, error) {
    var err error
    mod_info_list := ModInfoList{
        Destination: destination,
    }

    err = mod_info_list.listInstalledMods()
    if err != nil {
        log.Printf("ModInfoList ... error listing installed Mods: %s", err)
        return mod_info_list, err
    }

    return mod_info_list, nil
}

func (mod_info_list *ModInfoList) listInstalledMods() (error) {
    var err error
    mod_info_list.Mods = nil

    err = filepath.Walk(mod_info_list.Destination, func(path string, info os.FileInfo, err error) error {
        if !info.IsDir() && filepath.Ext(path) == ".zip" {
            err = fileLock.RLock(path)
            if err != nil && err == lockfile.ErrorAlreadyLocked {
                log.Println(err)
                return nil
            } else if err != nil {
                log.Printf("error locking file: %s", err)
                return err
            }
            defer fileLock.RUnlock(path)

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
            mod_info_list.Mods = append(mod_info_list.Mods, mod_info)
        }

        return nil
    })

    if err != nil {
        log.Printf("error while walking over the given dir: %s", err)
        return err
    }

    return nil
}

func (mod_info_list *ModInfoList) deleteMod(mod_name string) (error) {
    var err error

    //search for mod, that should be deleted
    for _, mod := range mod_info_list.Mods {
        if mod.Name == mod_name {
            filePath := mod_info_list.Destination + "/" + mod.FileName

            fileLock.LockW(filePath)
            //delete mod
            err = os.Remove(filePath)
            fileLock.Unlock(filePath)
            if err != nil {
                log.Printf("ModInfoList ... error when deleting mod: %s", err)
                return err
            }

            //reload mod-list
            err = mod_info_list.listInstalledMods()
            if err != nil {
                log.Printf("ModInfoList ... error while refreshing installedModList: %s", err)
                return err
            }

            return nil
        }
    }

    log.Printf("the mod-file doesnt exists!")
    return nil
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

            err = json.Unmarshal(byte_array, mod_info)
            if err != nil {
                log.Fatalln(err)
                return err
            }

            return nil
        }
    }

    return errors.New("info.json not found in zip-file")
}

func (mod_info_list *ModInfoList) createMod(mod_name string, file_name string, mod_file io.Reader) error {
    var err error

    //save uploaded file
    filePath := mod_info_list.Destination + "/" + file_name
    new_file, err := os.Create(filePath)
    if err != nil {
        log.Printf("error on creating new file - %s: %s", file_name, err)
        return err
    }
    defer new_file.Close()

    fileLock.LockW(filePath)

    _, err = io.Copy(new_file, mod_file)
    if err != nil {
        log.Printf("error on copying file to disk: %s", err)
        return err
    }

    err = new_file.Close()
    if err != nil {
        log.Printf("error on closing new created zip-file: %s", err)
        return err
    }

    fileLock.Unlock(filePath)

    //reload the list
    err = mod_info_list.listInstalledMods()
    if err != nil {
        log.Printf("error on listing mod-infos: %s", err)
        return err
    }

    return nil
}
