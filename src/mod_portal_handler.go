package main

import (
	"fmt"
	"log"
	"net/http"
)

func FactorioModPortalListMods(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	var statusCode int
	resp, err, statusCode = listModPortal()

	if err != nil {
		resp = fmt.Sprintf("Error in listing mods from mod portal: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
}
