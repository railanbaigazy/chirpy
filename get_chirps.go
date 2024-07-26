package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprint(err))
		return
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
