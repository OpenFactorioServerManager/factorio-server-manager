package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestListInstalledMods(t *testing.T) {
	modDir := "./mods"
	match := ""

	os.Mkdir(modDir, 0755)
	ioutil.WriteFile(modDir+"/mods_test.zip", []byte("TEST"), 0644)

	result, err := listInstalledMods(modDir)
	if err != nil {
		t.Errorf("Error listing mods test: %s", err)
	}
	fmt.Println(result)
	if len(result) < 1 {
		fmt.Println("No Mods installed")
		t.Errorf("Error listing mods test: %s", err)
	}
	for _, m := range result {
		if m == "mods_test.zip" {
			match = m
		}
	}

	if match == "" {
		t.Errorf("Error listing mods test: %s", err)
	}

	os.RemoveAll(modDir)
}
