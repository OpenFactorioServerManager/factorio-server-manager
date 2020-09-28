package main

import (
	"flag"
	"github.com/mroote/factorio-server-manager/api"
	"github.com/mroote/factorio-server-manager/bootstrap"
	"github.com/mroote/factorio-server-manager/factorio"
	"log"
	"net/http"
)

func main() {
	// parse command flags and pass them to bootstrap
	flags := parseFlags()
	bootstrap.SetFlags(flags)

	// get the all configs based on the flags
	config := bootstrap.GetConfig()

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

func parseFlags() bootstrap.Flags {
	confFile := flag.String("conf", "./conf.json", "Specify location of Factorio Server Manager bootstrap file.")
	factorioDir := flag.String("dir", "./", "Specify location of Factorio directory.")
	serverIP := flag.String("host", "0.0.0.0", "Specify IP for webserver to listen on.")
	factorioIP := flag.String("game-bind-address", "0.0.0.0", "Specify IP for Fcatorio gamer server to listen on.")
	factorioPort := flag.String("port", "8080", "Specify a port for the server.")
	factorioConfigFile := flag.String("bootstrap", "bootstrap/bootstrap.ini", "Specify location of Factorio bootstrap.ini file")
	factorioMaxUpload := flag.Int64("max-upload", 1024*1024*20, "Maximum filesize for uploaded files (default 20MB).")
	factorioBinary := flag.String("bin", "bin/x64/factorio", "Location of Factorio Server binary file")
	glibcCustom := flag.String("glibc-custom", "false", "By default false, if custom glibc is required set this to true and add glibc-loc and glibc-lib-loc parameters")
	glibcLocation := flag.String("glibc-loc", "/opt/glibc-2.18/lib/ld-2.18.so", "Location glibc ld.so file if needed (ex. /opt/glibc-2.18/lib/ld-2.18.so)")
	glibcLibLoc := flag.String("glibc-lib-loc", "/opt/glibc-2.18/lib", "Location of glibc lib folder (ex. /opt/glibc-2.18/lib)")
	autostart := flag.String("autostart", "false", "Autostart factorio server on bootup of FSM, default false [true/false]")
	flag.Parse()

	return bootstrap.Flags{
		ConfFile:           confFile,
		FactorioDir:        factorioDir,
		ServerIP:           serverIP,
		FactorioIP:         factorioIP,
		FactorioPort:       factorioPort,
		FactorioConfigFile: factorioConfigFile,
		FactorioMaxUpload:  factorioMaxUpload,
		FactorioBinary:     factorioBinary,
		GlibcCustom:        glibcCustom,
		GlibcLocation:      glibcLocation,
		GlibcLibLoc:        glibcLibLoc,
		Autostart:          autostart,
	}
}
