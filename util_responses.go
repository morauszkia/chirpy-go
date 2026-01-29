package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, statusCode int, message string, err error) {
	if err != nil {
		log.Println(err)
	}
	if statusCode >= 500 {
		log.Printf("Internal server error: %s", message)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	response := errorResponse{Error: message}
	respondWithJSON(w, statusCode, response)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	encodedResponse, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error encoding JSON response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(encodedResponse)
}