package main

import (
	"log"
	"net/http"
)


func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlerStart)

	server := &http.Server{
		Handler: mux,
		Addr: ":8080",
	  }
	  
	  log.Fatal(server.ListenAndServe())
}