package api

import (
	"github.com/mroote/factorio-server-manager/api/websocket"
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

	// Route for initializing websocket connection
	// Clients connecting to /ws establish websocket connection by upgrading
	// HTTP session.
	// Ensure user is logged in with the AuthorizeHandler middleware
	r.Path("/ws").
		Methods("GET").
		Name("Websocket").
		Handler(
			AuthorizeHandler(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						websocket.ServeWs(w, r)
					},
				),
			),
		)

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
		Auth := GetAuth()
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
	{
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
		"LoadModsFromSave",
		"POST",
		"/saves/mods",
		LoadModsFromSaveHandler,
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
		ModPortalListModsHandler,
	}, {
		"ModPortalGetModInfo",
		"GET",
		"/mods/portal/info/{mod}",
		ModPortalModInfoHandler,
	}, {
		"ModPortalInstallMod",
		"POST",
		"/mods/portal/install",
		ModPortalInstallHandler,
	}, {
		"ModPortalLogin",
		"POST",
		"/mods/portal/login",
		ModPortalLoginHandler,
	}, {
		"ModPortalLoginStatus",
		"GET",
		"/mods/portal/loginstatus",
		ModPortalLoginStatusHandler,
	}, {
		"ModPortalLogout",
		"GET",
		"/mods/portal/logout",
		ModPortalLogoutHandler,
	}, {
		"ModPortalInstallMultiple",
		"POST",
		"/mods/portal/install/multiple",
		ModPortalInstallMultipleHandler,
	},
	// Mods Stuff
	{
		"ListInstalledMods",
		"GET",
		"/mods/list",
		ListInstalledModsHandler,
	}, {
		"ToggleMod",
		"POST",
		"/mods/toggle",
		ModToggleHandler,
	}, {
		"DeleteMod",
		"POST",
		"/mods/delete",
		ModDeleteHandler,
	}, {
		"DeleteAllMods",
		"POST",
		"/mods/delete/all",
		ModDeleteAllHandler,
	}, {
		"UpdateMod",
		"POST",
		"/mods/update",
		ModUpdateHandler,
	}, {
		"UploadMod",
		"POST",
		"/mods/upload",
		ModUploadHandler,
	}, {
		"DownloadMods",
		"GET",
		"/mods/download",
		ModDownloadHandler,
	},
	// Mod Packs
	{
		"ModPacksList",
		"GET",
		"/mods/packs/list",
		ModPackListHandler,
	}, {
		"ModPackCreate",
		"POST",
		"/mods/packs/create",
		ModPackCreateHandler,
	}, {
		"ModPackDelete",
		"POST",
		"/mods/packs/{modpack}/delete",
		ModPackDeleteHandler,
	}, {
		"ModPackDownload",
		"GET",
		"/mods/packs/{modpack}/download",
		ModPackDownloadHandler,
	}, {
		"LoadModPack",
		"POST",
		"/mods/packs/{modpack}/load",
		ModPackLoadHandler,
	},
	// Mods inside Mod Packs
	{
		"ModPackListMods",
		"GET",
		"/mods/packs/{modpack}/list",
		ModPackModListHandler,
	}, {
		"ModPackToggleMod",
		"POST",
		"/mods/packs/{modpack}/mod/toggle",
		ModPackModToggleHandler,
	}, {
		"ModPackDeleteMod",
		"POST",
		"/mods/packs/{modpack}/mod/delete",
		ModPackModDeleteHandler,
	}, {
		"ModPackDeleteAllMod",
		"POST",
		"/mods/packs/{modpack}/mod/delete/all",
		ModPackModDeleteAllHandler,
	}, {
		"ModPackUpdateMod",
		"POST",
		"/mods/packs/{modpack}/mod/update",
		ModPackModUpdateHandler,
	}, {
		"ModPackUploadMod",
		"POST",
		"/mods/packs/{modpack}/mod/upload",
		ModPackModUploadHandler,
	}, {
		"ModPackModPortalInstallMod",
		"POST",
		"/mods/packs/{modpack}/portal/install",
		ModPackModPortalInstallHandler,
	}, {
		"ModPackModPortalInstallMultiple",
		"POST",
		"/mods/packs/{modpack}/portal/install/multiple",
		ModPackModPortalInstallMultipleHandler,
	},
}
