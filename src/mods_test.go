package main

import (
    "testing"
    "fmt"
)
/*
 *JUST FOR TESTING MY IMPLEMENTATION, NOT FOR TESTING THE PROGRAM
 */

//func TestModsListInstalledMods(t *testing.T) {
//    var err error
//    mods, err := newMods("../../factorio-server/mods")
//    if err != nil {
//        t.Errorf("failed while creating newMods: %s", err)
//    }
//
//    mods_result_list := mods.listInstalledMods()
//
//    by, err := json.MarshalIndent(mods_result_list, "", "  ")
//    if err != nil {
//        t.Error("error while marshaling the result")
//    }
//
//    //fmt.Print(by)
//    os.Stdout.Write(by)
//}

//func TestNewModPackList(t *testing.T) {
//    var err error
//    mod_pack_map, err := newModPackMap()
//    if err != nil {
//        t.Errorf("error: %s", err)
//    }
//    if mod_pack_map == nil {
//        fmt.Println("mod_pack_map is nil :(")
//    } else {
//        fmt.Println("mod_pack_map is not nil :)")
//    }
//
//    err = mod_pack_map.createModPack("baum")
//    if err != nil {
//        t.Errorf("error on creating. %s", err)
//    }
//    fmt.Println(mod_pack_map)
//}
