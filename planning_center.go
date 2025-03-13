package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

func (cfg *config) loginPC(w http.ResponseWriter, r *http.Request) {
	link := "https://api.planningcenteronline.com/oauth/authorize?client_id=" + cfg.pcClient + "&redirect_uri=" + cfg.pcRedirect + "&response_type=code&scope=services"

	t, err := template.ParseFiles("./views/pc.html")
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}
	
	err = t.Execute(w, link)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

func (cfg *config) planningcentercallback(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.RequestURI)
	if err != nil {
		return
	}
	code := parsedURL.Query().Get("code")
	fmt.Printf("Code: %v\n", code)
	accessParams := accessParameters {
		GrantType: "authorization_code",
		Code: code,
		ClientID: cfg.pcClient,
		ClientSecret: cfg.pcSecret,
		RedirectURI: cfg.pcRedirect,
	}

	requestBody, err := json.Marshal(accessParams)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}	
	authReq, err := http.NewRequest("POST", "https://api.planningcenteronline.com/oauth/token", bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}
	authReq.Header.Add("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(authReq)
	decoder := json.NewDecoder(resp.Body)
	authParams := authorizationParameters{}
	decoder.Decode(&authParams)

	fmt.Printf("Access token: %v\n", authParams.AccessToken)
	
	cookie := http.Cookie{
		Name: "jwt-token",
		Value: "test", 
		Path: "/",
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/dashboard", http.StatusPermanentRedirect)
}