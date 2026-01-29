package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/morauszkia/chirpy-go/internal/database"
)

type Chirp struct {
		Id			uuid.UUID	`json:"id"`
		CreatedAt	time.Time	`json:"created_at"`
		UpdatedAt	time.Time	`json:"updated_at"`
		Body		string		`json:"body"`
		UserId		uuid.UUID	`json:"user_id"`
	}

func formatChirpForResponse (chirp database.Chirp) Chirp {
	result := Chirp{
		Id: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserId: chirp.UserID,
	}
	return result
}

func (cfg *apiConfig) handlerCreateChirp (w http.ResponseWriter, r *http.Request) {
	type chirpBody struct {
		Body		string		`json:"body"`
		UserId		uuid.UUID	`json:"user_id"`
	}
	
	type response struct {
		Chirp
	}

	chirpData := chirpBody{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&chirpData); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body", err)
		return
	}
	defer r.Body.Close()

	validatedChirp, err := validateChirp(chirpData.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "", err)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: validatedChirp,
		UserID: chirpData.UserId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}
	payload := formatChirpForResponse(chirp)
	respondWithJSON(w, http.StatusCreated, payload)
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	payload := []Chirp{}
	for _, chirp := range chirps {
		payload = append(payload, formatChirpForResponse(chirp))
	}
	respondWithJSON(w, http.StatusOK, payload)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID", err)
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	payload := formatChirpForResponse(chirp)
	respondWithJSON(w, http.StatusOK, payload)
}