package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/derjabineli/lyricslink/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type config struct {
	db *database.Queries
	pcClient string
	pcSecret string
	pcRedirect string
	tokenSecret string
	pcSongRoute string
	tokenDuration time.Duration
}


func main() {
	godotenv.Load()
	PORT := os.Getenv("PORT")
	pcClient := os.Getenv("PC_CLIENTID")
	pcSecret := os.Getenv("PC_SECRET")
	pcRedirect := os.Getenv("PC_REDIRECT")
	tokenSecret := os.Getenv("TOKEN_SECRET")
	pcSongRoute := os.Getenv("PC_SONG_ROUTE")
	dbURL := os.Getenv("DB_URL")
  	db, err := sql.Open("postgres", dbURL)
  	if err != nil {
   		fmt.Print(err)
    	return 
  	}
	dbQueries := database.New(db)

	cfg := config{db: dbQueries, pcClient: pcClient, pcSecret: pcSecret, pcRedirect: pcRedirect, tokenSecret: tokenSecret, tokenDuration: time.Hour * 8, pcSongRoute: pcSongRoute}
	

	mux := http.NewServeMux()

	// PAGES
	mux.Handle("/", cfg.guestOnlyMiddleware(cfg.handlerHome))
	mux.Handle("/login", cfg.guestOnlyMiddleware(cfg.loginStatic))
	mux.HandleFunc("/static/", staticHandler)
	mux.Handle("/dashboard", cfg.authMiddleware(cfg.handlerDashboard)) 
	mux.Handle("/events/{id}", cfg.authMiddleware(cfg.handlerEvents)) 
	mux.Handle("/settings", cfg.authMiddleware(cfg.handlerSettings)) 
	mux.HandleFunc("GET /live/{id}", cfg.handlerLive)
	
	// API
	mux.Handle("POST /api/events", cfg.authMiddleware(cfg.addEvent))
	mux.Handle("PUT /api/events", cfg.authMiddleware(cfg.updateEventDate))
	mux.Handle("DELETE /api/events/{id}", cfg.authMiddleware(cfg.deleteEvent))
	mux.Handle("POST /api/songs", cfg.authMiddleware(cfg.getSongs))
	mux.Handle("GET /api/songs/{id}", cfg.authMiddleware(cfg.getArrangements))
	mux.Handle("GET /api/songs/sync", cfg.authMiddleware(cfg.syncPlanningCenterSongs))
	mux.Handle("POST /api/events_arrangements", cfg.authMiddleware(cfg.addArrangementToEvent))
	mux.Handle("PUT /api/events_arrangements/{id}", cfg.authMiddleware(cfg.updateEventArrangement))
	mux.Handle("DELETE /api/events_arrangements/{id}", cfg.authMiddleware(cfg.deleteEventArrangement))
	mux.HandleFunc("GET /api/logout", cfg.handlerLogout)

	// AUTH
	mux.HandleFunc("/pc/callback", cfg.loginWithPC)

	server := &http.Server{
		Handler: mux,
		Addr: ":" + PORT,
	  }
	  
	  fmt.Printf("Running on Port %v\n", PORT)
	  log.Fatal(server.ListenAndServe())
}