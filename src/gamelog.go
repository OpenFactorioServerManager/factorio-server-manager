package main

import (
	"log"

	"github.com/hpcloud/tail"
)

func tailLog(filename string) ([]string, error) {
	result := []string{}

	t, err := tail.TailFile(config.FactorioLog, tail.Config{Follow: false})
	if err != nil {
		log.Printf("Error tailing log %s", err)
		return result, err
	}

	for line := range t.Lines {
		result = append(result, line.Text)
	}

	return result, nil
}
