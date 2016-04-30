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

	r.Path("/login").
		Methods("GET").
		Name("Login").
		Handler(http.StripPrefix("/login", http.FileServer(http.Dir("./app/"))))

	// API subrouter
	// Serves all JSON REST handlers prefixed with /api
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
		Handler(CheckSession(http.StripPrefix("/mods", http.FileServer(http.Dir("./app/")))))
	r.Path("/saves").
		Methods("GET").
		Name("Saves").
		Handler(CheckSession(http.StripPrefix("/saves", http.FileServer(http.Dir("./app/")))))
	r.Path("/logs").
		Methods("GET").
		Name("Logs").
		Handler(CheckSession(http.StripPrefix("/logs", http.FileServer(http.Dir("./app/")))))
	r.Path("/config").
		Methods("GET").
		Name("Config").
		Handler(CheckSession(http.StripPrefix("/config", http.FileServer(http.Dir("./app/")))))
	r.Path("/server").
		Methods("GET").
		Name("Server").
		Handler(CheckSession(http.StripPrefix("/server", http.FileServer(http.Dir("./app/")))))
	r.PathPrefix("/").
		Methods("GET").
		Name("Index").
		Handler(http.FileServer(http.Dir("./app/")))

	return r
}

// Middleware returns a http.HandlerFunc which authenticates the users request
func CheckSession(h http.Handler) http.Handler {
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
		"LoginUser",
		"POST",
		"/login",
		LoginUser,
	}, {
		"LogoutUser",
		"GET",
		"/logout",
		LogoutUser,
	}, {
		"UserStatus",
		"GET",
		"/user/status",
		GetCurrentLogin,
	},
}
