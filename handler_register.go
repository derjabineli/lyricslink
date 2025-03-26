package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/derjabineli/lyricslink/internal/auth"
	"github.com/derjabineli/lyricslink/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type registerParameters struct {
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type successResponse struct {
	Success bool `json:"success"`
}

func (cfg *config) register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	registerParams := registerParameters{}
	decoder.Decode(&registerParams)

	userId := uuid.New()
	hashedPassword, err := auth.HashPassword(registerParams.Password)
	if err != nil {
		respondWithError(w, 400, "We apologize. There seems to be an issue with your password. Please try again or use a different password.")
		return
	}

	user, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: userId,
		FirstName: registerParams.FirstName,
		LastName: registerParams.LastName,
		Email: registerParams.Email,
		HashedPassword: hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		PcAuthorized: false,
	})

	if err != nil {
		errorMsg := handleDBError(err)
		respondWithError(w, http.StatusInternalServerError, errorMsg)
		return
	}

	cookie, err := auth.NewJWTCookie(user.ID, cfg.tokenSecret, cfg.tokenDuration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server error occured. Please try again")
		return
	}
	
	http.SetCookie(w, cookie)

	success := successResponse{
		Success: true,
	} 

	respondWithJSON(w, 200, success)
}

func handleDBError(err error) string {
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505":
			return "A user with this email already exists"
		default:
			return "⚠️ Other DB error occurred"
		}
	} else {
		return "Unknown error occured. Please try again"
	}
}