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
	payload := response{
		Chirp: Chirp{
			Id: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserId: chirp.UserID,
		},
	}
	respondWithJSON(w, http.StatusCreated, payload)
}