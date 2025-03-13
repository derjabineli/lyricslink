package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/derjabineli/lyricslink/internal/database"
)

type formattedEvent struct {
	Name string
	Date string
	ID string
}

func (cfg *config) handlerDashboard(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	events, err := cfg.db.GetEventsByUserId(context.Background(), userID)
	if err != nil {
		fmt.Print("Couldn't get events")
		return
	}

	formattedEvents := formatEvents(events)

	t, err := template.ParseFiles("./views/dashboard.html")
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}
	
	err = t.Execute(w, formattedEvents)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

func formatEvents(events []database.Event) []formattedEvent {
	formattedEvents := []formattedEvent{}

	for _, event := range events {
		formattedEvents = append(formattedEvents, formattedEvent{
			Name: event.Name,
			Date: event.Date.Format("January 2, 2006"),
			ID: event.ID.String(),
		})
	}

	return formattedEvents
} 