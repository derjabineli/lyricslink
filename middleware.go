package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/derjabineli/lyricslink/internal/auth"
	"github.com/google/uuid"
)

type contextKey string

const userIDKey contextKey = "userID"

func (cfg *config) getUserIDFromCookie(r *http.Request) (uuid.UUID, error) {
	cookie, err := r.Cookie("ll_user")
	if err != nil {
		return uuid.Nil, err
	}
	return auth.ValidateJWT(cookie.Value, cfg.tokenSecret) // Validate JWT
}

func (cfg *config) authMiddleware(next http.HandlerFunc) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := cfg.getUserIDFromCookie(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserIDFromContext(r *http.Request) (uuid.UUID, error) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("no user id present")
	}

	return userID, nil
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