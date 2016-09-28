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
	Middleware  func(http.Handler) http.Handler
}

type Routes []Route

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	// API subrouter
	// Serves all JSON REST handlers prefixed with /api
	s := r.PathPrefix("/api").Subrouter()
	for _, route := range apiRoutes {
		apiRoute := s.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name)
        if route.Middleware != nil {
		    apiRoute.Handler(route.Middleware(route.HandlerFunc))
        } else {
		    apiRoute.Handler(route.HandlerFunc)
        }
	}

	// Serves the frontend application from the app directory
	// Uses basic file server to serve index.html and Javascript application
	// Routes match the ones defined in React application
	r.Path("/login").
		Methods("GET").
		Name("Login").
		Handler(http.StripPrefix("/login", http.FileServer(http.Dir("./app/"))))
	r.Path("/settings").
		Methods("GET").
		Name("Settings").
		Handler(CheckAuth(http.StripPrefix("/settings", http.FileServer(http.Dir("./app/")))))
	r.Path("/mods").
		Methods("GET").
		Name("Mods").
		Handler(CheckAuth(http.StripPrefix("/mods", http.FileServer(http.Dir("./app/")))))
	r.Path("/saves").
		Methods("GET").
		Name("Saves").
		Handler(CheckAuth(http.StripPrefix("/saves", http.FileServer(http.Dir("./app/")))))
	r.Path("/logs").
		Methods("GET").
		Name("Logs").
		Handler(CheckAuth(http.StripPrefix("/logs", http.FileServer(http.Dir("./app/")))))
	r.Path("/config").
		Methods("GET").
		Name("Config").
		Handler(CheckAuth(http.StripPrefix("/config", http.FileServer(http.Dir("./app/")))))
	r.Path("/server").
		Methods("GET").
		Name("Server").
		Handler(CheckAuth(http.StripPrefix("/server", http.FileServer(http.Dir("./app/")))))
	r.PathPrefix("/").
		Methods("GET").
		Name("Index").
		Handler(http.FileServer(http.Dir("./app/")))

	return r
}

// Middleware returns a http.HandlerFunc which authenticates the users request
// Redirects user to login page if no session is found
func CheckAuth(h http.Handler) http.Handler {
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
        CheckAuth,
	}, {
		"ListMods",
		"GET",
		"/mods/list",
		ListMods,
        CheckAuth,
	}, {
		"ToggleMod",
		"GET",
		"/mods/toggle/{mod}",
		ToggleMod,
        CheckAuth,
	}, {
		"UploadMod",
		"POST",
		"/mods/upload",
		UploadMod,
        CheckAuth,
	}, {
		"RemoveMod",
		"GET",
		"/mods/rm/{mod}",
		RemoveMod,
        CheckAuth,
	}, {
		"DownloadMod",
		"GET",
		"/mods/dl/{mod}",
		DownloadMod,
        CheckAuth,
	}, {
		"ListSaves",
		"GET",
		"/saves/list",
		ListSaves,
        CheckAuth,
	}, {
		"DlSave",
		"GET",
		"/saves/dl/{save}",
		DLSave,
        CheckAuth,
	}, {
		"UploadSave",
		"POST",
		"/saves/upload",
		UploadSave,
        CheckAuth,
	}, {
		"RemoveSave",
		"GET",
		"/saves/rm/{save}",
		RemoveSave,
        CheckAuth,
	}, {
		"CreateSave",
		"GET",
		"/saves/create/{save}",
		CreateSaveHandler,
        CheckAuth,
	}, {
		"LogTail",
		"GET",
		"/log/tail",
		LogTail,
        CheckAuth,
	}, {
		"LoadConfig",
		"GET",
		"/config",
		LoadConfig,
        CheckAuth,
	}, {
		"StartServer",
		"GET",
		"/server/start",
		StartServer,
        CheckAuth,
	}, {
		"StartServer",
		"POST",
		"/server/start",
		StartServer,
        CheckAuth,
	}, {
		"StopServer",
		"GET",
		"/server/stop",
		StopServer,
        CheckAuth,
	}, {
		"RunningServer",
		"GET",
		"/server/status",
		CheckServer,
        CheckAuth,
	}, {
		"LoginUser",
		"POST",
		"/login",
		LoginUser,
        nil,
	}, {
		"LogoutUser",
		"GET",
		"/logout",
		LogoutUser,
        CheckAuth,
	}, {
		"StatusUser",
		"GET",
		"/user/status",
		GetCurrentLogin,
        CheckAuth,
	}, {
		"ListUsers",
		"GET",
		"/user/list",
		ListUsers,
        CheckAuth,
	}, {
		"AddUser",
		"POST",
		"/user/add",
		AddUser,
        CheckAuth,
	}, {
		"RemoveUser",
		"POST",
		"/user/remove",
		RemoveUser,
        CheckAuth,
	}, {
		"ListModPacks",
		"GET",
		"/mods/packs/list",
		ListModPacks,
        CheckAuth,
	}, {
		"DownloadModPack",
		"GET",
		"/mods/packs/dl/{modpack}",
		DownloadModPack,
        CheckAuth,
	}, {
		"DeleteModPack",
		"GET",
		"/mods/packs/rm/{modpack}",
		DeleteModPack,
        CheckAuth,
	}, {
		"CreateModPack",
		"POST",
		"/mods/packs/add",
		CreateModPackHandler,
        CheckAuth,
	},
}
