package main

import (
	"log"
	"path/filepath"

	"github.com/hpcloud/tail"
)

func IsClosed(ch <-chan Message) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func logSubscribe(client *Client, data interface{}) {
	go func() {
		logfile := filepath.Join(config.FactorioDir, "factorio-server-console.log")
		t, err := tail.TailFile(logfile, tail.Config{Follow: true})
		if err != nil {
			log.Printf("Error subscribing to tail log %s", err)
			return
		}

		for line := range t.Lines {
			if !IsClosed(client.send) {
				client.send <- Message{"log update", line.Text}
			} else {
				log.Printf("Channel was closed")
				return
			}
		}
	}()
}

func commandSend(client *Client, data interface{}) {
	if FactorioServ.Running {
		go func() {
			log.Printf("Received command: %v", data)

			reqId, err := FactorioServ.Rcon.Write(data.(string))
			if err != nil {
				log.Printf("Error sending rcon command: %s", err)
				return
			}

			log.Printf("Command send to Factorio: %s, with rcon request id: %v", data, reqId)

			if !IsClosed(client.send) {
				client.send <- Message{"receive command", data}
			} else {
				log.Printf("Channel was closed")
				return
			}
		}()
	}
}
