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

type EventViewData struct {
	ID       uuid.UUID        `json:"id"`
	Name     string           `json:"name"`
	Date     string           `json:"date"`
	Songs    []songParameters `json:"songs"`
	Livelink string           `json:"live_link"`
	User     UserViewData     `json:"user"`
}

type songParameters struct {
	ID                   uuid.UUID               `json:"id"`
	Event_Arrangement_Id uuid.UUID               `json:"eventArrangementId"`
	PC_ID                int                     `json:"pcId"`
	Admin                string                  `json:"admin"`
	Author               string                  `json:"author"`
	CCLI                 int                     `json:"ccli"`
	Copyright            string                  `json:"copyright"`
	Themes               string                  `json:"themes"`
	Title                string                  `json:"title"`
	Arrangements         []arrangementParameters `json:"arrangements"`
}

type arrangementParameters struct {
	ID         uuid.UUID     `json:"id"`
	Name       string        `json:"name"`
	Lyrics     template.HTML `json:"lyrics"`
	ChordChart string        `json:"chordChart"`
	SongID     uuid.UUID     `json:"songId"`
	IsSelected bool          `json:"isSelected"`
}

type updateEventParameters struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Date string    `json:"date"`
}

func (cfg *config) handlerEvents(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}
	// Needed to pass avatar url to frontend
	user, err := cfg.db.GetUserById(context.Background(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't fetch user")
		return
	}

	eventQuery := path.Base(r.URL.Path)

	eventID, err := uuid.Parse(eventQuery)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't find event")
		return
	}

	event, err := cfg.db.GetEventById(context.Background(), eventID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error finding event with id %v", eventID))
		return
	}
	formattedDate := event.Date.Format("2006-01-02")

	eventArrangments, err := cfg.db.GetArrangementsWithEventId(context.Background(), eventID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error finding arrangements with event id %v", eventID))
		return
	}

	// Live link to public lyric sheet
	livelink := fmt.Sprintf("/live/%v", eventID)
	eventParams := EventViewData{ID: eventID, Name: event.Name, Date: formattedDate, Livelink: livelink, Songs: []songParameters{}}


	for _, selectedArrangement := range eventArrangments {
		dbSong, err := cfg.db.GetSongById(context.Background(), selectedArrangement.SongID)
		if err != nil {
			fmt.Printf("Error: finding song with id %v\n", selectedArrangement.SongID)
			fmt.Printf("Error: %v\n", err)
			continue
		}
		song := songParameters{
			ID:           selectedArrangement.SongID,
			PC_ID:        int(dbSong.PcID.Int32),
			Event_Arrangement_Id: selectedArrangement.ID,
			Title:        dbSong.Title,
			Arrangements: []arrangementParameters{},
		}
		eventParams.Songs = append(eventParams.Songs, song)
		idx := len(eventParams.Songs) - 1

		available_arrangements, err := cfg.db.GetAvailableArrangements(context.Background(), database.GetAvailableArrangementsParams{SongID: dbSong.ID, ID: selectedArrangement.ArrangementID})
		if err != nil {
			fmt.Printf("Error: finding arrangements for song with id %v\n", selectedArrangement.SongID)
			fmt.Printf("Error: %v\n", err)
			continue
		}

		for _, a := range available_arrangements {
			song.Arrangements = append(song.Arrangements, arrangementParameters{
			ID:         a.ID,
			Name:       a.Name,
			Lyrics:     lyricSheetToHTML(a.Lyrics),
			ChordChart: a.ChordChart.String,
			SongID:     a.SongID,
			IsSelected: false,
		})
		}


		song.Arrangements = append(song.Arrangements, arrangementParameters{
			ID:         selectedArrangement.ID,
			Name:       selectedArrangement.ArrangementName,
			Lyrics:     lyricSheetToHTML(selectedArrangement.Lyrics),
			ChordChart: selectedArrangement.ChordChart.String,
			SongID:     selectedArrangement.SongID,
			IsSelected: true,
		})

		eventParams.Songs[idx] = song
	}

	eventParams.User.Avatar = user.Avatar

	eventJSON, _ := json.Marshal(eventParams)

	t, err := template.ParseFiles(path.Join("frontend", "views", "event.html"))
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}

	eventData := struct {
		Event template.JS
	}{
		Event: template.JS(eventJSON),
	}

	err = t.Execute(w, eventData)
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

func (cfg *config) updateEventDate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	eventParams := updateEventParameters{}
	decoder.Decode(&eventParams)

	formattedTime, _ := time.Parse("2006-01-02", eventParams.Date)

	updatedEvent, err := cfg.db.UpdateEventDate(context.Background(), database.UpdateEventDateParams{ID: eventParams.ID, Date: formattedTime})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update event")
	}

	respondWithJSON(w, http.StatusOK, updatedEvent)
}

func (cfg *config) deleteEvent(w http.ResponseWriter, r *http.Request) {
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

	err = cfg.db.DeleteEvent(context.Background(), database.DeleteEventParams{ID: eventID, UserID: userID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete event")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

