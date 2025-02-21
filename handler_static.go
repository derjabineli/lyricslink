package main

import (
	"net/http"
	"os"
	"path/filepath"
)

func staticHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "." + r.URL.Path
	file, _ := os.ReadFile(filePath)

	ext := filepath.Ext(filePath)
	if ext == ".css" {
		w.Header().Add("Content-Type", "text/css")
	} else if ext == ".js" {
		w.Header().Add("Content-Type", "text/javascript")
	}
	w.Write(file)
 }