package main

import (
	"flag"
	"log"
	"net/http"
)

type Config struct {
	FactorioDir string
	ServerIP    string
	FactorioLog string
}

var config Config

func main() {
	factorioDir := flag.String("dir", "./", "Specify location of Factorio config directory.")
	factorioIP := flag.String("host", "0.0.0.0:8080", "Specify IP and port for webserver to listen on.")
	flag.Parse()

	config.FactorioDir = *factorioDir
	config.ServerIP = *factorioIP
	config.FactorioLog = config.FactorioDir + "/factorio-current.log"

	router := NewRouter()

	log.Fatal(http.ListenAndServe(config.ServerIP, router))
}
