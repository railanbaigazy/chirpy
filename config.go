package main

import (
	"fmt"

	"github.com/railanbaigazy/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
}

func startDB() (apiConfig, error) {
	const filepathDB = "database.json"
	db, err := database.NewDB(filepathDB)
	if err != nil {
		return apiConfig{}, fmt.Errorf("failed to initialize database: %v", err)
	}
	return apiConfig{fileserverHits: 0, db: db}, nil
}
