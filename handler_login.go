package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

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