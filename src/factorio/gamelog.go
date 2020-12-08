package factorio

import (
	"github.com/hpcloud/tail"
	"github.com/mroote/factorio-server-manager/bootstrap"
	"log"
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
