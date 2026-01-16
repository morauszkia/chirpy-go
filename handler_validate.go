package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"slices"
)

const MAX_CHIRP_LENGTH = 140
var BAD_WORDS = []string{"kerfuffle", "sharbert", "fornax"}

func filterBadWords(chirp string) string {
	words := strings.Split(chirp, " ")
	for i, word := range words {
		if slices.Contains(BAD_WORDS, strings.ToLower(word)) {
			words[i] = "****"
		} 

	}
	return strings.Join(words, " ")
}


func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body string `json:"body"`
	}
	type response struct {
		CleanedBody string `json:"cleaned_body"`
	}

	chirpBody := requestBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirpBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode request body", err)
		return
	}

	if len(chirpBody.Body) > MAX_CHIRP_LENGTH {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	filteredBody := filterBadWords(chirpBody.Body)
	
	payload := response{CleanedBody: filteredBody}
	respondWithJSON(w, http.StatusOK, payload)
}

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