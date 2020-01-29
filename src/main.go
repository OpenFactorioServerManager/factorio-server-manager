package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	FactorioDir             string `json:"factorio_dir"`
	FactorioSavesDir        string `json:"saves_dir"`
	FactorioModsDir         string `json:"mods_dir"`
	FactorioModPackDir      string `json:"mod_pack_dir"`
	FactorioConfigFile      string `json:"config_file"`
	FactorioConfigDir       string `json:"config_directory"`
	FactorioLog             string `json:"logfile"`
	FactorioBinary          string `json:"factorio_binary"`
	FactorioRconPort        int    `json:"rcon_port"`
	FactorioRconPass        string `json:"rcon_pass"`
	FactorioCredentialsFile string `json:"factorio_credentials_file"`
	FactorioIP              string `json:"factorio_ip"`
	FactorioAdminFile       string `json:"-"`
	ServerIP                string `json:"server_ip"`
	ServerPort              string `json:"server_port"`
	MaxUploadSize           int64  `json:"max_upload_size"`
	Username                string `json:"username"`
	Password                string `json:"password"`
	DatabaseFile            string `json:"database_file"`
	CookieEncryptionKey     string `json:"cookie_encryption_key"`
	SettingsFile            string `json:"settings_file"`
	LogFile                 string `json:"log_file"`
	ConfFile                string
	glibcCustom             string
	glibcLocation           string
	glibcLibLoc             string
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
// config file will overwrite any provided flags
func loadServerConfig(f string) {
	file, err := os.Open(f)
	failOnError(err, "Error loading config file.")

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	failOnError(err, "Error decoding JSON config file.")

	config.FactorioRconPort = randomPort()
}

func parseFlags() {
	confFile := flag.String("conf", "./conf.json", "Specify location of Factorio Server Manager config file.")
	factorioDir := flag.String("dir", "./", "Specify location of Factorio directory.")
	serverIP := flag.String("host", "0.0.0.0", "Specify IP for webserver to listen on.")
	factorioIP := flag.String("game-bind-address", "0.0.0.0", "Specify IP for Fcatorio gamer server to listen on.")
	factorioPort := flag.String("port", "8080", "Specify a port for the server.")
	factorioConfigFile := flag.String("config", "config/config.ini", "Specify location of Factorio config.ini file")
	factorioMaxUpload := flag.Int64("max-upload", 1024*1024*20, "Maximum filesize for uploaded files (default 20MB).")
	factorioBinary := flag.String("bin", "bin/x64/factorio", "Location of Factorio Server binary file")
	glibcCustom := flag.String("glibc-custom", "false", "By default false, if custom glibc is required set this to true and add glibc-loc and glibc-lib-loc parameters")
	glibcLocation := flag.String("glibc-loc", "/opt/glibc-2.18/lib/ld-2.18.so", "Location glibc ld.so file if needed (ex. /opt/glibc-2.18/lib/ld-2.18.so)")
	glibcLibLoc := flag.String("glibc-lib-loc", "/opt/glibc-2.18/lib", "Location of glibc lib folder (ex. /opt/glibc-2.18/lib)")

	flag.Parse()
	config.glibcCustom = *glibcCustom
	config.glibcLocation = *glibcLocation
	config.glibcLibLoc = *glibcLibLoc
	config.ConfFile = *confFile
	config.FactorioDir = *factorioDir
	config.ServerIP = *serverIP
	config.FactorioIP = *factorioIP
	config.ServerPort = *factorioPort
	config.FactorioSavesDir = filepath.Join(config.FactorioDir, "saves")
	config.FactorioModsDir = filepath.Join(config.FactorioDir, "mods")
	config.FactorioModPackDir = "./mod_packs"
	config.FactorioConfigDir = filepath.Join(config.FactorioDir, "config")
	config.FactorioConfigFile = filepath.Join(config.FactorioDir, *factorioConfigFile)
	config.FactorioBinary = filepath.Join(config.FactorioDir, *factorioBinary)
	config.FactorioCredentialsFile = "./factorio.auth"
	config.FactorioAdminFile = "server-adminlist.json"
	config.MaxUploadSize = *factorioMaxUpload

	if runtime.GOOS == "windows" {
		appdata := os.Getenv("APPDATA")
		config.FactorioLog = filepath.Join(appdata, "Factorio", "factorio-current.log")
	} else {
		config.FactorioLog = filepath.Join(config.FactorioDir, "factorio-current.log")
	}
}

func main() {
	var err error

	// Parse configuration flags
	parseFlags()
	// Load server config from file
	loadServerConfig(config.ConfFile)
	// create mod-stuff
	modStartUp()

	// Initialize Factorio Server struct
	FactorioServ, err = initFactorio()
	if err != nil {
		log.Printf("Error occurred during FactorioServer initializaion: %v\n", err)
		return
	}

	// Initialize authentication system
	Auth = initAuth()
	Auth.CreateAuth(config.DatabaseFile, config.CookieEncryptionKey)
	Auth.CreateOrUpdateUser(config.Username, config.Password, "admin", "")

	// Initialize HTTP router
	router := NewRouter()
	log.Printf("Starting server on: %s:%s", config.ServerIP, config.ServerPort)
	log.Fatal(http.ListenAndServe(config.ServerIP+":"+config.ServerPort, router))

}
