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
	tokenDuration time.Duration
}


func main() {
	godotenv.Load()
	pcClient := os.Getenv("PC_CLIENTID")
	pcSecret := os.Getenv("PC_SECRET")
	pcRedirect := os.Getenv("PC_REDIRECT")
	tokenSecret := os.Getenv("TOKEN_SECRET")
	dbURL := os.Getenv("DB_URL")
  	db, err := sql.Open("postgres", dbURL)
  	if err != nil {
   		fmt.Print(err)
    	return 
  	}
	dbQueries := database.New(db)

	cfg := config{db: dbQueries, pcClient: pcClient, pcSecret: pcSecret, pcRedirect: pcRedirect, tokenSecret: tokenSecret, tokenDuration: time.Hour * 8}
	

	mux := http.NewServeMux()

	// PAGES
	mux.Handle("/", cfg.guestOnlyMiddleware(handlerHome))
	mux.Handle("/login", cfg.guestOnlyMiddleware(cfg.loginStatic))
	mux.HandleFunc("/static/", staticHandler)
	mux.Handle("/dashboard", cfg.authMiddleware(cfg.handlerDashboard))
	
	// API
	mux.HandleFunc("POST /api/register", cfg.register)
	mux.HandleFunc("POST /api/login", cfg.login)
	mux.Handle("POST /api/events", cfg.authMiddleware(cfg.addEvent))

	// AUTH
	mux.HandleFunc("/pc/callback", cfg.planningcentercallback)

	server := &http.Server{
		Handler: mux,
		Addr: ":3005",
	  }
	  
	  fmt.Print("Running on Port 3005")
	  log.Fatal(server.ListenAndServe())
}