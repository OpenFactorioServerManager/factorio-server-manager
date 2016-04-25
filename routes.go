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
		s.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	// Serves the frontend application from the app directory
	// Uses basic file server to serve index.html and Javascript application
	// Routes match the ones defined in React application
	r.Path("/mods").
		Methods("GET").
		Name("Mods").
		Handler(http.StripPrefix("/mods", http.FileServer(http.Dir("./app/"))))
	r.Path("/saves").
		Methods("GET").
		Name("Saves").
		Handler(http.StripPrefix("/saves", http.FileServer(http.Dir("./app/"))))
	r.Path("/logs").
		Methods("GET").
		Name("Logs").
		Handler(http.StripPrefix("/logs", http.FileServer(http.Dir("./app/"))))
	r.Path("/config").
		Methods("GET").
		Name("Config").
		Handler(http.StripPrefix("/config", http.FileServer(http.Dir("./app/"))))
	r.PathPrefix("/").
		Methods("GET").
		Name("Index").
		Handler(http.FileServer(http.Dir("./app/")))

	return r
}

// Defines all API REST endpoints
// All routes are prefixed with /api
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
		"UploadMod",
		"POST",
		"/mods/upload",
		UploadMod,
	}, {
		"RemoveMod",
		"GET",
		"/mods/rm/{mod}",
		RemoveMod,
	}, {
		"DownloadMod",
		"GET",
		"/mods/dl/{mod}",
		DownloadMod,
	}, {
		"ListSaves",
		"GET",
		"/saves/list",
		ListSaves,
	}, {
		"DlSave",
		"GET",
		"/saves/dl/{save}",
		DLSave,
	}, {
		"UploadSave",
		"POST",
		"/saves/upload",
		UploadSave,
	}, {
		"RemoveSave",
		"GET",
		"/saves/rm/{save}",
		RemoveSave,
	}, {
		"CreateSave",
		"GET",
		"/saves/create/{save}",
		CreateSaveHandler,
	}, {
		"LogTail",
		"GET",
		"/log/tail",
		LogTail,
	}, {
		"LoadConfig",
		"GET",
		"/config",
		LoadConfig,
	},
}
