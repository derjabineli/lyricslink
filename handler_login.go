package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/derjabineli/lyricslink/internal/auth"
)

type loginParameters struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (cfg *config) loginStatic(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./frontend/views/login.html")
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
		respondWithError(w, http.StatusUnauthorized, "No user with that email or password found")
		return
	}

	err = auth.CheckPasswordHash(loginCreds.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No user with that email or password found")
		return
	}

	cookie, err := auth.NewJWTCookie(user.ID, cfg.tokenSecret, cfg.tokenDuration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was an issue logging you in. Please try again")
		return
	}

	http.SetCookie(w, cookie)

	success := successResponse{
		Success: true,
	} 

	respondWithJSON(w, 200, success)
}
