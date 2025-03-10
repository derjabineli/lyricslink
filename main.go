package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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

	cfg := config{db: dbQueries, pcClient: pcClient, pcSecret: pcSecret, pcRedirect: pcRedirect, tokenSecret: tokenSecret}
	

	mux := http.NewServeMux()

	// PAGES
	mux.HandleFunc("/", handlerHome)
	mux.HandleFunc("/login", cfg.loginStatic)
	mux.HandleFunc("/static/", staticHandler)
	mux.HandleFunc("/dashboard", handlerDashboard)
	
	// API
	mux.HandleFunc("/api/register", cfg.register)
	mux.HandleFunc("/api/login", cfg.login)

	// AUTH
	mux.HandleFunc("/pc/callback", cfg.planningcentercallback)

	server := &http.Server{
		Handler: mux,
		Addr: ":3005",
	  }
	  
	  fmt.Print("Running on Port 3005")
	  log.Fatal(server.ListenAndServe())
}