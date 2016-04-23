package main

import (
	"encoding/json"
	"log"

	"github.com/go-ini/ini"
)

func loadConfig(filename string) ([]byte, error) {
	log.Printf("Loading config file: %s", filename)
	cfg, err := ini.Load(filename)
	if err != nil {
		log.Printf("Error loading config.ini file: %s", err)
		return nil, err
	}

	result := map[string]map[string]string{}

	sections := cfg.Sections()
	sectionNames := cfg.SectionStrings()
	log.Printf("Appending sections %s to JSON response", sectionNames)
	for _, s := range sections {
		sectionName := s.Name()
		if sectionName == "DEFAULT" {
			continue
		}
		result[sectionName] = map[string]string{}
		result[sectionName] = s.KeysHash()
	}
	log.Printf("Encoding config.ini to JSON")
	resp, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error marshaling config.ini to JSON: %s", err)
		return nil, err
	}

	return resp, nil
}
