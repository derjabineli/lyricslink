package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/derjabineli/lyricslink/internal/database"
	"github.com/google/uuid"
)

type arrangementsResponse struct {
	ID 			uuid.UUID 	`json:"id"`
	Name 		string		`json:"name"`
	Lyrics 		string 		`json:"lyrics"`
	ChordChart 	string		`json:"chord_chart"`
	SongID 		uuid.UUID	`json:"song_id"`
}

type eventArrangementBody struct {
	EventID 		uuid.UUID	`json:"event_id"`
	ArrangementID 		uuid.UUID	`json:"arrangement_id"`
}

func (cfg *config) getArrangements(w http.ResponseWriter, r *http.Request) {
	songPath := path.Base(r.URL.Path)
	songID, err := uuid.Parse(songPath)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "No resources available")
	}

	arrangements, err := cfg.db.GetArrangementWithSongId(context.Background(), songID) 
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error")
	}

	response := []arrangementsResponse{}
	for _, a := range arrangements {
		arrangement := arrangementsResponse{
			ID: a.ID,
			Name: a.Name,
			Lyrics: a.Lyrics,
			ChordChart: getSqlStringValue(a.ChordChart),
			SongID: a.SongID,
		}

		response = append(response, arrangement)
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *config) addArrangementToEvent(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	eventArrangement := eventArrangementBody{}
	decoder.Decode(&eventArrangement)

	_, err := cfg.db.AddArrangementToEvent(context.Background(), database.AddArrangementToEventParams{
		EventID: eventArrangement.EventID,
		ArrangementID: eventArrangement.ArrangementID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't add arrangement to event")
		return
	}
}

func (cfg *config) updateEventArrangement(w http.ResponseWriter, r *http.Request) {
	eventPath := path.Base(r.URL.Path)
	eventArrangementID, err := uuid.Parse(eventPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "No event id found")
		return
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	eventArrangement := eventArrangementBody{}
	decoder.Decode(&eventArrangement)
	
	_, err = cfg.db.UpdateArrangement(context.Background(), database.UpdateArrangementParams{
		ID: eventArrangementID, ArrangementID: eventArrangement.ArrangementID,
	})
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couldn't update arrangement")
		return
	}

	w.WriteHeader(200)
}

func (cfg *config) deleteEventArrangement(w http.ResponseWriter, r *http.Request) {
	eventPath := path.Base(r.URL.Path)
	eventID, err := uuid.Parse(eventPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "No event id found")
		return
	}

	userID, err := getUserIDFromContext(r)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Print(userID)

	err = cfg.db.DeleteEventArrangement(context.Background(), eventID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete event")
		return 
	}

	w.WriteHeader(http.StatusNoContent)
}