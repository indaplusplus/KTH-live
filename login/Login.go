package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

var LOGIN_API_URL string = "https://login2.datasektionen.se"
var LOGIN_COMPLETE_URL string = "http://localhost:8021/loginComplete"
var LOGIN_API_KEY string = os.Getenv("LOGIN_API_KEY")

type jsonData struct {
	Token string `json:"token,omitempty"`
	User  string `json:"user,omitempty"`
}

type verifyData struct {
	First   string `json:"first_name,omitempty"`
	Last    string `json:"last_name,omitempty"`
	Email   string `json:"emails,omitempty"`
	User    string `json:"user,omitempty"`
	Ugkthid string `json:"ugkthid,omitempty"`
}

func login(w http.ResponseWriter, r *http.Request) {
	url := LOGIN_API_URL + "/login?callback=" + LOGIN_COMPLETE_URL + "?token="
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func loginComplete(w http.ResponseWriter, r *http.Request) {
	var token string = getTokenFromURL(r)
	url := LOGIN_API_URL + "/verify/" + token + ".json?api_key=" + LOGIN_API_KEY
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	var tmp verifyData
	json.NewDecoder(resp.Body).Decode(&tmp)
	if tmp.User == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	loginToken(token, tmp.User)
	json.NewEncoder(w).Encode(jsonData{Token: token})
}

func logout(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromJson(r)
	logoutToken(token)
	w.WriteHeader(http.StatusOK)
}

func isLoggedin(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromJson(r)
	if containsToken(token) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromJson(r)
	user, has := getLoggedInUser(token)
	if has {
		json.NewEncoder(w).Encode(jsonData{User: user})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	//TODO if needed
}

func getTokenFromJson(r *http.Request) string {
	var data jsonData
	json.NewDecoder(r.Body).Decode(&data)
	return data.Token
}

func getTokenFromURL(r *http.Request) string {
	return string(regexp.MustCompile(`token=\w{22}`).Find([]byte(fmt.Sprint(r.URL))))[6:]
}
