package main

import (
	"encoding/json"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userReq := userRequest{}
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	mailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !mailRegex.MatchString(userReq.Email) {
		respondWithError(w, http.StatusBadRequest, "invalid email")
		return
	}

	hashPassword := validatePassword(w, userReq.Password)
	if hashPassword == nil {
		return
	}

	user, err := cfg.db.CreateUser(userReq.Email, hashPassword)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, 201, user)
}

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr, err := getTokenString(w, r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	userID, err := getUserIDByToken(cfg, tokenStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	userReq := userRequest{}
	err = json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashPassword := validatePassword(w, userReq.Password)
	if hashPassword == nil {
		return
	}

	user, err := cfg.db.UpdateUser(userID, userReq.Email, hashPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func validatePassword(w http.ResponseWriter, password string) []byte {
	if len(password) < 5 {
		respondWithError(w, http.StatusBadRequest, "weak password")
		return nil
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return nil
	}
	return hashPassword
}
