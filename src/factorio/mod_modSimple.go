package factorio

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type ModSimple struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type ModSimpleList struct {
	Mods        []ModSimple `json:"mods"`
	Destination string      `json:"-"`
}

func newModSimpleList(destination string) (ModSimpleList, error) {
	var err error

	modSimpleList := ModSimpleList{
		Destination: destination,
	}

	err = modSimpleList.listInstalledMods()
	if err != nil {
		log.Printf("ModSimpleList ... error list installed mods: %s", err)
		return modSimpleList, err
	}

	return modSimpleList, nil
}

func (modSimpleList *ModSimpleList) listInstalledMods() error {
	var err error

	file, err := ioutil.ReadFile(modSimpleList.Destination + "/mod-list.json")
	if os.IsNotExist(err) {
		log.Println("no mod-list.json found ... create new one ...")

		//mod-list.json does not exist, create one
		baseMod := ModSimple{
			Name:    "base",
			Enabled: true,
		}
		modSimpleList.Mods = append(modSimpleList.Mods, baseMod)
		err = modSimpleList.saveModInfoJson()
		if err != nil {
			log.Printf("error saving mod-list.json: %s", err)
			return err
		}

		return nil
	}
	if err != nil {
		log.Printf("ModSimpleList ... error read or write the mod-info.json: %s", err)
		return err
	}

	err = json.Unmarshal(file, modSimpleList)
	if err != nil {
		log.Printf("ModSimpleList ... error while decode mod-info.json: %s", err)
		return err
	}

	return nil
}

func (modSimpleList *ModSimpleList) saveModInfoJson() error {
	var err error

	//build json of current state
	newJson, _ := json.MarshalIndent(modSimpleList, "", "    ")

	err = ioutil.WriteFile(modSimpleList.Destination+"/mod-list.json", newJson, 0664)
	if err != nil {
		log.Printf("error when writing new mod-list: %s", err)
		return err
	}

	return nil
}

func (modSimpleList *ModSimpleList) deleteMod(modName string) error {
	var err error

	for index, mod := range modSimpleList.Mods {
		if mod.Name == modName {
			slice1 := modSimpleList.Mods[:index]
			slice2 := modSimpleList.Mods[index+1:]
			var newModList []ModSimple
			newModList = append(newModList, slice1...)
			newModList = append(newModList, slice2...)
			modSimpleList.Mods = newModList
			break
		}
	}

	err = modSimpleList.saveModInfoJson()
	if err != nil {
		log.Printf("error when saving new mod_list: %s", err)
		return err
	}

	return nil
}

func (modSimpleList *ModSimpleList) CheckModExists(modName string) bool {
	for _, singleMod := range modSimpleList.Mods {
		if singleMod.Name == modName {
			return true
		}
	}

	return false
}

func (modSimpleList *ModSimpleList) createMod(modName string) error {
	var err error

	newModSimple := ModSimple{
		Name:    modName,
		Enabled: true,
	}

	modSimpleList.Mods = append(modSimpleList.Mods, newModSimple)

	err = modSimpleList.saveModInfoJson()
	if err != nil {
		log.Printf("error when saving new Info.json: %s", err)
		return err
	}

	//reloading not necessary, just changed it live xD

	return nil
}

func (modSimpleList *ModSimpleList) ToggleMod(modName string) (error, bool) {
	var err error
	var newEnabled bool

	var found bool
	for index, mod := range modSimpleList.Mods {
		if mod.Name == modName {
			newEnabled = !modSimpleList.Mods[index].Enabled
			modSimpleList.Mods[index].Enabled = newEnabled
			found = true
			break
		}
	}

	if !found {
		return errors.New("mod is not installed"), newEnabled
	}

	err = modSimpleList.saveModInfoJson()
	if err != nil {
		log.Printf("error on savin new ModSimpleList: %s", err)
		return err, newEnabled
	}

	//i changed it already don't need to reload it

	return nil, newEnabled
}
