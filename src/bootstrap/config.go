package bootstrap

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type Flags struct {
	ConfFile           *string
	FactorioDir        *string
	ServerIP           *string
	FactorioIP         *string
	FactorioPort       *string
	FactorioConfigFile *string
	FactorioMaxUpload  *int64
	FactorioBinary     *string
	GlibcCustom        *string
	GlibcLocation      *string
	GlibcLibLoc        *string
	Autostart          *string
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
}

var instantiated Config
var flags Flags
var once sync.Once

func GetConfig() Config {
	once.Do(func() {
		instantiated = mapFlags()
		instantiated.loadServerConfig()
	})
	return instantiated
}

func SetFlags(parsedFlags Flags) {
	flags = parsedFlags
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

func mapFlags() Config {
	var config = Config{
		Autostart:               *flags.Autostart,
		GlibcCustom:             *flags.GlibcCustom,
		GlibcLocation:           *flags.GlibcLocation,
		GlibcLibLoc:             *flags.GlibcLibLoc,
		ConfFile:                *flags.ConfFile,
		FactorioDir:             *flags.FactorioDir,
		ServerIP:                *flags.ServerIP,
		ServerPort:              *flags.FactorioPort,
		FactorioIP:              *flags.FactorioIP,
		FactorioSavesDir:        filepath.Join(*flags.FactorioDir, "saves"),
		FactorioModsDir:         filepath.Join(*flags.FactorioDir, "mods"),
		FactorioModPackDir:      "./mod_packs",
		FactorioConfigDir:       filepath.Join(*flags.FactorioDir, "config"),
		FactorioConfigFile:      filepath.Join(*flags.FactorioDir, *flags.FactorioConfigFile),
		FactorioBinary:          filepath.Join(*flags.FactorioDir, *flags.FactorioBinary),
		FactorioCredentialsFile: "./factorio.auth",
		FactorioAdminFile:       "server-adminlist.json",
		MaxUploadSize:           *flags.FactorioMaxUpload,
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
