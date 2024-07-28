package main

import (
	"errors"
	"net/http"
	"strings"
)

func getTokenString(w http.ResponseWriter, r *http.Request) (string, error) {
	headerText := r.Header.Get("Authorization")
	if headerText == "" {
		return "", errors.New("authorization header missing")
	}

	prefix := "Bearer "
	if !strings.HasPrefix(headerText, prefix) {
		respondWithError(w, http.StatusUnauthorized, "invalid token format")
		return "", errors.New("invalid token format")
	}
	tokenStr := strings.TrimPrefix(headerText, prefix)
	return tokenStr, nil
}
