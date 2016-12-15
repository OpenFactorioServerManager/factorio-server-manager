package main

import (
	"log"
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
	// Serves all JSON REST handlers prefixed with /api
	s := r.PathPrefix("/api").Subrouter()
	for _, route := range apiRoutes {
		s.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(AuthorizeHandler(route.HandlerFunc))
	}

	// The login handler does not check for authentication.
	s.Path("/login").
		Methods("POST").
		Name("LoginUser").
		HandlerFunc(LoginUser)

	// Serves the frontend application from the app directory
	// Uses basic file server to serve index.html and Javascript application
	// Routes match the ones defined in React frontend application
	r.Path("/login").
		Methods("GET").
		Name("Login").
		Handler(http.StripPrefix("/login", http.FileServer(http.Dir("./app/"))))
	r.Path("/settings").
		Methods("GET").
		Name("Settings").
		Handler(AuthorizeHandler(http.StripPrefix("/settings", http.FileServer(http.Dir("./app/")))))
	r.Path("/mods").
		Methods("GET").
		Name("Mods").
		Handler(AuthorizeHandler(http.StripPrefix("/mods", http.FileServer(http.Dir("./app/")))))
	r.Path("/saves").
		Methods("GET").
		Name("Saves").
		Handler(AuthorizeHandler(http.StripPrefix("/saves", http.FileServer(http.Dir("./app/")))))
	r.Path("/logs").
		Methods("GET").
		Name("Logs").
		Handler(AuthorizeHandler(http.StripPrefix("/logs", http.FileServer(http.Dir("./app/")))))
	r.Path("/config").
		Methods("GET").
		Name("Config").
		Handler(AuthorizeHandler(http.StripPrefix("/config", http.FileServer(http.Dir("./app/")))))
	r.Path("/server").
		Methods("GET").
		Name("Server").
		Handler(AuthorizeHandler(http.StripPrefix("/server", http.FileServer(http.Dir("./app/")))))
	r.PathPrefix("/").
		Methods("GET").
		Name("Index").
		Handler(http.FileServer(http.Dir("./app/")))

	return r
}

// Middleware returns a http.HandlerFunc which authenticates the users request
// Redirects user to login page if no session is found
func AuthorizeHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := Auth.aaa.Authorize(w, r, true); err != nil {
			log.Printf("Unauthenticated request %s %s %s", r.Method, r.Host, r.RequestURI)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		h.ServeHTTP(w, r)
	})
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
	}, {
		"StartServer",
		"GET",
		"/server/start",
		StartServer,
	}, {
		"StartServer",
		"POST",
		"/server/start",
		StartServer,
	}, {
		"StopServer",
		"GET",
		"/server/stop",
		StopServer,
	}, {
		"RunningServer",
		"GET",
		"/server/status",
		CheckServer,
	}, {
		"LogoutUser",
		"GET",
		"/logout",
		LogoutUser,
	}, {
		"StatusUser",
		"GET",
		"/user/status",
		GetCurrentLogin,
	}, {
		"ListUsers",
		"GET",
		"/user/list",
		ListUsers,
	}, {
		"AddUser",
		"POST",
		"/user/add",
		AddUser,
	}, {
		"RemoveUser",
		"POST",
		"/user/remove",
		RemoveUser,
	}, {
		"ListModPacks",
		"GET",
		"/mods/packs/list",
		ListModPacks,
	}, {
		"DownloadModPack",
		"GET",
		"/mods/packs/dl/{modpack}",
		DownloadModPack,
	}, {
		"DeleteModPack",
		"GET",
		"/mods/packs/rm/{modpack}",
		DeleteModPack,
	}, {
		"CreateModPack",
		"POST",
		"/mods/packs/add",
		CreateModPackHandler,
	}, {
		"GetServerSettings",
		"GET",
		"/settings",
		GetServerSettings,
	}, {
		"UpdateServerSettings",
		"POST",
		"/settings/update",
		UpdateServerSettings,
	},
}
