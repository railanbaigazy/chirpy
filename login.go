package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	loginReq := loginRequest{}
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprint(err))
		return
	}

	user, err := cfg.db.Login(loginReq.Email, loginReq.Password)
	if err != nil {
		respondWithError(w, 401, fmt.Sprint(err))
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}
