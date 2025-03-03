package main

import (
	"html/template"
	"log"
	"net/http"
)

func handlerDashboard(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./template/dashboard.html")
		if err != nil {
			http.Error(w, "Error loading login page", http.StatusInternalServerError)
			log.Println("Template parsing error:", err)
			return
		}
	
		err = t.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error rendering login page", http.StatusInternalServerError)
			log.Println("Template execution error:", err)
		}
}