package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/derjabineli/lyricslink/internal/database"
)

type PCOAuthAuthorizationCodeRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

type PCOAuthAuthorizationResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int    `json:"created_at"`
}

type PCOAuthRefreshRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	GrantType    string `json:"grant_type"`
}

type PCOAuthRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    int    `json:"created_at"`
}

func getPCToken(requestBody []byte) (PCOAuthAuthorizationResponse, error) {
	authParams := PCOAuthAuthorizationResponse{}

	client := &http.Client{}
	authReq, err := http.NewRequest(http.MethodPost, "https://api.planningcenteronline.com/oauth/token", bytes.NewBuffer(requestBody))
	if err != nil {
		return authParams, err
	}
	authReq.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(authReq)
	if err != nil {
		return authParams, err
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&authParams)
	if err != nil {
		fmt.Println(err.Error())
	}
	return authParams, nil
}

func (cfg *config) refreshSessionToken(token database.UserSession) (string, error) {
	oAuthRefreshRequest := PCOAuthRefreshRequest{
		ClientID:     cfg.pcClient,
		ClientSecret: cfg.pcSecret,
		RefreshToken: token.RefreshToken,
		GrantType:    "refresh_token",
	}
	jsonRequest, err := json.Marshal(oAuthRefreshRequest)
	if err != nil {
		return "", err
	}

	authParams, err := getPCToken(jsonRequest)
	if err != nil {
		return "", err
	}

	newToken, err := cfg.db.UpdateUserToken(context.Background(),
		database.UpdateUserTokenParams{
			AccessToken:  authParams.AccessToken,
			RefreshToken: authParams.RefreshToken,
			Scope:        authParams.Scope,
			ID:           token.ID,
	})
	if err != nil {
		return "", err
	}

	return newToken.AccessToken, nil
}

