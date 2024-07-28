package main

import (
	"errors"
	"net/http"
	"strings"
)

func getTokenString(r *http.Request) (string, error) {
	headerText := r.Header.Get("Authorization")
	if headerText == "" {
		return "", errors.New("authorization header missing")
	}

	prefix := "Bearer "
	if !strings.HasPrefix(headerText, prefix) {
		return "", errors.New("invalid token format")
	}
	tokenStr := strings.TrimPrefix(headerText, prefix)
	return tokenStr, nil
}

func getApiKey(r *http.Request) (string, error) {
	headerText := r.Header.Get("Authorization")
	if headerText == "" {
		return "", errors.New("authorization header missing")
	}

	prefix := "ApiKey "
	if !strings.HasPrefix(headerText, prefix) {
		return "", errors.New("invalid apikey format")
	}
	apiKey := strings.TrimPrefix(headerText, prefix)
	return apiKey, nil
}
