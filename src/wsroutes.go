package main

import (
	"log"

	"github.com/hpcloud/tail"
)

func logSubscribe(client *Client, data interface{}) {
	t, err := tail.TailFile(config.FactorioLog, tail.Config{Follow: true})
	if err != nil {
		log.Printf("Error subscribing to tail log %s", err)
		return
	}

	for line := range t.Lines {
		client.send <- Message{"log update", line.Text}
	}
}
