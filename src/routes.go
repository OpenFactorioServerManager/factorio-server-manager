package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	//TODO Proper origin check
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Handler func(*Client, interface{})

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

type WSRouter struct {
	rules map[string]Handler
}

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	ws := NewWSRouter()

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

	// Route for initializing websocket connection
	// Clients connecting to /ws establish websocket connection by upgrading
	// HTTP session.
	// Ensure user is logged in with the AuthorizeHandler middleware
	r.Path("/ws").
		Methods("GET").
		Name("Websocket").
		Handler(AuthorizeHandler(ws))
	ws.Handle("command send", commandSend)
	ws.Handle("log subscribe", logSubscribe)

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
	r.Path("/console").
		Methods("GET").
		Name("Server").
		Handler(AuthorizeHandler(http.StripPrefix("/console", http.FileServer(http.Dir("./app/")))))
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

func NewWSRouter() *WSRouter {
	return &WSRouter{
		rules: make(map[string]Handler),
	}
}

func (ws *WSRouter) Handle(msgName string, handler Handler) {
	ws.rules[msgName] = handler
}

func (ws *WSRouter) FindHandler(msgName string) (Handler, bool) {
	handler, found := ws.rules[msgName]
	return handler, found
}

func (ws *WSRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error opening ws connection: %s", err)
		return
	}
	client := NewClient(socket, ws.FindHandler)
	defer client.Close()
	go client.Write()
	client.Read()
}

// Defines all API REST endpoints
// All routes are prefixed with /api
var apiRoutes = Routes{
	Route{
		"ListInstalledMods",
		"GET",
		"/mods/list/installed",
		listInstalledModsHandler,
	}, {
		"LoginFactorioModPortal",
		"POST",
		"/mods/factorio/login",
		LoginFactorioModPortal,
	}, {
		"SearchModPortal",
		"GET",
		"/mods/search",
		ModPortalSearchHandler,
	}, {
		"GetModDetails",
		"POST",
		"/mods/details",
		ModPortalDetailsHandler,
	}, {
		"ModPortalInstall",
		"POST",
		"/mods/install",
		ModPortalInstallHandler,
	}, {
		"ToggleMod",
		"POST",
		"/mods/toggle",
		ToggleModHandler,
	}, {
		"DeleteMod",
		"POST",
		"/mods/delete",
		DeleteModHandler,
	}, {
        "UpdateMod",
        "POST",
        "/mods/update",
        UpdateModHandler,
    }, {
		"UploadMod",
		"POST",
		"/mods/upload",
		UploadModHandler,
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
		"KillServer",
		"GET",
		"/server/kill",
		KillServer,
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
		ListModPacksHandler,
	}, {
		"DownloadModPack",
		"GET",
		"/mods/packs/download/{modpack}",
		DownloadModPackHandler,
	}, {
		"DeleteModPack",
		"POST",
		"/mods/packs/delete",
		DeleteModPackHandler,
	}, {
		"CreateModPack",
		"POST",
		"/mods/packs/create",
		CreateModPackHandler,
	}, {
		"LoadModPack",
		"POST",
		"/mods/packs/load",
		LoadModPackHandler,
	}, {
	    "ModPackToggleMod",
	    "POST",
	    "/mods/packs/mod/toggle",
	    ModPackToggleModHandler,
    }, {
    	"ModPackDeleteMod",
    	"POST",
    	"/mods/packs/mod/delete",
    	ModPackDeleteModHandler,
	}, {
		"ModPackUpdateMod",
		"POST",
		"/mods/packs/mod/update",
		ModPackUpdateModHandler,
	} ,{
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
