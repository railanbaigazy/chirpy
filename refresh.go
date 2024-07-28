package main

import "net/http"

func (cfg *apiConfig) refreshAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenStr, err := getTokenString(w, r)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	refreshResp, err := cfg.db.RefreshAccessToken(tokenStr, []byte(cfg.jwtSecret))
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, refreshResp)
}
