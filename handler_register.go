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

	cookie, err := newJWT(user.ID, cfg.tokenSecret, cfg.tokenDuration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server error occured. Please try again")
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

func newJWT(id uuid.UUID, tokenSecret string, expiresIn time.Duration) (*http.Cookie, error) {
	jwtToken, err := auth.MakeJWT(id, tokenSecret, expiresIn)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     "ll_user",
		Value:    jwtToken,
		HttpOnly: true, // Make the cookie inaccessible to JavaScript
		Secure:   false, // Ensure the cookie is only sent over HTTPS
		SameSite: http.SameSiteLaxMode, // Protect against CSRF attacks
		Expires:  time.Now().Add(24 * time.Hour), // Set cookie expiration
		Path:     "/", // Define cookie scope
	}

	return cookie, nil
}