package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/derjabineli/lyricslink/internal/database"
	"github.com/google/uuid"
)

// Root struct representing the API response
type PlanningCenterSong struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes struct {
		Admin                    string    `json:"admin"`
		Author                   string    `json:"author"`
		CCLINumber               int       `json:"ccli_number"`
		Copyright                string    `json:"copyright"`
		CreatedAt                time.Time `json:"created_at"`
		Hidden                   bool      `json:"hidden"`
		LastScheduledAt          time.Time `json:"last_scheduled_at"`
		LastScheduledShortDates  string    `json:"last_scheduled_short_dates"`
		Notes                    *string   `json:"notes"` // Nullable field
		Themes                   string    `json:"themes"`
		Title                    string    `json:"title"`
		UpdatedAt                time.Time `json:"updated_at"`
	} `json:"attributes"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

type PlanningCenterSongDetails struct {
	Links struct {
		Self string `json:"self"`
	} `json:"links"`  
	Data []PlanningCenterArrangement `json:"data"`
}



func (cfg *config) syncUserSongs(url, accessToken string, userID uuid.UUID) {
	method := "GET"
  
	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)
  
	if err != nil {
	  fmt.Println(err)
	  return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken))
  
	res, err := client.Do(req)
	if err != nil {
	  fmt.Println(err.Error())
	  return
	}
	defer res.Body.Close()
  
	body, err := io.ReadAll(res.Body)
	if err != nil {
	  fmt.Println(err.Error())
	  return
	}
  
   pcSongData := PlanningCenterResponse{}
   json.Unmarshal(body, &pcSongData)
   for _, song := range pcSongData.Data {
	cfg.savePCSongToDB(&song, accessToken, userID)
   }
  
   if pcSongData.Links.Next != "" {
	  cfg.syncUserSongs(pcSongData.Links.Next, accessToken, userID)
   }
  }

func (cfg *config) savePCSongToDB(song *PlanningCenterSong, accessToken string, userID uuid.UUID) {
	pcId, err := strconv.Atoi(song.ID)
	if err != nil {
		return
	}

	songID, err := cfg.db.AddPCSong(context.Background(), database.AddPCSongParams{
		Title: song.Attributes.Title,
		Themes: validateSqlNullString(song.Attributes.Themes),
   		CopyRight: validateSqlNullString(song.Attributes.Copyright),
    	CcliNumber: validateSqlNullInt32(song.Attributes.CCLINumber),
    	Author: validateSqlNullString(song.Attributes.Author),
    	Admin: validateSqlNullString(song.Attributes.Admin),
    	PcID: validateSqlNullInt32(pcId),
	})
	if err != nil {
		fmt.Print(err)
		return
	}

	cfg.db.CreateUserSongRelation(context.Background(), database.CreateUserSongRelationParams{
		UserID: userID,
		SongID: songID,
	})


	url := fmt.Sprintf("https://api.planningcenteronline.com/services/v2/songs/%v/arrangements", song.ID)

	method := "GET"
  
	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)
  
	if err != nil {
	  fmt.Println(err.Error())
	  return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken))
  
  
	body, err := fetchFromPC(client, req) 
	if err != nil {
		return
	}

	arrangements := PlanningCenterSongDetails{}

	json.Unmarshal(body, &arrangements)

	for _, a := range arrangements.Data {
		cfg.savePCArrangementToDB(a, songID)
	}
}

func fetchFromPC(client *http.Client, req *http.Request) ([]byte, error) {
	res, err := client.Do(req)
	if err != nil {
	  fmt.Println(err.Error())
	  return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		retryAfter := res.Header.Get("Retry-After")
		if retryAfter != "" {
			// If Retry-After is in seconds, sleep for that amount of time
			sleepTime, err := strconv.Atoi(retryAfter)
			if err == nil {
				fmt.Printf("GO ROUTINE: Rate limit exceeded. Retrying after %d seconds...\n", sleepTime)
				time.Sleep(time.Duration(sleepTime) * time.Second)

				// Retry the request after sleeping
				res, err = client.Do(req)
				if err != nil {
					fmt.Printf("GO ROUTINE: Error executing retry request: %v\n", err)
					return nil, err
				}
				defer res.Body.Close()
				// Log status code after retry
				fmt.Printf("GO ROUTINE: Response Status after retry: %v\n", res.Status)
			} else {
				fmt.Println("GO ROUTINE: Invalid Retry-After header value")
			}
		}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
	  fmt.Println(err.Error())
	  return nil, err
	}
	return body, nil
}