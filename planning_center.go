package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

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
	ID string 					`json:"id"`
	Attributes struct {
		LoginIdentifier string 	`json:"login_identifier"`
		Avatar string 			`json:"avatar"`
	} `json:"attributes"`
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

	userID, err := cfg.getUserIDWithPCDetails(authParams.AccessToken)

	if err != nil {
		redirectURL := fmt.Sprintf("/settings?status=error&message=%v", err.Error())
		http.Redirect(w, r, redirectURL, http.StatusPermanentRedirect)
		return
	}

	go cfg.syncUserSongs(cfg.pcSongRoute, authParams.AccessToken, userID)

	redirectURL := fmt.Sprintf("/settings?status=success&message=%v", "successfully synced account. your songs may take a few minutes to appear")
	http.Redirect(w, r, redirectURL, http.StatusPermanentRedirect)
}

func (cfg *config) getUserIDWithPCDetails(accessToken string) (uuid.UUID, error){
	method := "GET"
  
	client := &http.Client {
	}

	req, err := http.NewRequest(method, "https://api.planningcenteronline.com/people/v2/me", nil)
  
	if err != nil {
	  fmt.Println(err)
	  return uuid.Nil, errors.New("there was an error retrieving your planning center information. please try again")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken))

	res, err := client.Do(req)
	if err != nil {
	  fmt.Println(err.Error())
	  return uuid.Nil, errors.New("there was an error retrieving your planning center information. please try again")
	}
	defer res.Body.Close()
  
	body, err := io.ReadAll(res.Body)
	if err != nil {
	  fmt.Println(err.Error())
	  return uuid.Nil, errors.New("there was an error retrieving your planning center information. please try again")
	}

	pcUserDetails := PCUserParameters{}
	err = json.Unmarshal(body, &pcUserDetails)
	if err != nil {
		return uuid.Nil, errors.New("there was an error retrieving your planning center information. please try again")
	}

	planningCenterID, err := strconv.Atoi(pcUserDetails.Data.ID)
	if err != nil {
		return uuid.Nil, errors.New("there was an error retrieving your planning center information. please try again")
	}

	user, err := cfg.db.GetUserByPCID(context.Background(), validateSqlNullInt32(planningCenterID))
	if errors.Is(err, sql.ErrNoRows) {
		user, err = cfg.db.GetUserByEmail(context.Background(), pcUserDetails.Data.Attributes.LoginIdentifier)
		if err != nil {
			return uuid.Nil, errors.New("your account couldn't be synced. please ensure that your planning center account login email matches your lyriclink email")
		}
	}

	cfg.db.UpdatePlanningCenterUser(context.Background(), database.UpdatePlanningCenterUserParams{
		Avatar: validateSqlNullString(pcUserDetails.Data.Attributes.Avatar),
		PcID: validateSqlNullInt32(planningCenterID),
		ID: user.ID,
	})

	return user.ID, nil
}

func redirect_after_pc_sync(w http.ResponseWriter, r *http.Request) {
	redirectStatus := r.URL.Query().Get("status")
	redirectMessage := r.URL.Query().Get("message")

	fmt.Print(redirectMessage, redirectStatus)

	
	http.Redirect(w, r, "/settings", http.StatusFound)
}