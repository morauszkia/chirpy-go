package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID		`json:"id"`
	CreatedAt time.Time		`json:"created_at"`
	UpdatedAt time.Time		`json:"updated_at"`
	Email     string		`json:"email"`
	}

func (cfg *apiConfig) handlerCreateUser (w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email	string	`json:"email"`
	}
	type response struct {
		User
	}

	userData := requestBody{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userData); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), userData.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	payload := response{
		User: User {
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		},
	}
	respondWithJSON(w, http.StatusCreated, payload)
}