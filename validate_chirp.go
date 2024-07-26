package main

import (
	"encoding/json"
	"net/http"
)

type ChirpReq struct {
	Body string `json:"body"`
}

type ErrorResp struct {
	Error string `json:"error"`
}

type SuccessResp struct {
	Valid bool `json:"valid"`
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	chirp := ChirpReq{}
	err := json.NewDecoder(r.Body).Decode(&chirp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResp{Error: "Invalid request body"})
		return
	}

	if len(chirp.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResp{Error: "Chirp is too long"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResp{Valid: true})
}
