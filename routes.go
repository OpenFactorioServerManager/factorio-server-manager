package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	// API subrouter
	// Serves all REST handlers prefixed with /api
	s := r.PathPrefix("/api").Subrouter()
	for _, route := range apiRoutes {
		s.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	// Serves the frontend application from the app directory
	// Uses basic file server to server index.html and Javascript application
	r.
		PathPrefix("/").
		Methods("GET").
		Name("Index").
		Handler(http.FileServer(http.Dir("./app/")))

	return r
}

var apiRoutes = Routes{
	Route{
		"ListInstalledMods",
		"GET",
		"/mods/list/installed",
		ListInstalledMods,
	}, {
		"ListMods",
		"GET",
		"/mods/list",
		ListMods,
	}, {
		"ToggleMod",
		"GET",
		"/mods/toggle/{mod}",
		ToggleMod,
	}, {
		"ListSaves",
		"GET",
		"/saves/list",
		ListSaves,
	}, {
		"LogTail",
		"GET",
		"/log/tail",
		LogTail,
	},
}
