package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Database reset not allowed", nil)
		return
	}
	cfg.fileserverHits.Store(0)
	if err := cfg.db.DeleteUsers(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete users", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Users table cleared \n"))
	w.Write([]byte("Metrics reset \n"))
}