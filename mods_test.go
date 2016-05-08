package main

import (
	"fmt"
	"testing"
)

var (
	modDir = "/home/mitch/bin/factorio/mods"
)

func TestListInstalledMods(t *testing.T) {
	result, err := listInstalledMods(modDir)
	if err != nil {
		t.Errorf("Error listing mods test: %s", err)
	}
	if len(result) < 1 {
		fmt.Println("No Mods installed")
	}
}
