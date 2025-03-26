package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type jsonErrorResponse struct {
    Error string `json:"error"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
    response, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(code)

    _, err = w.Write(response)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
        return err
	}
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
    return respondWithJSON(w, code, jsonErrorResponse{Error: msg})
}