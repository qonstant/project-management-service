package app

import (
	"log"
	"net/http"
	"os"
	"project-management-service/internal/config"
	"project-management-service/internal/database"
	"project-management-service/internal/handlers"
)

func Run() {
	// Load configuration
	configs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize the database
	database.InitDB()

	// Set up dependencies
	deps := handler.Dependencies{
		DB:      database.DB,
		Configs: configs,
	}

	// Initialize the handler
	h, err := handler.New(deps, handler.WithHTTPHandler())
	if err != nil {
		log.Fatalf("Error initializing handlers: %v", err)
	}

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, h.HTTP); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
