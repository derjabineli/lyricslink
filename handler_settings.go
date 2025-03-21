package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type userJSON struct {
	ID             uuid.UUID	`json:"id"`
    FirstName      string		`json:"firstName"`
    LastName       string		`json:"lastName"`
    Email          string		`json:"email"`
    CreatedAt      string	`json:"createdAt"`
    PcAuthorized   bool			`json:"pcAuthorized"`
}

func (cfg *config) handlerSettings(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	user, err := cfg.db.GetUserById(context.Background(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error finding user")
	}

	formattedDate := user.CreatedAt.Format("January 2, 2006")

	response := userJSON{
		ID: userID,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		CreatedAt: formattedDate,
		PcAuthorized: user.PcAuthorized,
	}

	t, err := template.ParseFiles("./frontend/views/settings.html")
	if err != nil {
		http.Error(w, "Error loading login page", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}

	err = t.Execute(w, response)
	if err != nil {
		http.Error(w, "Error rendering login page", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}