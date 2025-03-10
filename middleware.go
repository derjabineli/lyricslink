package main

import (
	"log"
	"net/http"

	"github.com/derjabineli/lyricslink/internal/auth"
	"github.com/google/uuid"
)

func (cfg *config) getUserIDFromCookie(r *http.Request) (uuid.UUID, error) {
	cookie, err := r.Cookie("ll_user")
	if err != nil {
		return uuid.Nil, err
	}
	return auth.ValidateJWT(cookie.Value, cfg.tokenSecret) // Validate JWT
}

func (cfg *config) authMiddleware(next http.HandlerFunc) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := cfg.getUserIDFromCookie(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		log.Printf("Authenticated user ID: %v", id)
		next.ServeHTTP(w, r)
	})
}

func (cfg *config) guestOnlyMiddleware(next http.HandlerFunc) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := cfg.getUserIDFromCookie(r)
		if err == nil {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}