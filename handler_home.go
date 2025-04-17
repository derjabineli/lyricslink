package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func (cfg *config) handlerHome(w http.ResponseWriter, r *http.Request) {
	loginLink := fmt.Sprintf("https://api.planningcenteronline.com/oauth/authorize?client_id=%v&redirect_uri=%v&response_type=code&scope=services people", cfg.pcClient, cfg.pcRedirect)

	t, err := template.ParseFiles("./frontend/views/home.html")
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