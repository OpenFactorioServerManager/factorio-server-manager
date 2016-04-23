package main

import (
	"flag"
	"log"
	"net/http"
)

type Config struct {
	FactorioDir        string
	FactorioSavesDir   string
	FactorioModsDir    string
	FactorioConfigFile string
	FactorioLog        string
	ServerIP           string
	ServerPort         string
	MaxUploadSize      int64
}

var config Config

func loadFlags() {
	factorioDir := flag.String("dir", "./", "Specify location of Factorio directory.")
	factorioIP := flag.String("host", "0.0.0.0", "Specify IP for webserver to listen on.")
	factorioPort := flag.String("port", "8080", "Specify a port for the server.")
	factorioConfigFile := flag.String("config", "config/config.ini", "Specify location of Factorio config.ini file")
	factorioMaxUpload := flag.Int64("max-upload", 100000, "Maximum filesize for uploaded files.")

	flag.Parse()

	config.FactorioDir = *factorioDir
	config.FactorioSavesDir = config.FactorioDir + "/saves"
	config.FactorioModsDir = config.FactorioDir + "/mods"
	config.FactorioConfigFile = config.FactorioDir + "/" + *factorioConfigFile
	config.ServerIP = *factorioIP
	config.ServerPort = *factorioPort
	config.FactorioLog = config.FactorioDir + "/factorio-current.log"
	config.MaxUploadSize = *factorioMaxUpload
}

func main() {
	loadFlags()

	router := NewRouter()

	log.Fatal(http.ListenAndServe(config.ServerIP+":"+config.ServerPort, router))
}
