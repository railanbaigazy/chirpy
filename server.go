package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func runServer() error {
	godotenv.Load(".env")
	port := os.Getenv("PORT")

	mux := http.NewServeMux()
	mux.Handle("/app/*", http.StripPrefix("/app", http.FileServer(http.Dir("app"))))
	mux.HandleFunc("/healthz", readinessHandler)

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

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Println("error writing response:", err)
	}
}
