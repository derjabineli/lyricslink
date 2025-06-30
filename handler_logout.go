package main

import (
	"context"
	"net/http"
	"time"

	"github.com/derjabineli/lyricslink/internal/auth"
)

func (cfg *config) handlerLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	refreshCookie, err := r.Cookie(auth.RefreshCookieName)
	if err != nil {
		http.Error(w, "Missing refresh cookie", http.StatusBadRequest)
		return
	}

	sessionID, err := auth.ExtractSessionID(refreshCookie.Value, cfg.tokenSecret)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	if err := cfg.db.RevokeSession(context.Background(), sessionID); err != nil {
		http.Error(w, "Could not revoke token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.RefreshCookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Path:     "/",
	})

	w.WriteHeader(http.StatusNoContent)
}

