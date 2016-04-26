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
	FactorioBinary     string
	ServerIP           string
	ServerPort         string
	MaxUploadSize      int64
}

var (
	config       Config
	FactorioServ *FactorioServer
)

func loadFlags() {
	factorioDir := flag.String("dir", "./", "Specify location of Factorio directory.")
	factorioIP := flag.String("host", "0.0.0.0", "Specify IP for webserver to listen on.")
	factorioPort := flag.String("port", "8080", "Specify a port for the server.")
	factorioConfigFile := flag.String("config", "config/config.ini", "Specify location of Factorio config.ini file")
	factorioMaxUpload := flag.Int64("max-upload", 1024*1024*20, "Maximum filesize for uploaded files (default 20MB).")
	factorioBinary := flag.String("bin", "bin/x64/factorio", "Location of Factorio Server binary file")

	flag.Parse()

	config.FactorioDir = *factorioDir
	config.FactorioSavesDir = config.FactorioDir + "/saves"
	config.FactorioModsDir = config.FactorioDir + "/mods"
	config.FactorioConfigFile = config.FactorioDir + "/" + *factorioConfigFile
	config.FactorioBinary = config.FactorioDir + "/" + *factorioBinary
	config.ServerIP = *factorioIP
	config.ServerPort = *factorioPort
	config.FactorioLog = config.FactorioDir + "/factorio-current.log"
	config.MaxUploadSize = *factorioMaxUpload
}

func main() {
	loadFlags()

	FactorioServ = initFactorio()
	FactorioServ.Port = 12345
	FactorioServ.Savefile = "testingsaves"

	router := NewRouter()

	log.Fatal(http.ListenAndServe(config.ServerIP+":"+config.ServerPort, router))
}
