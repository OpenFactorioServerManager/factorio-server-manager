package main

import (
    "path/filepath"
    "log"
    "os"
    "errors"
    "io/ioutil"
    "io"
)

type ModPackMap map[string]*ModPack
type ModPack struct {
    Mods Mods
}

type ModPackResult struct {
    Name string `json:"name"`
    Mods ModsResultList `json:"mods"`
}
type ModPackResultList struct {
    ModPacks []ModPackResult `json:"mod_packs"`
}

func newModPackMap() (ModPackMap, error) {
    var err error
    //var mod_pack_map ModPackMap
    mod_pack_map := make(ModPackMap)

    err = mod_pack_map.reload()
    if err != nil {
        log.Printf("error on loading the modpacks: %s", err)
        return mod_pack_map, err
    }

    return mod_pack_map, nil
}

func newModPack(mod_pack_folder string) (*ModPack, error) {
    var err error
    var mod_pack ModPack

    mod_pack.Mods, err = newMods(mod_pack_folder)
    if err != nil {
        log.Printf("error on loading mods in mod_pack_dir: %s", err)
        return &mod_pack, err
    }

    return &mod_pack, err
}

func (mod_pack_map *ModPackMap) reload() error {
    var err error
    new_mod_pack_map := make(ModPackMap)

    err = filepath.Walk(config.FactorioModPackDir, func(path string, info os.FileInfo, err error) error {
        if path == config.FactorioModPackDir || !info.IsDir() {
            return nil
        }

        mod_pack_name := filepath.Base(path)

        new_mod_pack_map[mod_pack_name], err = newModPack(path)
        if err != nil {
            log.Printf("error on creating newModPack: %s", err)
            return err
        }

        return nil
    })
    if err != nil {
        log.Printf("error on walking over the ModDir: %s", err)
        return err
    }

    *mod_pack_map = new_mod_pack_map

    return nil
}

func (mod_pack_map *ModPackMap) listInstalledModPacks() ModPackResultList {
    var mod_pack_result_list ModPackResultList

    for mod_pack_name, mod_pack := range *mod_pack_map {
        var mod_pack_result ModPackResult
        mod_pack_result.Name = mod_pack_name
        mod_pack_result.Mods = mod_pack.Mods.listInstalledMods()

        mod_pack_result_list.ModPacks = append(mod_pack_result_list.ModPacks, mod_pack_result)
    }

    return mod_pack_result_list
}

func (mod_pack_map *ModPackMap) createModPack(mod_pack_name string) error {
    var err error

    mod_pack_folder := filepath.Join(config.FactorioModPackDir, mod_pack_name)

    if mod_pack_map.checkModPackExists(mod_pack_name) == true {
        log.Printf("ModPack %s already existis", mod_pack_name)
        return errors.New("ModPack " + mod_pack_name + " already exists, please choose a different name")
    }

    source_file_info, err := os.Stat(config.FactorioModsDir)
    if err != nil {
        log.Printf("error when reading factorioModsDir. %s", err)
        return err
    }

    //Create the modPack-folder
    err = os.MkdirAll(mod_pack_folder, source_file_info.Mode())
    if err != nil {
        log.Printf("error on creating the new ModPack directory: %s", err)
        return err
    }

    files, err := ioutil.ReadDir(config.FactorioModsDir)
    if err != nil {
        log.Printf("error on reading the dactorio mods dir: %s", err)
        return err
    }

    for _, file := range files {
        if file.IsDir() == false {
            source_filepath := filepath.Join(config.FactorioModsDir, file.Name())
            destination_filepath := filepath.Join(mod_pack_folder, file.Name())

            source_file, err := os.Open(source_filepath)
            if err != nil {
                log.Printf("error on opening source_filepath: %s", err)
                return err
            }
            defer source_file.Close()

            destination_file, err := os.Create(destination_filepath)
            if err != nil {
                log.Printf("error on creating destination_filepath: %s", err)
                return err
            }
            defer destination_file.Close()

            _, err = io.Copy(destination_file, source_file)
            if err != nil {
                log.Printf("error on copying data from source to destination: %s", err)
                return err
            }

            source_file.Close()
            destination_file.Close()
        }
    }

    //reload the ModPackList
    err = mod_pack_map.reload()
    if err != nil {
        log.Printf("error on reloading ModPack: %s", err)
        return err
    }

    return nil
}

func (mod_pack_map *ModPackMap) checkModPackExists(mod_pack_name string) bool {
    for mod_pack_id := range *mod_pack_map {
        if mod_pack_id == mod_pack_name {
            return true
        }
    }

    return false
}

func (mod_pack_map *ModPackMap) deleteModPack(mod_pack_name string) error {
    var err error

    mod_pack_dir := filepath.Join(config.FactorioModPackDir, mod_pack_name)

    err = os.RemoveAll(mod_pack_dir)
    if err != nil {
        log.Printf("error on removing the ModPack: %s", err)
        return err
    }

    err = mod_pack_map.reload()
    if err != nil {
        log.Printf("error on reloading the ModPackList: %s", err)
        return err
    }

    return nil
}

func (mod_pack *ModPack) loadModPack() error {
    var err error

    //get filemode, so it can be restored
    file_info, err := os.Stat(config.FactorioModsDir)
    if err != nil {
        log.Printf("error on trying to save folder infos: %s", err)
        return err
    }
    folder_mode := file_info.Mode()

    //clean factorio mod directory
    err = os.RemoveAll(config.FactorioModsDir)
    if err != nil {
        log.Printf("error on removing the factorio mods dir: %s", err)
        return err
    }

    err = os.Mkdir(config.FactorioModsDir, folder_mode)
    if err != nil {
        log.Printf("error on recreating mod dir: %s", err)
        return err
    }

    //copy the modpack folder to the normal mods directory
    err = filepath.Walk(mod_pack.Mods.ModInfoList.Destination, func(path string, info os.FileInfo, err error) error {
        if info.IsDir() {
            return nil
        }
        new_file, err := os.Create(filepath.Join(config.FactorioModsDir, info.Name()))
        if err != nil {
            log.Printf("error on creting mod file: %s", err)
            return err
        }
        defer new_file.Close()

        old_file, err := os.Open(path)
        if err != nil {
            log.Printf("error on opening modFile: %s", err)
            return err
        }
        defer old_file.Close()

        _ ,err = io.Copy(new_file, old_file)
        if err != nil {
            log.Printf("error on copying data to the new file: %s", err)
            return err
        }

        return nil
    })
    if err != nil {
        log.Printf("error on copying the mod pack: %s", err)
        return err
    }

    //mods, err := newMods(config.FactorioModsDir)

    return nil
}

//func modPackToggleMod(mod_pack_name string, mod_name string) (ModPackList, error) {
//    //var err error
//    //var mod_pack_list ModPackList
//    //
//    //mod_pack := ModPack{
//    //    Name: mod_pack_name,
//    //}
//    //
//    //temp_dir, err := mod_pack.create_temp_dir()
//    //if err != nil {
//    //    log.Printf("error when creating temp_dir: %s", err)
//    //    return mod_pack_list, err
//    //}
//    //defer os.RemoveAll(temp_dir)
//    //
//    //var mods_list ModsList
//    //err = mods_list.listInstalledMods(temp_dir)
//    //if err != nil {
//    //    log.Printf("error on listing mods in temp_dir: %s", err)
//    //    return mod_pack_list, err
//    //}
//    //
//    //log.Print(mods_list)
//
//    return ModPackList{}, nil
//}
