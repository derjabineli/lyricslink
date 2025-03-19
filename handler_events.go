package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/derjabineli/lyricslink/internal/database"
	"github.com/google/uuid"
)

type newEventParameters struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type eventParameters struct {
	Name string
	Date string
	Songs map[uuid.UUID]songParameters
}

type songParameters struct {
	ID uuid.UUID
	PC_ID int
	Admin string
	Author string
	CCLI int
	Copyright string
	Themes string
	Title string
	Arrangements []arrangementParameters
}

type arrangementParameters struct {
	ID uuid.UUID
	Name string
	Lyrics template.HTML
	ChordChart string
	SongID uuid.UUID
	IsSelected bool
}

func (cfg *config) handlerEvents(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Print(userID)
	eventQuery := path.Base(r.URL.Path)

	eventID, err := uuid.Parse(eventQuery)
	if err != nil {
		fmt.Println(err)
	}
	
	event, err := cfg.db.GetEventById(context.Background(), eventID)
	if err != nil {
		fmt.Println(err)
	}
	formattedDate := event.Date.Format("2006-01-02")

	arrangements, _ := cfg.db.GetArrangementsWithEventId(context.Background(), eventID)

	eventParams := eventParameters{Name: event.Name, Date: formattedDate, Songs: map[uuid.UUID]songParameters{}}

	for _, a := range arrangements {
		song, exists := eventParams.Songs[a.SongID]
		if !exists {
			dbSong, _ := cfg.db.GetSongById(context.Background(), a.SongID)

			song = songParameters{
				ID:          a.SongID,
				PC_ID: int(dbSong.PcID.Int32),
				Title: dbSong.Title,
				Arrangements: []arrangementParameters{},
			}
		}
	
		song.Arrangements = append(song.Arrangements,  arrangementParameters{
			ID:         a.ID,
			Name:       a.Name,
			Lyrics:     lyricSheetToHTML(a.Lyrics),
			ChordChart: a.ChordChart.String,
			SongID:     a.SongID, 
			IsSelected: a.IsSelected,
		})
	
		eventParams.Songs[a.SongID] = song
	}

	t, err := template.ParseFiles(path.Join("frontend", "views", "event.html"))
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}
	
	err = t.Execute(w, eventParams)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

func (cfg *config) addEvent(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	eventParams := newEventParameters{}
	decoder.Decode(&eventParams)

	formattedTime, _ := time.Parse("2006-01-02", eventParams.Date)

	event, err := cfg.db.CreateEvent(context.Background(), database.CreateEventParams{Name: eventParams.Name, Date: formattedTime, UserID: userID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create event")
		return
	}

	respondWithJSON(w, http.StatusCreated, event)
}