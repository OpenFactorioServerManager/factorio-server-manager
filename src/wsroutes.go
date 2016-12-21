package main

import (
	"log"

	"github.com/hpcloud/tail"
)

func logSubscribe(client *Client, data interface{}) {
	go func() {
		t, err := tail.TailFile(config.FactorioLog, tail.Config{Follow: true})
		if err != nil {
			log.Printf("Error subscribing to tail log %s", err)
			return
		}

		for line := range t.Lines {
			client.send <- Message{"log update", line.Text}
		}
	}()
}

func commandSend(client *Client, data interface{}) {
	go func() {
		log.Printf("Received command: %v", data)

		req_id, err := FactorioServ.Rcon.Write(data.(string))
		if err != nil {
			log.Printf("Error sending rcon command: %s", err)
			return
		}

		log.Printf("Command send to Factorio: %s, with rcon request id: %v", data, req_id)

		client.send <- Message{"receive command", data}
	}()
}
