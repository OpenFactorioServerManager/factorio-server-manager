package factorio

import (
	"github.com/mroote/factorio-server-manager/bootstrap"
	"log"
	"strconv"

	"github.com/majormjr/rcon"
)

func connectRC() error {
	var err error
	config := bootstrap.GetConfig()
	rconAddr := config.ServerIP + ":" + strconv.Itoa(config.FactorioRconPort)
	server, err := GetFactorioServer()
	server.Rcon, err = rcon.Dial(rconAddr, config.FactorioRconPass)
	if err != nil {
		log.Printf("Cannot create rcon session: %s", err)
		return err
	}
	log.Printf("rcon session established on %s", rconAddr)

	return nil
}
