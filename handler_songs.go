package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/derjabineli/lyricslink/internal/database"
	"github.com/google/uuid"
)

type searchParameters struct {
	Query string `json:"query"`
}

type songResultParameters struct {
	ID         uuid.UUID 		`json:"id"`
    PcID       int	`json:"pc_id"`
    Admin      string	`json:"admin"`
    Author     string	`json:"author"`
    CcliNumber int	`json:"ccli_number"`
    CopyRight  string	`json:"copy_right"`
    Themes     string	`json:"themes"`
    Title      string			`json:"title"`
}


func (cfg *config) getSongs(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized page")
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	searchParams := searchParameters{}
	decoder.Decode(&searchParams)

	dbSongResults, err := cfg.db.SearchSongs(context.Background(), database.SearchSongsParams{
		Title: searchParams.Query + "%", 
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	songResults := []songResultParameters{}
	for _, song := range dbSongResults {
		songResults = append(songResults, songResultParameters{
			ID: song.ID,
			PcID: getInt32Value(song.PcID),
			Admin: getSqlStringValue(song.Admin),
			Author: getSqlStringValue(song.Author),
			CcliNumber: getInt32Value(song.CcliNumber),
			CopyRight: getSqlStringValue(song.CopyRight),
			Themes: getSqlStringValue(song.Themes),
			Title: song.Title,
		})
	}

	respondWithJSON(w, 200, songResults)
}
