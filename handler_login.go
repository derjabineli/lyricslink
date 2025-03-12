package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/derjabineli/lyricslink/internal/auth"
)

type accessParameters struct {
	GrantType string `json:"grant_type"`
	Code string `json:"code"`
	ClientID string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI string `json:"redirect_uri"`
}

type authorizationParameters struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope string `json:"scope"`
	CreatedAt int `json:"created_at"`
}

type loginParameters struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (cfg *config) loginStatic(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./views/login.html")
		if err != nil {
			http.Error(w, "Error loading login page", http.StatusInternalServerError)
			log.Println("Template parsing error:", err)
			return
		}
	
		err = t.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error rendering login page", http.StatusInternalServerError)
			log.Println("Template execution error:", err)
		}
}

func (cfg *config) login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	loginCreds := loginParameters{}
	decoder.Decode(&loginCreds)

	user, err := cfg.db.GetUserByEmail(context.Background(), loginCreds.Email)
	if err != nil {
		respondWithError(w, 401, "User does not exist")
	}

	err = auth.CheckPasswordHash(loginCreds.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Wrong Password")
	}

	cookie, err := newJWT(user.ID, cfg.tokenSecret, cfg.tokenDuration)
	if err != nil {
		respondWithError(w, 401, err.Error())
	}

	http.SetCookie(w, cookie)

	success := successResponse{
		Success: true,
	} 

	respondWithJSON(w, 200, success)
}

func (cfg *config) planningcentercallback(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.RequestURI)
	if err != nil {
		return
	}
	code := parsedURL.Query().Get("code")
	fmt.Printf("Code: %v", code)
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