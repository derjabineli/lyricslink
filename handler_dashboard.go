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

type FormattedEvent struct {
	Link string `json:"link"`
	Name string `json:"name"`
	Date string	`json:"date"`
	ID string	`json:"id"`
}

type UserViewData struct {
	Avatar string	`json:"avatar"`
}

type DashboardViewData struct {
	User UserViewData		`json:"user"`
	Events []FormattedEvent `json:"events"`
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
	data := DashboardViewData{
		User: UserViewData{
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

func formatEvents(events []database.Event) []FormattedEvent {
	formattedEvents := []FormattedEvent{}

	for _, event := range events {
		link := fmt.Sprintf("events/%v", event.ID)

		formattedEvents = append(formattedEvents, FormattedEvent{
			Link: link,
			Name: event.Name,
			Date: event.Date.Format("January 2, 2006"),
			ID: event.ID.String(),
		})
	}

	return formattedEvents
} 