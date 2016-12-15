package main

import (
	"fmt"
	"log"

	"github.com/james4k/rcon"
)

func connectRC(addr, pass string) {
	rc, err := rcon.Dial(addr, pass)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	defer rc.Close()

	req_id, err := rc.Write("Factorio Server Manager Connected")
	if err != nil {
		log.Printf("Error: %s", err)
	}

	log.Printf("Establish RCON connection with request id: %s", req_id)

	for {
		resp, req_id, err := rc.Read()
		if err != nil {
			log.Printf("Error: %s", err)
		}

		fmt.Println(resp, req_id)
	}

}
