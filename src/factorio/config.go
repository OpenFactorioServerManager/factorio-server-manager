package factorio

import (
	"log"

	"github.com/go-ini/ini"
)

// Loads bootstrap.ini file from the factorio bootstrap directory
func LoadConfig(filename string) (map[string]map[string]string, error) {
	log.Printf("Loading bootstrap file: %s", filename)
	cfg, err := ini.Load(filename)
	if err != nil {
		log.Printf("Error loading bootstrap.ini file: %s", err)
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
	log.Printf("Encoding bootstrap.ini to JSON")

	return result, nil
}
