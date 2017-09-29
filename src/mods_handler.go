package main

import (
    "net/http"
    "fmt"
    "encoding/json"
    "log"
    "github.com/gorilla/mux"
    "archive/zip"
    "path/filepath"
    "os"
    "io"
    "errors"
)

// Returns JSON response of all mods installed in factorio/mods
func listInstalledModsHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    resp := JSONResponse{
        Success: false,
    }

    w.Header().Set("Content-Type", "application/json;charset=UTF-8")

    mods, err := newMods(config.FactorioModsDir)

    if err != nil {
        w.WriteHeader(500)
        resp.Data = fmt.Sprintf("Error in ListInstalledMods handler: %s", err)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error in list mods: %s", err)
        }
        return
    }

    resp.Data = mods.listInstalledMods().ModsResult
    resp.Success = true

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("Error in list mods: %s", err)
    }
}

// Returns JSON response with success or error-message
func LoginFactorioModPortal(w http.ResponseWriter, r *http.Request) {
    var err error
    resp := JSONResponse{
        Success: false,
    }

    w.Header().Set("Content-Type", "application/json;charset=UTF-8")

    username := r.FormValue("username")
    password := r.FormValue("password")

    login_status, err, statusCode := factorioLogin(username, password)
    if login_status == "" && err == nil {
        resp.Data = true
    }

    w.WriteHeader(statusCode)

    if err != nil {
        resp.Data = fmt.Sprintf("Error in getUserToken or LoginFactorioModPortal handler: %s", err)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error in Factorio-Login: %s", err)
        }
        return
    }

    resp.Success = true

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("Error in Factorio-Login: %s", err)
    }
}

func LoginstatusFactorioModPortal(w http.ResponseWriter, r *http.Request) {
    var err error
    resp := JSONResponse{
        Success: false,
    }

    var credentials FactorioCredentials
    resp.Data, err = credentials.load()

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        resp.Data = fmt.Sprintf("Error in getUserToken or LoginFactorioModPortal handler: %s", err)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error in Factorio-Login: %s", err)
        }
        return
    }

    resp.Success = true

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("Error in Factorio-Login: %s", err)
    }
}

//Returns JSON response with the found mods
func ModPortalSearchHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    resp := JSONResponse{
        Success: false,
    }

    w.Header().Set("Content-Type", "application/json;charset=UTF-8")

    //Get Data out of the request
    search_keyword := r.FormValue("search")

    var statusCode int
    resp.Data, err, statusCode = searchModPortal(search_keyword)

    w.WriteHeader(statusCode)

    if err != nil {
        w.WriteHeader(500)
        resp.Data = fmt.Sprintf("Error in searchModPortal: %s", err)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error in searchModPortal: %s", err)
        }
        return
    }

    resp.Success = true

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("Error in ModPortalSearch: %s", err)
    }
}

//Returns JSON response with the mod details
func ModPortalDetailsHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    resp := JSONResponse{
        Success: false,
    }

    w.Header().Set("Content-Type", "application/json;charset=UTF-8")

    //Get Data out of the request
    mod_id := r.FormValue("mod_id")

    var statusCode int
    resp.Data, err, statusCode = getModDetails(mod_id)

    w.WriteHeader(statusCode)

    if err != nil {
        resp.Data = fmt.Sprintf("Error in searchModPortal: %s", err)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error in searchModPortal: %s", err)
        }
        return
    }

    resp.Success = true

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("Error in ModPortalSearch: %s", err)
    }
}

func ModPortalInstallHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    resp := JSONResponse{
        Success: false,
    }

    w.Header().Set("Content-Type", "application/json;charset=UTF-8")

    //Get Data out of the request
    downloadUrl := r.FormValue("link")
    filename := r.FormValue("filename")
    mod_name := r.FormValue("modName")

    mods, err := newMods(config.FactorioModsDir)
    if err == nil {
        err = mods.downloadMod(downloadUrl, filename, mod_name)
    }

    if err != nil {
        w.WriteHeader(500)
        resp.Data = fmt.Sprintf("Error in installMod: %s", err)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error in installMod: %s", err)
        }
        return
    }

    resp.Data = mods.listInstalledMods()
    resp.Success = true

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("Error in ModPortalInstallHandler: %s", err)
    }
}

func ToggleModHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    resp := JSONResponse{
        Success: false,
    }

    w.Header().Set("Content-Type", "application/json;charset=UTF-8")

    //Get Data out of the request
    mod_name := r.FormValue("mod_name")

    mods, err := newMods(config.FactorioModsDir)
    if err == nil {
        err = mods.ModSimpleList.toggleMod(mod_name)
    }

    if err != nil {
        w.WriteHeader(500)
        resp.Data = fmt.Sprintf("Error in listInstalledModsByFolder: %s", err)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error in listInstalledModsByFolder: %s", err)
        }
        return
    }

    resp.Data = mods.listInstalledMods().ModsResult
    resp.Success = true

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("Error in ToggleModHandler: %s", err)
    }
}

func DeleteModHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    resp := JSONResponse{
        Success: false,
    }

    w.Header().Set("Content-Type", "application/json;charset=UTF-8")

    //Get Data out of the request
    mod_name := r.FormValue("mod_name")

    mods, err := newMods(config.FactorioModsDir)
    if err == nil {
        mods.deleteMod(mod_name)
    }

    if err != nil {
        w.WriteHeader(500)
        resp.Data = fmt.Sprintf("Error in deleteMod: %s", err)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error in DeleteModHandler: %s", err)
        }
        return
    }

    resp.Data = mods.listInstalledMods().ModsResult
    resp.Success = true

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("Error in DeleteModHandler: %s", err)
    }
}

func UpdateModHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    resp := JSONResponse{
        Success: false,
    }

    w.Header().Set("Content-Type", "application/json;charset=UTF-8")

    //Get Data out of the request
    mod_name := r.FormValue("mod_name")
    download_url := r.FormValue("downloadUrl")
    file_name := r.FormValue("filename")

    mods, err := newMods(config.FactorioModsDir)
    if err == nil {
        err = mods.updateMod(mod_name, download_url, file_name)
    }

    if err != nil {
        w.WriteHeader(500)
        resp.Data = fmt.Sprintf("Error in deleteMod: %s", err)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error in DeleteModHandler: %s", err)
        }
        return
    }

    resp.Data = mods.listInstalledMods()
    resp.Success = true

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("Error in DeleteModHandler: %s", err)
    }
}

func UploadModHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    resp := JSONResponseFileInput{
        Success: false,
    }

    w.Header().Set("Content-Type", "application/json;charset=UTF-8")

    r.ParseMultipartForm(32 << 20)

    mods, err := newMods(config.FactorioModsDir)
    if err == nil {
        for file_key, mod_file := range r.MultipartForm.File["mod_file"] {
            err = mods.uploadMod(mod_file)
            if err != nil {
                resp.ErrorKeys = append(resp.ErrorKeys, file_key)
                resp.Error = "An error occurred during upload or saving, pls check manually, if all went well and delete invalid files. (This program also could be crashed)"
            }
        }
    }

    if err != nil {
        w.WriteHeader(500)
        resp.Data = fmt.Sprintf("Error in uploadMod, listing mods wasn't successful: %s", err)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error in uploadMod, listing mods wasn't successful: %s", err)
        }
        return
    }

    resp.Data = mods.listInstalledMods()
    resp.Success = true

    if err = json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("Error in UploadModHandler: %s", err)
    }
}

func ListModPacksHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    resp := JSONResponse{
        Success: false,
    }

    w.Header().Set("Content-Type", "application/json;charset=UTF-8")

    mod_pack_map, err := newModPackMap()

    if err != nil {
        w.WriteHeader(500)
        w.WriteHeader(http.StatusInternalServerError)
        resp.Data = fmt.Sprintf("Error listing modpack files: %s", err)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error listing modpacks: %s", err)
        }
        return
    }

    resp.Data = mod_pack_map.listInstalledModPacks()
    resp.Success = true

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("Error listing saves: %s", err)
    }
}

func CreateModPackHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    resp := JSONResponse{
        Success: false,
    }

    name := r.FormValue("name")

    w.Header().Set("Content-Type", "application/json;charset=UTF-8")

    mod_pack_map, err := newModPackMap()
    if err == nil {
        err = mod_pack_map.createModPack(name)
    }

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        resp.Data = fmt.Sprintf("Error creating modpack file: %s", err)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error creating modpack: %s", err)
        }
        return
    }

    resp.Data = mod_pack_map.listInstalledModPacks()
    resp.Success = true

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("Error creating modpack response: %s", err)
    }
}

func DownloadModPackHandler(w http.ResponseWriter, r *http.Request) {
    var err error

    vars := mux.Vars(r)
    modpack := vars["modpack"]

    mod_pack_map, err := newModPackMap()
    if err != nil {
        log.Printf("error on loading modPacks: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    if mod_pack_map.checkModPackExists(modpack) {
        zip_writer := zip.NewWriter(w)
        defer zip_writer.Close()

        //iterate over folder and create everything in the zip
        err = filepath.Walk(filepath.Join(config.FactorioModPackDir, modpack), func(path string, info os.FileInfo, err error) error {
            if info.IsDir() == false {
                writer, err := zip_writer.Create(info.Name())
                if err != nil {
                    log.Printf("error on creating new file inside zip: %s", err)
                    return err
                }

                file, err := os.Open(path)
                if err != nil {
                    log.Printf("error on opening modfile: %s", err)
                    return err
                }
                defer file.Close()

                _, err = io.Copy(writer, file)
                if err != nil {
                    log.Printf("error on copying file into zip: %s", err)
                    return err
                }
            }

            return nil
        })
        if err != nil {
            log.Printf("error on walking over the modpack: %s", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
    } else {
        log.Printf("requested modPack doesnt exist")
        w.WriteHeader(http.StatusNotFound)
        return
    }

    writer_header := w.Header()
    writer_header.Set("Content-Type", "application/zip;charset=UTF-8")
    writer_header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", modpack + ".zip"))
}

func DeleteModPackHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    resp := JSONResponse{
        Success: false,
    }

    w.Header().Set("Content-Type", "application/json;charset=UTF-8")

    name := r.FormValue("name")

    mod_pack_map, err := newModPackMap()
    if err == nil {
        err = mod_pack_map.deleteModPack(name)
    }

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        resp.Data = fmt.Sprintf("Error deleting modpack file: %s", err)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error deleting modpack: %s", err)
        }
        return
    }

    resp.Data = mod_pack_map.listInstalledModPacks()
    resp.Success = true

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("Error creating delete modpack response: %s", err)
    }
}

func LoadModPackHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    resp := JSONResponse{
        Success: false,
    }

    w.Header().Set("Content-Type", "application/json;charset=UTF-8")

    name := r.FormValue("name")

    mod_pack_map, err := newModPackMap()
    if err == nil {
        mod_pack_map[name].loadModPack()
    }

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        resp.Data = fmt.Sprintf("Error loading modpack file: %s", err)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error loading modpack: %s", err)
        }
        return
    }

    resp.Data = mod_pack_map[name].Mods.listInstalledMods()
    resp.Success = true

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("Error creating loading modpack response: %s", err)
    }
}

func ModPackToggleModHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    resp := JSONResponse{
        Success: false,
    }

    w.Header().Set("Content-Type", "application/json;charset=UTF-8")

    mod_name := r.FormValue("mod_name")
    mod_pack_name := r.FormValue("mod_pack")

    mod_pack_map, err := newModPackMap()
    if err == nil {
        err = mod_pack_map[mod_pack_name].Mods.ModSimpleList.toggleMod(mod_name)
    }
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        resp.Data = fmt.Sprintf("Error loading modpack file: %s", err)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error loading modpack: %s", err)
        }
        return
    }

    resp.Data = mod_pack_map.listInstalledModPacks()
    resp.Success = true

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("Error creating loading modpack response: %s", err)
    }
}

func ModPackDeleteModHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    resp := JSONResponse{
        Success: false,
    }

    w.Header().Set("Content-Type", "application/json;charset=UTF-8")

    mod_name := r.FormValue("mod_name")
    mod_pack_name := r.FormValue("mod_pack_name")

    mod_pack_map, err := newModPackMap()
    if err == nil {
        if mod_pack_map.checkModPackExists(mod_pack_name) {
            err = mod_pack_map[mod_pack_name].Mods.deleteMod(mod_name)
        } else {
            err = errors.New("ModPack " + mod_pack_name + " does not exist")
        }
    }
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        resp.Data = fmt.Sprintf("Error loading modpack file: %s", err)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error loading modpack: %s", err)
        }
        return
    }

    resp.Data = mod_pack_map.listInstalledModPacks()
    resp.Success = true

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("Error creating loading modpack response: %s", err)
    }

    return
}

func ModPackUpdateModHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    resp := JSONResponse{
        Success: false,
    }

    w.Header().Set("Content-Type", "application/json;charset=UTF-8")

    //Get Data out of the request
    mod_name := r.FormValue("mod_name")
    download_url := r.FormValue("downloadUrl")
    file_name := r.FormValue("filename")
    mod_pack_name := r.FormValue("mod_pack_name")

    mod_pack_map, err := newModPackMap()
    if err == nil {
        if mod_pack_map.checkModPackExists(mod_pack_name){
            err = mod_pack_map[mod_pack_name].Mods.updateMod(mod_name, download_url, file_name)
        } else {
            err = errors.New("ModPack " + mod_pack_name + "does not exist")
        }
    }

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        resp.Data = fmt.Sprintf("Error in deleteMod: %s", err)
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error in DeleteModHandler: %s", err)
        }
        return
    }

    resp.Data = mod_pack_map.listInstalledModPacks()
    resp.Success = true

    if err := json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("Error in DeleteModHandler: %s", err)
    }
}