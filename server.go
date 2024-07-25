package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func runServer() error {
	godotenv.Load(".env")
	port := os.Getenv("PORT")

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server: %v", err)
	}
	return nil
}
