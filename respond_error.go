package main

import (
	"encoding/json"
	"net/http"
)

type ErrorResp struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResp{Error: msg})
}
