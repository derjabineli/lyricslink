package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
		return
	}

	organizationID, err := cfg.db.GetOrganizationIDByUserID(r.Context(), userID)
	if err != nil {
		fmt.Printf("Couldn't get Org ID.\n error: %v\n", err.Error())
		respondWithError(w, http.StatusUnauthorized, "Internal Server Error")
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	searchParams := searchParameters{}
	decoder.Decode(&searchParams)

	dbSongResults, err := cfg.db.SearchSongs(context.Background(), database.SearchSongsParams{
		Title: searchParams.Query + "%", 
		OrganizationID: organizationID,
	})
	if err != nil {
		fmt.Printf("Couldn't get songs.\n error: %v\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
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

func (cfg *config) syncPlanningCenterSongs(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	ctx := context.Background()

	organizationID, err := cfg.db.GetOrganizationIDByUserID(ctx, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "there was an error syncing your songs")
		return
	}

	token, err := cfg.db.GetTokenByUserID(ctx, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "there was an error syncing your songs")
		return
	}

	tokenExpiration := token.UpdatedAt.Add(time.Second*time.Duration(token.ExpiresIn))
	if time.Now().After(tokenExpiration){
		fmt.Println("Token expired!")
		return
	}

	err = cfg.fetchandSyncSongs(cfg.pcSongRoute, token.AccessToken, organizationID)
	if err != nil {
		fmt.Printf("Error encountered error: %v\n", err.Error())
		respondWithError(w, http.StatusUnauthorized, "Unauthorized Planning Center Account")
		return
	}
	respondWithJSON(w, http.StatusAccepted, jsonServerResponse{Success: "Successfully synced Planning Center Songs"})
}