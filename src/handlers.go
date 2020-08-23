package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const readHttpBodyError = "Could not read the Request Body."

type JSONResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,string"`
}
type JSONResponseFileInput struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,string"`
	Error     string      `json:"error"`
	ErrorKeys []int       `json:"errorkeys"`
}

func WriteResponse(w http.ResponseWriter, data interface{}) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error writing response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ReadRequestBody(w http.ResponseWriter, r *http.Request, resp *interface{}) (body []byte, err error) {
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		*resp = fmt.Sprintf("%s: %s", readHttpBodyError, err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

// Lists all save files in the factorio/saves directory
func ListSaves(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	savesList, err := listSaves(config.FactorioSavesDir)
	if err != nil {
		resp.Success = false
		resp.Data = fmt.Sprintf("Error listing save files: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing saves: %s", err)
		}
		return
	}

	loadLatest := Save{Name: "Load Latest"}
	savesList = append(savesList, loadLatest)

	resp.Data = savesList

	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error listing saves: %s", err)
	}
}

func DLSave(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")

	vars := mux.Vars(r)
	save := vars["save"]
	saveName := filepath.Join(config.FactorioSavesDir, save)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", save))
	log.Printf("%s downloading: %s", r.Host, saveName)

	http.ServeFile(w, r, saveName)
}

func UploadSave(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	switch r.Method {
	case "GET":
		resp.Data = "Unsupported method"
		resp.Success = false
		if err = json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing mods: %s", err)
		}
	case "POST":
		log.Println("Uploading save file")

		r.ParseMultipartForm(32 << 20)

		for _, saveFile := range r.MultipartForm.File["savefile"] {
			file, err := saveFile.Open()
			if err != nil {
				resp.Success = false
				resp.Data = err.Error()
				json.NewEncoder(w).Encode(resp)
				log.Printf("Error in upload save formfile: %s", err.Error())
				return
			}
			defer file.Close()

			out, err := os.Create(filepath.Join(config.FactorioSavesDir, saveFile.Filename))
			if err != nil {
				resp.Success = false
				resp.Data = err.Error()
				json.NewEncoder(w).Encode(resp)
				log.Printf("Error in out: %s", err)
				return
			}
			defer out.Close()

			_, err = io.Copy(out, file)
			if err != nil {
				resp.Success = false
				resp.Data = err.Error()
				json.NewEncoder(w).Encode(resp)
				log.Printf("Error in io copy: %s", err)
				return
			}

			log.Printf("Uploaded save file: %s", saveFile.Filename)
			resp.Data = "File '" + saveFile.Filename + "' uploaded successfully"
			resp.Success = true
			json.NewEncoder(w).Encode(resp)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Deletes provided save
func RemoveSave(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	vars := mux.Vars(r)
	name := vars["save"]

	save, err := findSave(name)
	if err != nil {
		resp.Data = fmt.Sprintf("Error removing save: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error removing save %s", err)
		}
	}

	err = save.remove()
	if err == nil {
		// save was removed
		resp.Data = fmt.Sprintf("Removed save: %s", save.Name)
		resp.Success = true
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error removing save %s", err)
		}
	} else {
		log.Printf("Error in remove save handler: %s", err)
		resp.Data = fmt.Sprintf("Error in remove save handler: %s", err)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error removing save: %s", err)
		}
	}
}

// Launches Factorio server binary with --create flag to create save
// Url must include save name for creation of savefile
func CreateSaveHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	vars := mux.Vars(r)
	saveName := vars["save"]

	if saveName == "" {
		log.Printf("Error creating save, no name provided: %s", err)
		resp.Data = "No save name provided."
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding save handler response: %s", err)
		}
		return
	}

	saveFile := filepath.Join(config.FactorioSavesDir, saveName)
	cmdOut, err := createSave(saveFile)
	if err != nil {
		log.Printf("Error creating save: %s", err)
		resp.Data = "Error creating savefile."
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding save handler response: %s", err)
		}
		return
	}

	resp.Success = true
	resp.Data = fmt.Sprintf("Save %s created successfully. Command output: \n%s", saveName, cmdOut)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding save response: %s", err)
	}
}

// LogTail returns last lines of the factorio-current.log file
func LogTail(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	resp.Data, err = tailLog(config.FactorioLog)
	if err != nil {
		resp.Data = fmt.Sprintf("Could not tail %s: %s", config.FactorioLog, err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Could not tail %s: %s", config.FactorioLog, err)
		}
		return
	}

	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error tailing logfile: %s", err)
	}
}

// LoadConfig returns JSON response of config.ini file
func LoadConfig(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	configContents, err := loadConfig(config.FactorioConfigFile)
	if err != nil {
		log.Printf("Could not retrieve config.ini: %s", err)
		resp.Data = "Error getting config.ini"
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error tailing logfile: %s", err)
		}
	} else {
		resp.Data = configContents
		resp.Success = true
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding config file JSON reponse: %s", err)
	}

	log.Printf("Sent config.ini response")
}

func StartServer(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if FactorioServ.Running {
		resp.Data = "Factorio server is already running"
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding JSON response: %s", err)
		}
		return
	}

	switch r.Method {
	case "GET":
		log.Printf("GET not supported for startserver handler")
		resp.Data = "Unsupported method"
		resp.Success = false
		if err = json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing mods: %s", err)
		}
	case "POST":
		log.Printf("Starting Factorio server.")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error in starting factorio server handler body: %s", err)
			return
		}

		log.Printf("Starting Factorio server with settings: %v", string(body))

		err = json.Unmarshal(body, &FactorioServ)
		if err != nil {
			log.Printf("Error unmarshaling server settings JSON: %s", err)
			return
		}

		// Check if savefile was submitted with request to start server.
		if FactorioServ.Savefile == "" {
			log.Printf("Error starting Factorio server: no save file provided")
			resp.Success = false
			resp.Data = fmt.Sprintf("Error starting Factorio server: %s", "No save file provided")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error encoding config file JSON reponse: %s", err)
			}
			return
		}

		go func() {
			err = FactorioServ.Run()
			if err != nil {
				log.Printf("Error starting Factorio server: %+v", err)
				return
			}
		}()

		timeout := 0
		for timeout <= 3 {
			time.Sleep(1 * time.Second)
			if FactorioServ.Running {
				resp.Data = fmt.Sprintf("Factorio server with save: %s started on port: %d", FactorioServ.Savefile, FactorioServ.Port)
				resp.Success = true
				log.Printf("Factorio server started on port: %v", FactorioServ.Port)
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					log.Printf("Error encoding config file JSON reponse: %s", err)
				}
				break
			} else {
				log.Printf("Did not detect running Factorio server attempt: %+v", timeout)
			}

			timeout++
		}
		if FactorioServ.Running == false {
			log.Printf("Error starting Factorio server: %s", err)
			resp.Data = fmt.Sprintf("Error starting Factorio server: %s", err)
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error encoding start server JSON response: %s", err)
			}
		}
	}
}

func StopServer(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if FactorioServ.Running {
		err := FactorioServ.Stop()
		if err != nil {
			log.Printf("Error in stop server handler: %s", err)
			resp.Data = fmt.Sprintf("Error in stop server handler: %s", err)
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error encoding config file JSON reponse: %s", err)
			}
			return
		}

		log.Printf("Stopped Factorio server.")
		resp.Success = true
		resp.Data = fmt.Sprintf("Factorio server stopped")
	} else {
		resp.Data = "Factorio server is not running"
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding config file JSON reponse: %s", err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding config file JSON reponse: %s", err)
	}
}

func KillServer(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if FactorioServ.Running {
		err := FactorioServ.Kill()
		if err != nil {
			log.Printf("Error in kill server handler: %s", err)
			resp.Data = fmt.Sprintf("Error in kill server handler: %s", err)
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error encoding config file JSON reponse: %s", err)
			}
			return
		}

		log.Printf("Killed Factorio server.")
		resp.Success = true
		resp.Data = fmt.Sprintf("Factorio server killed")
	} else {
		resp.Data = "Factorio server is not running"
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding config file JSON reponse: %s", err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding config file JSON reponse: %s", err)
	}
}

func CheckServer(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if FactorioServ.Running {
		resp.Success = true
		status := map[string]string{}
		status["status"] = "running"
		status["port"] = strconv.Itoa(FactorioServ.Port)
		status["savefile"] = FactorioServ.Savefile
		status["address"] = FactorioServ.BindIP
		resp.Data = status
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding config file JSON reponse: %s", err)
		}
	} else {
		resp.Success = true
		status := map[string]string{}
		status["status"] = "stopped"
		resp.Data = status
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding config file JSON reponse: %s", err)
		}
	}
}

func FactorioVersion(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	status := map[string]string{}
	status["version"] = FactorioServ.Version.String()
	status["base_mod_version"] = FactorioServ.BaseModVersion

	resp.Data = status

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error loading Factorio Version: %s", err)
	}
}

// Unmarshall the User object from the given bytearray
// This function has side effects (it will write to resp and to w, in case of an error)
func UnmarshallUserJson(body []byte, resp *interface{}, w http.ResponseWriter) (user User, err error) {
	err = json.Unmarshal(body, &user)
	if err != nil {
		*resp = fmt.Sprintf("Unable to parse the request body: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusBadRequest)
	}
	return
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	// add resp to the response
	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	body, err := ReadRequestBody(w, r, &resp)
	if err != nil {
		return
	}

	user, err := UnmarshallUserJson(body, &resp, w)
	if err != nil {
		return
	}

	log.Printf("Logging in user: %s", user.Username)

	err = Auth.aaa.Login(w, r, user.Username, user.Password, "/")
	if err != nil {
		resp = fmt.Sprintf("Error loggin in user: %s, error: %s", user.Username, err)
		log.Println(resp)
		return
	}

	log.Printf("User: %s, logged in successfully", user.Username)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if err = Auth.aaa.Logout(w, r); err != nil {
		log.Printf("Error logging out current user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp = "User logged out successfully."
}

func GetCurrentLogin(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp interface{}

	// add resp to the response
	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	user, err := Auth.aaa.CurrentUser(w, r)
	if err != nil {
		resp = fmt.Sprintf("Error getting user status: %s, error: %s", user.Username, err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp = user
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	users, err := Auth.listUsers()
	if err != nil {
		resp = fmt.Sprintf("Error listing users: %s", err)
		log.Println(resp)
		return
	}

	resp = users
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	body, err := ReadRequestBody(w, r, &resp)
	if err != nil {
		return
	}

	log.Printf("Adding user: %v", string(body))

	user, err := UnmarshallUserJson(body, &resp, w)
	if err != nil {
		return
	}

	err = Auth.addUser(user.Username, user.Password, user.Email, user.Role)
	if err != nil {
		resp = fmt.Sprintf("Error in adding user {%s}: %s", user.Username, err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp = fmt.Sprintf("User: %s successfully added.", user.Username)
}

func RemoveUser(w http.ResponseWriter, r *http.Request) {
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	body, err := ReadRequestBody(w, r, &resp)
	if err != nil {
		return
	}

	user, err := UnmarshallUserJson(body, &resp, w)
	if err != nil {
		return
	}

	err = Auth.removeUser(user.Username)
	if err != nil {
		resp = fmt.Sprintf("Error in removing user {%s}, error: %s", user.Username, err)
		log.Println(resp)
		return
	}

	resp = fmt.Sprintf("User: %s successfully removed.", user.Username)
}

// GetServerSettings returns JSON response of server-settings.json file
func GetServerSettings(w http.ResponseWriter, r *http.Request) {
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	resp = FactorioServ.Settings

	log.Printf("Sent server settings response")
}

func UpdateServerSettings(w http.ResponseWriter, r *http.Request) {
	var resp interface{}

	defer func() {
		WriteResponse(w, resp)
	}()

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	body, err := ReadRequestBody(w, r, &resp)
	if err != nil {
		return
	}
	log.Printf("Received settings JSON: %s", body)

	err = json.Unmarshal(body, &FactorioServ.Settings)
	if err != nil {
		resp = fmt.Sprintf("Error unmarhaling server settings JSON: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	settings, err := json.MarshalIndent(&FactorioServ.Settings, "", "  ")
	if err != nil {
		resp = fmt.Sprintf("Failed to marshal server settings: %s", err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = ioutil.WriteFile(filepath.Join(config.FactorioConfigDir, config.SettingsFile), settings, 0644)
	if err != nil {
		resp = fmt.Sprintf("Failed to save server settings: %v\n", err)
		log.Println(resp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Saved Factorio server settings in server-settings.json")

	if (FactorioServ.Version.Greater(Version{0, 17, 0})) {
		// save admins to adminJson
		admins, err := json.MarshalIndent(FactorioServ.Settings["admins"], "", "  ")
		if err != nil {
			resp = fmt.Sprintf("Failed to marshal admins-Setting: %s", err)
			log.Println(resp)
			return
		}

		err = ioutil.WriteFile(filepath.Join(config.FactorioConfigDir, config.FactorioAdminFile), admins, 0664)
		if err != nil {
			resp = fmt.Sprintf("Failed to save admins: %s", err)
			log.Println(resp)
			return
		}
	}

	resp = fmt.Sprintf("Settings successfully saved")
}
