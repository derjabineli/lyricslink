package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/derjabineli/lyricslink/internal/database"
	"github.com/google/uuid"
)

type PlanningCenterResponse struct {
	Links struct {
		Self string `json:"self"`
		Next string `json:"next"`
	} `json:"links"`
	Data     []PlanningCenterSong `json:"data"`
	Included []any  `json:"included"` // Kept generic since "included" is empty
	Meta     struct {
		TotalCount int `json:"total_count"`
		Count      int `json:"count"`
		Next       struct {
			Offset int `json:"offset"`
		} `json:"next"`
		CanOrderBy  []string `json:"can_order_by"`
		CanQueryBy  []string `json:"can_query_by"`
		Parent      struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"parent"`
	} `json:"meta"`
}

type PlanningCenterError struct {
	Errors []PlanningCenterErrorDetails	`json:"errors"`
}

type PlanningCenterErrorDetails struct {
	Code string	`json:"code"`
	Detail string	`json:"detail"`
}

type PCaccessParameters struct {
	GrantType string `json:"grant_type"`
	Code string `json:"code"`
	ClientID string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI string `json:"redirect_uri"`
}

type PCauthorizationParameters struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope string `json:"scope"`
	CreatedAt int `json:"created_at"`
}

type PCUserParameters struct {
	Data PCUserDataParameters `json:"data"`
}

type PCUserDataParameters struct {
	Attributes struct {
		LoginIdentifier string `json:"login_identifier"`
		Avatar string 			`json:"avatar"`
	} `json:"attributes"`
}

func (cfg *config) loginPC(w http.ResponseWriter, r *http.Request) {
	link := "https://api.planningcenteronline.com/oauth/authorize?client_id=" + cfg.pcClient + "&redirect_uri=" + cfg.pcRedirect + "&response_type=code&scope=services people"


	t, err := template.ParseFiles("./frontend/views/pc.html")
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}
	
	err = t.Execute(w, link)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

func (cfg *config) planningcentercallback(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.RequestURI)
	if err != nil {
		return
	}
	code := parsedURL.Query().Get("code")
	accessParams := PCaccessParameters {
		GrantType: "authorization_code",
		Code: code,
		ClientID: cfg.pcClient,
		ClientSecret: cfg.pcSecret,
		RedirectURI: cfg.pcRedirect,
	}

	requestBody, err := json.Marshal(accessParams)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusPermanentRedirect)
		return
	}

	client := &http.Client{}	
	authReq, err := http.NewRequest("POST", "https://api.planningcenteronline.com/oauth/token", bytes.NewBuffer(requestBody))
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusPermanentRedirect)
		return
	}
	authReq.Header.Add("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(authReq)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusPermanentRedirect)
		return
	}
	decoder := json.NewDecoder(resp.Body)
	authParams := PCauthorizationParameters{}
	decoder.Decode(&authParams)

	fmt.Printf("Access Token: %v\n", authParams.AccessToken)

	userID, err := cfg.getUserIDWithPCEmail(authParams.AccessToken)
	fmt.Print(userID)

	// go cfg.syncUserSongs(cfg.pcSongRoute, authParams.AccessToken, userID)
	// http.Redirect(w, r, "/dashboard", http.StatusPermanentRedirect)
}

func (cfg *config) getUserIDWithPCEmail(accessToken string) (uuid.UUID, error){
	method := "GET"
  
	client := &http.Client {
	}

	req, err := http.NewRequest(method, "https://api.planningcenteronline.com/people/v2/me", nil)
  
	if err != nil {
	  fmt.Println(err)
	  return uuid.Nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken))

	res, err := client.Do(req)
	if err != nil {
	  fmt.Println(err.Error())
	  return uuid.Nil, err
	}
	defer res.Body.Close()
  
	body, err := io.ReadAll(res.Body)
	if err != nil {
	  fmt.Println(err.Error())
	  return uuid.Nil, err
	}

	pcUserDetails := PCUserParameters{}
	err = json.Unmarshal(body, &pcUserDetails)
	if err != nil {
		return uuid.Nil, err
	}

	user, err := cfg.db.GetUserByEmail(context.Background(), pcUserDetails.Data.Attributes.LoginIdentifier)
	if err != nil {
		return uuid.Nil, err
	}

	if pcUserDetails.Data.Attributes.Avatar != "" {
		cfg.db.UpdateUserAvatar(context.Background(), database.UpdateUserAvatarParams{Avatar: sql.NullString{String: pcUserDetails.Data.Attributes.Avatar, Valid: true}})
	}

	return user.ID, nil
}