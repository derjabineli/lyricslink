package main

import (
	"fmt"
	"net/http"

	"github.com/derjabineli/lyricslink/internal/auth"
)

func (cfg *config) authMiddleware(next http.HandlerFunc) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("ll_user")
		if err != nil {
			http.Redirect(w, r, "/", 301)
			return
		}
		fmt.Printf("Cookie: %v\n", cookie)
		id, err := auth.ValidateJWT(cookie.Value, cfg.tokenSecret)
		if err != nil {}
		fmt.Printf("Id: %v\n", id)
		next.ServeHTTP(w, r)
	})
}