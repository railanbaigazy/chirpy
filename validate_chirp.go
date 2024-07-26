package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ChirpReq struct {
	Body string `json:"body"`
}

type SuccessResp struct {
	CleanedBody string `json:"cleaned_body"`
}

var profanes []string = []string{"kerfuffle", "sharbert", "fornax"}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	chirp := ChirpReq{}
	err := json.NewDecoder(r.Body).Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := cleanText(chirp.Body)
	respondWithJSON(w, http.StatusOK, SuccessResp{CleanedBody: cleanedBody})
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
