package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprint(err))
		return
	}
	sort.Slice(chirps, func(i, j int) bool { return chirps[i].ID < chirps[j].ID })
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("chirpid"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id")
		return
	}
	chirp, err := cfg.db.GetChirpByID(id)
	if err != nil {
		respondWithError(w, 404, fmt.Sprint(err))
		return
	}
	respondWithJSON(w, http.StatusOK, chirp)
}
