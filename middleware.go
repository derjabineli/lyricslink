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

func (cfg *config) getUserIDFromCookie(r *http.Request, tokenName string, tokenType auth.TokenType) (uuid.UUID, error) {
	cookie, err := r.Cookie(tokenName)
	if err != nil {
		return uuid.Nil, errors.New("no token found")
	}
	return auth.ValidateJWT(cookie.Value, cfg.tokenSecret, tokenType) // Validate JWT
}

func (cfg *config) authMiddleware(next http.HandlerFunc) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := cfg.getUserIDFromCookie(r, "ll_user", auth.AccessTokenType)
		if err == nil {
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		userID, err = cfg.getUserIDFromCookie(r, "ll_refresh", auth.RefreshTokenType)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		newAccessCookie, err := auth.NewAccessTokenCookie(userID, cfg.tokenSecret)
		if err != nil {
			return
		}
		http.SetCookie(w, newAccessCookie)

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
		_, err := cfg.getUserIDFromCookie(r, "ll_user", auth.AccessTokenType)
		if err == nil {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		_, err = cfg.getUserIDFromCookie(r, "ll_refresh", auth.RefreshTokenType)
		if err == nil {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}