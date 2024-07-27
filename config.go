package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/railanbaigazy/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
}

func startDB() (apiConfig, error) {
	const filepathDB = "database.json"

	isDebug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *isDebug {
		fmt.Println("Debug mode enabled. Deleting the database...")
		if err := os.Remove(filepathDB); err != nil {
			log.Print("Database is already deleted")
		} else {
			log.Print("Database is successfully deleted")
		}
	}
	db, err := database.NewDB(filepathDB)
	if err != nil {
		return apiConfig{}, fmt.Errorf("failed to initialize database: %v", err)
	}
	log.Print("Config is created")
	return apiConfig{fileserverHits: 0, db: db}, nil
}
