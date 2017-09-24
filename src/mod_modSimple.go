package main

import (
    "log"
    "io/ioutil"
    "encoding/json"
)

type ModSimple struct {
    Name    string `json:"name"`
    Enabled bool   `json:"enabled"`
}

type ModSimpleList struct {
    Mods    []ModSimple   `json:"mods"`
    Destination string `json:"-"`
}

func newModSimpleList(destination string) (ModSimpleList, error) {
    var err error

    mod_simple_list := ModSimpleList{
        Destination: destination,
    }

    err = mod_simple_list.listInstalledMods()
    if err != nil {
        log.Printf("ModSimpleList ... error list installed mods: %s", err)
        return mod_simple_list, err
    }

    return mod_simple_list, nil
}

func (mod_simple_list *ModSimpleList) listInstalledMods() (error) {
    var err error

    file, err := ioutil.ReadFile(mod_simple_list.Destination + "/mod-list.json")
    if err != nil {
        log.Printf("ModSimpleList ... error read the mod-info.json: %s", err)
        return err
    }

    err = json.Unmarshal(file, mod_simple_list)
    if err != nil {
        log.Printf("ModSimpleList ... error while decode mod-info.json: %s", err)
        return err
    }

    return nil
}

func (mod_simple_list *ModSimpleList) saveModInfoJson() (error) {
    var err error

    //build json of current state
    new_json, _ := json.Marshal(mod_simple_list)

    err = ioutil.WriteFile(mod_simple_list.Destination + "/mod-list.json", new_json, 0664)
    if err != nil {
        log.Printf("error when writing new mod-list: %s", err)
        return err
    }

    return nil
}

func (mod_simple_list *ModSimpleList) deleteMod(mod_name string) (error) {
    var err error

    for index, mod := range mod_simple_list.Mods {
        if mod.Name == mod_name {
            slice1 := mod_simple_list.Mods[:index]
            slice2 := mod_simple_list.Mods[index + 1:]
            var new_mod_list []ModSimple
            new_mod_list = append(new_mod_list, slice1...)
            new_mod_list = append(new_mod_list, slice2...)
            mod_simple_list.Mods = new_mod_list
            break
        }
    }

    err = mod_simple_list.saveModInfoJson()
    if err != nil {
        log.Printf("error when saving new mod_list: %s", err)
        return err
    }

    return nil
}

func (mod_simple_list *ModSimpleList) checkModExists(mod_name string) (bool) {
    for _, single_mod := range mod_simple_list.Mods {
        if single_mod.Name == mod_name {
            return true
        }
    }

    return false
}

func (mod_simple_list *ModSimpleList) createMod(mod_name string) (error) {
    var err error

    new_mod_simple := ModSimple{
        Name: mod_name,
        Enabled: true,
    }

    mod_simple_list.Mods = append(mod_simple_list.Mods, new_mod_simple)

    err = mod_simple_list.saveModInfoJson()
    if err != nil {
        log.Printf("error when saving new Info.json: %s", err)
        return err
    }

    //reloading not necessary, just changed it live xD

    return nil
}


func (mod_simple_list *ModSimpleList) toggleMod(mod_name string) error {
    var err error

    for index, mod := range mod_simple_list.Mods {
        if mod.Name == mod_name {
            mod_simple_list.Mods[index].Enabled = !mod_simple_list.Mods[index].Enabled
            break
        }
    }

    err = mod_simple_list.saveModInfoJson()
    if err != nil {
        log.Printf("error on savin new ModSimpleList: %s", err)
        return err
    }

    //i changed it already don't need to reload it

    return nil
}
