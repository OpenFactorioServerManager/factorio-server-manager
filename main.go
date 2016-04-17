package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type Config struct {
	FactorioDir string
}

var config Config

func main() {
	factorioDir := flag.String("dir", "./", "Specify location of Factorio config directory.")
	flag.Parse()

	config.FactorioDir = *factorioDir

	fmt.Println(listInstalledMods())

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
