package bootstrap

import (
	"encoding/json"
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
)

type Flags struct {
	ConfFile           string `long:"conf" default:"./conf.json" description:"Specify location of Factorio Server Manager config file."`
	FactorioDir        string `long:"dir" default:"./" description:"Specify location of Factorio directory."`
	ServerIP           string `long:"host" default:"0.0.0.0" description:"Specify IP for webserver to listen on."`
	FactorioIP         string `long:"game-bind-address" default:"0.0.0.0" description:"Specify IP for Fcatorio gamer server to listen on."`
	FactorioPort       string `long:"port" default:"8080" description:"Specify a port for the server."`
	FactorioConfigFile string `long:"config" default:"config/config.ini" description:"Specify location of Factorio config.ini file"`
	FactorioMaxUpload  int64  `long:"max-upload" default:"20.971.520" description:"Maximum filesize for uploaded files (default 20MB)."`
	FactorioBinary     string `long:"bin" default:"bin/x64/factorio" description:"Location of Factorio Server binary file"`
	GlibcCustom        string `long:"glibc-custom" default:"false" description:"By default false, if custom glibc is required set this to true and add glibc-loc and glibc-lib-loc parameters"`
	GlibcLocation      string `long:"glibc-loc" default:"/opt/glibc-2.18/lib/ld-2.18.so" description:"Location glibc ld.so file if needed (ex. /opt/glibc-2.18/lib/ld-2.18.so)"`
	GlibcLibLoc        string `long:"glibc-lib-loc" default:"/opt/glibc-2.18/lib" description:"Location of glibc lib folder (ex. /opt/glibc-2.18/lib)"`
	Autostart          string `long:"autostart" default:"false" description:"Autostart factorio server on bootup of FSM, default false [true/false]"`
	ModPackDir         string `long:"mod-pack-dir" default:"./mod_packs" description:"Directory to store mod packs."`
}

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
	GlibcCustom             string
	GlibcLocation           string
	GlibcLibLoc             string
	Autostart               string
	ConsoleCacheSize        int `json:"console_cache_size"` // the amount of cached lines, inside the factorio output cache
}

var instantiated Config

func NewConfig(args []string) Config {
	var opts Flags
	_, err := flags.NewParser(&opts, flags.IgnoreUnknown).ParseArgs(args)
	if err != nil {
		failOnError(err, "Failed to parse arguments")
	}
	instantiated = mapFlags(opts)
	instantiated.loadServerConfig()

	abs, err := filepath.Abs(instantiated.FactorioModPackDir)
	println(abs)

	return instantiated
}

func GetConfig() Config {
	return instantiated
}

// Loads server configuration files
// JSON config file contains default values,
// config file will overwrite any provided flags
func (config *Config) loadServerConfig() {
	file, err := os.Open(config.ConfFile)
	failOnError(err, "Error loading config file.")

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	failOnError(err, "Error decoding JSON config file.")

	config.FactorioRconPort = randomPort()
}

// Returns random port to use for rcon connection
func randomPort() int {
	return rand.Intn(45000-40000) + 40000
}

func mapFlags(flags Flags) Config {
	var config = Config{
		Autostart:               flags.Autostart,
		GlibcCustom:             flags.GlibcCustom,
		GlibcLocation:           flags.GlibcLocation,
		GlibcLibLoc:             flags.GlibcLibLoc,
		ConfFile:                flags.ConfFile,
		FactorioDir:             flags.FactorioDir,
		ServerIP:                flags.ServerIP,
		ServerPort:              flags.FactorioPort,
		FactorioIP:              flags.FactorioIP,
		FactorioSavesDir:        filepath.Join(flags.FactorioDir, "saves"),
		FactorioModsDir:         filepath.Join(flags.FactorioDir, "mods"),
		FactorioModPackDir:      flags.ModPackDir,
		FactorioConfigDir:       filepath.Join(flags.FactorioDir, "config"),
		FactorioConfigFile:      filepath.Join(flags.FactorioDir, flags.FactorioConfigFile),
		FactorioBinary:          filepath.Join(flags.FactorioDir, flags.FactorioBinary),
		FactorioCredentialsFile: "./factorio.auth",
		FactorioAdminFile:       "server-adminlist.json",
		MaxUploadSize:           flags.FactorioMaxUpload,
	}

	if runtime.GOOS == "windows" {
		appdata := os.Getenv("APPDATA")
		config.FactorioLog = filepath.Join(appdata, "Factorio", "factorio-current.log")
	} else {
		config.FactorioLog = filepath.Join(config.FactorioDir, "factorio-current.log")
	}

	return config
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
