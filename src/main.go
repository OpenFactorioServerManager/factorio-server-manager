package main

import (
	"github.com/mroote/factorio-server-manager/api"
	"github.com/mroote/factorio-server-manager/bootstrap"
	"github.com/mroote/factorio-server-manager/factorio"
	"log"
	"net/http"
	"os"
)

func main() {
	// get the all configs based on the flags
	config := bootstrap.NewConfig(os.Args)

	// setup the required files for the mods
	factorio.ModStartUp()

	// Initialize Factorio Server struct
	_, err := factorio.GetFactorioServer()
	if err != nil {
		log.Printf("Error occurred during Server initializaion: %v\n", err)
		return
	}

	// Initialize authentication system
	api.GetAuth()

	// Initialize HTTP router
	router := api.NewRouter()
	log.Printf("Starting server on: %s:%s", config.ServerIP, config.ServerPort)
	log.Fatal(http.ListenAndServe(config.ServerIP+":"+config.ServerPort, router))

}
