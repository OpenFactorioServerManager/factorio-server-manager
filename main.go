package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Config struct {
	FactorioDir        string `json:"factorio_dir"`
	FactorioSavesDir   string `json:"saves_dir"`
	FactorioModsDir    string `json:"mods_dir"`
	FactorioConfigFile string `json:"config_file"`
	FactorioLog        string `json:"logfile"`
	FactorioBinary     string `json:"factorio_binary"`
	ServerIP           string `json:"server_ip"`
	ServerPort         string `json:"server_port"`
	MaxUploadSize      int64  `json:"max_upload_size"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	DatabaseFile       string `json:"database_file"`
}

var (
	config       Config
	FactorioServ *FactorioServer
	Auth         *AuthHTTP
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// Loads server configuration files
// JSON config file contains default values,
// any provided flags to the binary will overwrite values in config file.
func loadServerConfig(f string) {
	file, err := os.Open(f)
	failOnError(err, "Error loading config file.")

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	parseFlags()
}

func parseFlags() {
	factorioDir := flag.String("dir", "./", "Specify location of Factorio directory.")
	factorioIP := flag.String("host", "0.0.0.0", "Specify IP for webserver to listen on.")
	factorioPort := flag.String("port", "8080", "Specify a port for the server.")
	factorioConfigFile := flag.String("config", "config/config.ini", "Specify location of Factorio config.ini file")
	factorioMaxUpload := flag.Int64("max-upload", 1024*1024*20, "Maximum filesize for uploaded files (default 20MB).")
	factorioBinary := flag.String("bin", "bin/x64/factorio", "Location of Factorio Server binary file")

	flag.Parse()

	config.FactorioDir = *factorioDir
	config.ServerIP = *factorioIP
	config.ServerPort = *factorioPort
	config.FactorioSavesDir = config.FactorioDir + "/saves"
	config.FactorioModsDir = config.FactorioDir + "/mods"
	config.FactorioConfigFile = config.FactorioDir + "/" + *factorioConfigFile
	config.FactorioBinary = config.FactorioDir + "/" + *factorioBinary
	config.FactorioLog = config.FactorioDir + "/factorio-current.log"
	config.MaxUploadSize = *factorioMaxUpload
}

func main() {
	loadServerConfig("./conf.json")

	FactorioServ = initFactorio()

	Auth = initAuth()
	Auth.createAuthDb(config.DatabaseFile)
	Auth.createRoles()
	err := Auth.createInitialUser(config.Username, config.Password, "admin", "")
	if err != nil {
		log.Printf("Error creating user: %s", err)
	}

	router := NewRouter()

	fmt.Printf("Starting server on: %s:%s", config.ServerIP, config.ServerPort)
	log.Fatal(http.ListenAndServe(config.ServerIP+":"+config.ServerPort, router))
}
