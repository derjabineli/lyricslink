package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func staticHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "./frontend" + r.URL.Path
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Print(err)
	}

	ext := filepath.Ext(filePath)
	if ext == ".css" {
		w.Header().Add("Content-Type", "text/css")
	} else if ext == ".js" {
		w.Header().Add("Content-Type", "text/javascript")
	}
	w.Write(file)
 }