package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/derjabineli/lyricslink/internal/database"
)

type formattedEvent struct {
	Link string `json:"link"`
	Name string `json:"name"`
	Date string	`json:"date"`
	ID string	`json:"id"`
}

type userParameters struct {
	Avatar string	`json:"avatar"`
}

type dashboardParameters struct {
	User userParameters		`json:"user"`
	Events []formattedEvent `json:"events"`
}

func (cfg *config) handlerDashboard(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	user, err := cfg.db.GetUserById(context.Background(), userID)
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
	data := dashboardParameters{
		User: userParameters{
			Avatar: user.Avatar,
		},
		Events: formattedEvents,
	}
	dashboardJSON, _ := json.Marshal(data)

	eventData := struct {
		Data template.JS
	}{
		Data: template.JS(dashboardJSON),
	}

	t, err := template.ParseFiles("./frontend/views/dashboard.html")
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}
	
	err = t.Execute(w, eventData)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

func formatEvents(events []database.Event) []formattedEvent {
	formattedEvents := []formattedEvent{}

	for _, event := range events {
		link := fmt.Sprintf("events/%v", event.ID)

		formattedEvents = append(formattedEvents, formattedEvent{
			Link: link,
			Name: event.Name,
			Date: event.Date.Format("January 2, 2006"),
			ID: event.ID.String(),
		})
	}

	return formattedEvents
} 