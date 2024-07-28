package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func runServer() error {
	const filepathRoot = "dist"
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	apiCfg, err := startDB()
	if err != nil {
		return fmt.Errorf("error starting database: %s", err)
	}
	fileserverHandler := apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))),
	)

	mux.Handle("/app/*", fileserverHandler)

	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("GET /api/reset", apiCfg.resetMetricsHandler)

	mux.HandleFunc("POST /api/chirps", apiCfg.createChirpHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpid}", apiCfg.getChirpByIDHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpid}", apiCfg.deleteChirpHandler)

	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)
	mux.HandleFunc("PUT /api/users", apiCfg.updateUserHandler)

	mux.HandleFunc("POST /api/refresh", apiCfg.refreshAccessTokenHandler)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeRefreshTokenHandler)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.upgradeUserHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	err = server.ListenAndServe()
	log.Println("Starting server on port", port)
	if err != nil {
		return fmt.Errorf("error starting server: %s", err)
	}
	return nil
}
