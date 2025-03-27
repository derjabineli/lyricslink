package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/google/uuid"
)

type LiveSetJSON struct {
	Songs []LiveSongJSON `json:"songs"`
}

type LiveSongJSON struct {
	Title string 		`json:"title"`
	Author string 		`json:"author"`
	Lyrics string		`json:"lyrics"`
	ChordChart string 	`json:"chord_chart"`
	CCLI int			`json:"ccli"`
	CopyRight string	`json:"copy_right"`
}

func (cfg *config) handlerLive(w http.ResponseWriter, r *http.Request) {
	eventQuery := path.Base(r.URL.Path)

	eventID, err := uuid.Parse(eventQuery)
	if err != nil {
		fmt.Println(err)
	}

	arrangements, err := cfg.db.GetArrangementsAndSongsWithEventId(context.Background(), eventID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "No arrangements found")
		return
	}

	data := LiveSetJSON{}

	for _, a := range arrangements {
		song := LiveSongJSON{Title: a.Title, Author: getSqlStringValue(a.Author), Lyrics: strings.ReplaceAll(a.Lyrics, "\n", "<br>"), ChordChart: getSqlStringValue(a.ChordChart), CCLI: getInt32Value(a.CcliNumber), CopyRight: getSqlStringValue(a.CopyRight)}

		data.Songs = append(data.Songs, song)
	}

	liveJSON, err := json.Marshal(data)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve lyrics")
		return
	}

	eventData := struct {
		Data template.JS
	}{
		Data: template.JS(liveJSON),
	}

	t, err := template.ParseFiles("./frontend/views/live.html")
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