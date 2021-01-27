package factorio

import (
	"log"
	"strconv"

	"github.com/OpenFactorioServerManager/factorio-server-manager/bootstrap"

	"github.com/OpenFactorioServerManager/rcon"
)

func connectRC() error {
	var err error
	config := bootstrap.GetConfig()
	rconAddr := config.ServerIP + ":" + strconv.Itoa(config.FactorioRconPort)
	server := GetFactorioServer()
	server.Rcon, err = rcon.Dial(rconAddr, config.FactorioRconPass)
	if err != nil {
		log.Printf("Cannot create rcon session: %s", err)
		return err
	}
	log.Printf("rcon session established on %s", rconAddr)

	return nil
}
