package main

import (
    "net/http"
    "net/url"
    "log"
    "io/ioutil"
)

type LoginErrorResponse struct {
    Message string  `json:"message"`
    Status  int     `json:"status"`
}
type LoginSuccessResponse struct {
    UserKey []string  `json:""`
}
//Log the user into factorio, so mods can be downloaded
func getUserToken(username string, password string) (string, error, int) {
    resp, get_err := http.PostForm("https://auth.factorio.com/api-login",
        url.Values{"require_game_ownership": {"true"}, "username": {username}, "password": {password}})
    if get_err != nil {
        log.Fatal(get_err)
        return "error", get_err, 500
    }

    //get the response-text
    text, err_io := ioutil.ReadAll(resp.Body)
    resp.Body.Close()

    text_string := string(text)

    if err_io != nil {
        log.Fatal(err_io)
        return "error", err_io, resp.StatusCode
    }

    return text_string, nil, resp.StatusCode
}


//Search inside the factorio mod portal
func searchModPortal(keyword string) (string, error, int) {
    //resp, get_err := http.Get
    req, err := http.NewRequest(http.MethodGet, "https://mods.factorio.com/api/mods", nil)
    if err != nil {
        return "error", err, 500
    }

    query := req.URL.Query()
    query.Add("q", keyword)
    req.URL.RawQuery = query.Encode()

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "error", err, 500
    }

    text, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    if err != nil {
        return "error", err, 500
    }

    text_string := string(text)

    return text_string, nil, resp.StatusCode
}

func getModDetails(modId string) (string, error, int) {
    var err error
    new_link := "https://mods.factorio.com/api/mods/" + modId
    resp, err := http.Get(new_link)

    if err != nil {
        return "error", err, 500
    }

    //get the response-text
    text, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()

    text_string := string(text)

    if err != nil {
        log.Fatal(err)
        return "error", err, resp.StatusCode
    }

    return text_string, nil, resp.StatusCode
}
