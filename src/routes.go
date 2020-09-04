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
	ws.Handle("server status subscribe", serverStatusSubscribe)

	// Serves the frontend application from the app directory
	// Uses basic file server to serve index.html and Javascript application
	// Routes match the ones defined in React frontend application
	r.Path("/login").
		Methods("GET").
		Name("Login").
		Handler(http.StripPrefix("/login", http.FileServer(http.Dir("./app/"))))
	r.Path("/saves").
		Methods("GET").
		Name("Saves").
		Handler(AuthorizeHandler(http.StripPrefix("/saves", http.FileServer(http.Dir("./app/")))))
	r.Path("/mods").
		Methods("GET").
		Name("Mods").
		Handler(AuthorizeHandler(http.StripPrefix("/mods", http.FileServer(http.Dir("./app/")))))
	r.Path("/server-settings").
		Methods("GET").
		Name("Server settings").
		Handler(AuthorizeHandler(http.StripPrefix("/server-settings", http.FileServer(http.Dir("./app/")))))
	r.Path("/game-settings").
		Methods("GET").
		Name("Game settings").
		Handler(AuthorizeHandler(http.StripPrefix("/game-settings", http.FileServer(http.Dir("./app/")))))
	r.Path("/console").
		Methods("GET").
		Name("Console").
		Handler(AuthorizeHandler(http.StripPrefix("/console", http.FileServer(http.Dir("./app/")))))
	r.Path("/logs").
		Methods("GET").
		Name("Logs").
		Handler(AuthorizeHandler(http.StripPrefix("/logs", http.FileServer(http.Dir("./app/")))))
	r.Path("/user-management").
		Methods("GET").
		Name("User management").
		Handler(AuthorizeHandler(http.StripPrefix("/user-management", http.FileServer(http.Dir("./app/")))))
	r.Path("/help").
		Methods("GET").
		Name("Help").
		Handler(AuthorizeHandler(http.StripPrefix("/help", http.FileServer(http.Dir("./app/")))))

	// catch all route
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
		"DeleteAllMods",
		"POST",
		"/mods/delete/all",
		DeleteAllModsHandler,
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
		"DownloadMods",
		"GET",
		"/mods/download",
		DownloadModsHandler,
	}, {
		"LoadModsFromSave",
		"POST",
		"/mods/save/load",
		LoadModsFromSaveHandler,
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
		"FactorioVersion",
		"GET",
		"/server/facVersion",
		FactorioVersion,
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
	// Mod Portal Stuff
	{
		"ModPortalListAllMods",
		"GET",
		"/mods/portal/list",
		FactorioModPortalListModsHandler,
	}, {
		"ModPortalGetModInfo",
		"GET",
		"/mods/portal/info/{mod}",
		FactorioModPortalModInfoHandler,
	}, {
		"ModPortalInstallMod",
		"POST",
		"/mods/portal/install",
		FactorioModPortalInstallHandler,
	}, {
		"ModPortalLogin",
		"POST",
		"/mods/portal/login",
		FactorioModPortalLoginHandler,
	}, {
		"ModPortalLoginStatus",
		"POST",
		"/mods/portal/loginstatus",
		FactorioModPortalLoginStatusHandler,
	}, {
		"ModPortalLogout",
		"GET",
		"/mods/portal/logout",
		FactorioModPortalLogoutHandler,
	}, {
		"ModPortalInstallMultiple",
		"POST",
		"/mods/portal/install/multiple",
		ModPortalInstallMultipleHandler,
	},
}
