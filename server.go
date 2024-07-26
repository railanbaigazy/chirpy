package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
}

func runServer() error {
	const filepathRoot = "dist"
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	apiCfg := apiConfig{fileserverHits: 0}
	fileserverHandler := apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))),
	)

	mux.Handle("/app/*", fileserverHandler)
	mux.HandleFunc("GET /healthz", readinessHandler)
	mux.HandleFunc("GET /metrics", apiCfg.metricsHandler)
	mux.HandleFunc("GET /reset", apiCfg.resetMetricsHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	err := server.ListenAndServe()
	log.Println("Starting server on port", port)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
	return nil
}
