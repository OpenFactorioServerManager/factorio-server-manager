package main

import (
	"log"
	"strconv"

	"github.com/majormjr/rcon"
)

func connectRC() error {
	var err error
	rconAddr := config.ServerIP + ":" + strconv.Itoa(config.FactorioRconPort)
	FactorioServ.Rcon, err = rcon.Dial(rconAddr, config.FactorioRconPass)
	if err != nil {
		log.Printf("Cannot create rcon session: %s", err)
		return err
	}
	log.Printf("rcon session established on %s", rconAddr)

	return nil
}
