package app

import (
	"log"
	"net/http"
	"os"
	"project-management-service/internal/database"
	"github.com/go-chi/chi/v5"
)

func Run() {
	// Initialize the database
	database.InitDB()

	// Create a new chi router
	router := chi.NewRouter()

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}