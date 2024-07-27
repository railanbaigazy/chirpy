package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type userRequest struct {
	Email string `json:"email"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userReq := userRequest{}
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprint(err))
		return
	}

	mailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !mailRegex.MatchString(userReq.Email) {
		respondWithError(w, http.StatusBadRequest, "invalid email")
		return
	}

	user, err := cfg.db.CreateUser(userReq.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprint(err))
		return
	}

	respondWithJSON(w, 201, user)
}
