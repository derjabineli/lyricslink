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
	mux.Handle("/", cfg.guestOnlyMiddleware(handlerHome))
	mux.Handle("/login", cfg.guestOnlyMiddleware(cfg.loginStatic))
	mux.HandleFunc("/static/", staticHandler)
	mux.HandleFunc("/loginPC", cfg.loginPC)
	mux.Handle("/dashboard", cfg.authMiddleware(cfg.handlerDashboard)) 
	mux.Handle("/events/{id}", cfg.authMiddleware(cfg.handlerEvents)) 
	mux.Handle("/settings", cfg.authMiddleware(cfg.handlerSettings)) 
	
	// API
	mux.Handle("POST /api/events", cfg.authMiddleware(cfg.addEvent))
	mux.Handle("PUT /api/events", cfg.authMiddleware(cfg.updateEventDate))
	mux.Handle("POST /api/songs", cfg.authMiddleware(cfg.getSongs))
	mux.Handle("GET /api/songs/{id}", cfg.authMiddleware(cfg.getArrangements))
	mux.Handle("POST /api/event_arrangements", cfg.authMiddleware(cfg.addArrangementToEvent))
	mux.Handle("PUT /api/event_arrangements", cfg.authMiddleware(cfg.updateEventArrangement))

	// AUTH
	mux.HandleFunc("/pc/callback", cfg.planningcentercallback)
	mux.HandleFunc("POST /api/register", cfg.register)
	mux.HandleFunc("POST /api/login", cfg.login)

	server := &http.Server{
		Handler: mux,
		Addr: ":3005",
	  }
	  
	  fmt.Print("Running on Port 3005\n")
	  log.Fatal(server.ListenAndServe())
}