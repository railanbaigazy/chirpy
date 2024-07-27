package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
		respondWithError(w, http.StatusBadRequest, fmt.Sprint(err))
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
		respondWithError(w, http.StatusBadRequest, fmt.Sprint(err))
		return
	}

	respondWithJSON(w, 201, user)
}

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	headerText := r.Header.Get("Authorization")
	if headerText == "" {
		respondWithError(w, http.StatusUnauthorized, "authorization header missing")
		return
	}

	prefix := "Bearer "
	if !strings.HasPrefix(headerText, prefix) {
		respondWithError(w, http.StatusUnauthorized, "invalid token format")
		return
	}
	tokenStr := strings.TrimPrefix(headerText, prefix)

	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 	return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		// }
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || claims.ExpiresAt == nil || claims.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "invalid token claims")
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid user ID")
		return
	}

	userReq := userRequest{}
	err = json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprint(err))
		return
	}

	hashPassword := validatePassword(w, userReq.Password)
	if hashPassword == nil {
		return
	}

	user, err := cfg.db.UpdateUser(userID, userReq.Email, hashPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprint(err))
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
		respondWithError(w, http.StatusBadRequest, fmt.Sprint(err))
		return nil
	}
	return hashPassword
}
