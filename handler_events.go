package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/derjabineli/lyricslink/internal/database"
)

type eventParameters struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

func (cfg *config) addEvent(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	eventParams := eventParameters{}
	decoder.Decode(&eventParams)

	formattedTime, _ := time.Parse("2006-01-02", eventParams.Date)

	event, err := cfg.db.CreateEvent(context.Background(), database.CreateEventParams{Name: eventParams.Name, Date: formattedTime, UserID: userID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create event")
	}

	respondWithJSON(w, http.StatusCreated, event)
}