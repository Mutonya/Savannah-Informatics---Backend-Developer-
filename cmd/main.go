package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// Initialize logger
	log.Println("Starting Savannah backend service...")

	// Create HTTP server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server
	log.Printf("Server listening on port %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to Savannah Backend!")
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status": "healthy"}`)
	})

	return mux
}
