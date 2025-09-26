package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/derjabineli/lyricslink/internal/auth"
	"github.com/derjabineli/lyricslink/internal/database"
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

type PCUserParameters struct {
	Data PCUserDataParameters `json:"data"`
}

type PCUserDataParameters struct {
	ID string 					`json:"id"`
	Attributes struct {
		LoginIdentifier string 	`json:"login_identifier"`
		Avatar string 			`json:"avatar"`
		FirstName string		`json:"first_name"`
		LastName string 		`json:"last_name"`
		SiteAdministrator bool 	`json:"site_administrator"`
	} `json:"attributes"`
	Links struct {
		Organization string 	`json:"organization"`
	} `json:"links"`
}

type PCOrganizationParameters struct {
	Data struct {
		Type string `json:"type"`
		ID 	 string 	`json:"id"`
		Attributes struct {
			AvatarURL 		string	`json:"avatar_url"`
			ContactWebsite	string	`json:"contact_website"`
			CountryCode		string	`json:"country_code"`
			CreatedAt 		string	`json:"created_at"`
			DateFormat 		string 	`json:"date_format"`
			Name 			string 	`json:"name"`
			TimeZone 		string 	`json:"time_zone"`
		} `json:"attributes"`
	} `json:"data"`
}

func (cfg *config) loginStatic(w http.ResponseWriter, r *http.Request) {
	loginLink := fmt.Sprintf("https://api.planningcenteronline.com/oauth/authorize?client_id=%v&redirect_uri=%v&response_type=code&scope=services people", cfg.pcClient, cfg.pcRedirect)

	t, err := template.ParseFiles("./frontend/views/login.html")
	if err != nil {
		http.Error(w, "Error loading login page", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}

	err = t.Execute(w, loginLink)
	if err != nil {
		http.Error(w, "Error rendering login page", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

func (cfg *config) loginWithPC(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.RequestURI)
	if err != nil {
		return
	}
	code := parsedURL.Query().Get("code")
	accessParams := PCOAuthAuthorizationCodeRequest {
		GrantType: "authorization_code",
		Code: code,
		ClientID: cfg.pcClient,
		ClientSecret: cfg.pcSecret,
		RedirectURI: cfg.pcRedirect,
	}

	requestBody, err := json.Marshal(accessParams)
	if err != nil {
		fmt.Printf("Error when marshalling access params req body.\n err: %v\n", err.Error())
		http.Redirect(w, r, "/error", http.StatusPermanentRedirect)
		return
	}

	authParams, err := getPCToken(requestBody)
	if err != nil {
		fmt.Printf("Error when getting PC Token.\n err: %v\n", err.Error())
		http.Redirect(w, r, "/error", http.StatusPermanentRedirect)
		return
	}

	pcUserData, err := getPCUserData(authParams.AccessToken)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	user, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		FirstName: pcUserData.Data.Attributes.FirstName,
		LastName:  pcUserData.Data.Attributes.LastName,
		Email:  pcUserData.Data.Attributes.LoginIdentifier,
		Avatar:  pcUserData.Data.Attributes.Avatar,
		PcID:  pcUserData.Data.ID,
		Administrator:  pcUserData.Data.Attributes.SiteAdministrator,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	organization, err := cfg.getPCOrganizationData(authParams.AccessToken, pcUserData.Data.Links.Organization)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = cfg.db.CreateUserOrgRelation(context.Background(), database.CreateUserOrgRelationParams{UserID: user.ID, OrganizationID: organization.ID})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	session, err := cfg.db.AddUserSession(context.Background(), database.AddUserSessionParams{
		UserID: user.ID, 
		AccessToken: authParams.AccessToken, 
		TokenType: authParams.TokenType, 
		ExpiresIn: int32(authParams.ExpiresIn), 
		RefreshToken: authParams.RefreshToken, 
		Scope: authParams.Scope})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	accessCookie, err := auth.NewAccessTokenCookie(user.ID, session.ID, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was an issue logging you in. Please try again")
		return
	}
	refreshCookie, err := auth.NewRefreshTokenCookie(user.ID, session.ID, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was an issue logging you in. Please try again")
		return
	}

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)

	http.Redirect(w, r, "/dashboard", http.StatusPermanentRedirect)
}

func getPCUserData(bearerToken string) (PCUserParameters, error) {
	client := &http.Client{}	
	req, err := http.NewRequest(http.MethodGet, "https://api.planningcenteronline.com/people/v2/me", nil)
	if err != nil {
		return PCUserParameters{}, err
	}
	req.Header.Add("Authorization", "Bearer " + bearerToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Couldn't make request %v\n", err.Error())
		return PCUserParameters{}, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	userParams := PCUserParameters{}
	err = decoder.Decode(&userParams)
	if err != nil {
		fmt.Printf("Couldn't decode %v\n", err.Error())
		return PCUserParameters{}, err
	}

	return userParams, nil
}

func (cfg *config) getPCOrganizationData(bearerToken string, org_url string) (database.PlanningCenterOrganization, error) {
	client := &http.Client{}	
	
	if org_url == "" {
		org_url = "https://api.planningcenteronline.com/people/v2/organization"
	}

	req, err := http.NewRequest(http.MethodGet, org_url, nil)
	if err != nil {
		return database.PlanningCenterOrganization{}, err
	}
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Organization couldn't make request %v\n", err.Error())
		return database.PlanningCenterOrganization{}, err
	}
	decoder := json.NewDecoder(resp.Body)
	orgParams := PCOrganizationParameters{}
	err = decoder.Decode(&orgParams)
	if err != nil {
		fmt.Printf("Organization couldn't decode %v\n", err.Error())
		return database.PlanningCenterOrganization{}, err
	}	
	
	organization, err := cfg.db.CreatePlanningCenterOrganization(context.Background(), database.CreatePlanningCenterOrganizationParams{PcID: orgParams.Data.ID, Name: orgParams.Data.Attributes.Name})
	if err != nil {
		return database.PlanningCenterOrganization{}, err
	}

	return organization, err
}
