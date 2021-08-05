package factorio

import (
	"log"

	"github.com/OpenFactorioServerManager/factorio-server-manager/bootstrap"
	"github.com/hpcloud/tail"
)

func TailLog() ([]string, error) {
	result := []string{}

	config := bootstrap.GetConfig()

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
