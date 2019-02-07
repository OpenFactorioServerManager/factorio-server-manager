package main

import (
	"bytes"
	"encoding/json"
	"github.com/go-test/deep"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestModSettings0_16(t *testing.T) {
	// Read dat and compare to JSON
	file, err := os.Open(filepath.Join("factorio_mod_settings_testfiles", "mod_settings_0.16.dat"))
	if err != nil {
		t.Fatalf("could not open mod-settings.dat: %s", err)
	}

	var modData FModData
	err = modData.Decode(file)
	if err != nil {
		t.Fatalf("could not decode FModData: %s", err)
	}

	modDataJson, err := ioutil.ReadFile(filepath.Join("factorio_mod_settings_testfiles", "mod_settings_0.16.json"))
	if err != nil {
		t.Fatalf("could not read json-file: %s", err)
	}

	var test interface{}
	err = json.Unmarshal(modDataJson, &test)
	if err != nil {
		t.Fatalf("could not Unmarshal JSON: %s", err)
	}

	diff := deep.Equal(modData.Data, test)
	if len(diff) > 0 {
		t.Fatalf("Data has %d differences: %s", len(diff), diff)
	}


	// Change some value
	modData.Data.(map[string]interface{})["runtime-per-user"].(map[string]interface{})["folk-fill-fuel-stack-size"].(map[string]interface{})["value"] = 150
	test.(map[string]interface{})["runtime-per-user"].(map[string]interface{})["folk-fill-fuel-stack-size"].(map[string]interface{})["value"] = float64(150)

	// write new data
	newBytes, err := modData.Encode()
	newBytesReader := bytes.NewReader(newBytes)
	if err != nil {
		t.Fatalf("couldn't Encode modData: %s", err)
	}

	var newData FModData
	err = newData.Decode(newBytesReader)
	if err != nil {
		t.Fatalf("couldn't Decode newBytes: %s", err)
	}

	diff2 := deep.Equal(newData.Data, test)
	if len(diff2) > 0 {
		t.Fatalf("Data has %d differences: %s", len(diff2), diff2)
	}
}
