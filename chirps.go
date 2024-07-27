package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type chirpRequest struct {
	Body string `json:"body"`
}

var profanes []string = []string{"kerfuffle", "sharbert", "fornax"}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	chirpReq := chirpRequest{}
	err := json.NewDecoder(r.Body).Decode(&chirpReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(chirpReq.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := cleanText(chirpReq.Body)
	chirp, err := cfg.db.CreateChirp(cleanedBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprint(err))
	}
	respondWithJSON(w, 201, chirp)
}

func cleanText(body string) string {
	words := strings.Fields(body)
	cleanedBody := []string{}
	for _, word := range words {
		isProfane := false
		for _, profane := range profanes {
			if strings.ToLower(word) == profane {
				cleanedBody = append(cleanedBody, "****")
				isProfane = true
				break
			}
		}
		if !isProfane {
			cleanedBody = append(cleanedBody, word)
		}
	}
	return strings.Join(cleanedBody, " ")
}

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
