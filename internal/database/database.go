package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"project-management-service/internal/config"

	_ "github.com/lib/pq"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var DB *sql.DB

func InitDB() {
	var err error

	// Load configuration
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Print DB_SOURCE for debugging
	fmt.Println("DB_SOURCE:", config.DBSource)

	// Open a connection to the database using DB_SOURCE directly
	DB, err = sql.Open("postgres", config.DBSource)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Check if the connection to the database is working
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	// Run database migrations
	if err := runMigrations(config.DBSource); err != nil {
		log.Fatalf("Could not run database migrations: %v", err)
	}

	log.Println("Connected to the database successfully!")
}

func runMigrations(dataSourceName string) error {
	// Get the absolute path for the migrations directory
	dir, err := filepath.Abs("./db/migrations")
	if err != nil {
		return fmt.Errorf("could not get absolute path for migrations: %w", err)
	}

	// Check if the migrations directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("Migrations directory '%s' does not exist, skipping migrations", dir)
		return nil // No migrations to run
	}

	// Print the contents of the migrations directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("could not read migrations directory: %w", err)
	}

	log.Println("Migration files:")
	for _, file := range files {
		log.Println(file.Name())
	}

	// Perform database migrations using golang-migrate
	if !strings.Contains(dataSourceName, "://") {
		return errors.New("undefined data source name " + dataSourceName)
	}

	m, err := migrate.New(fmt.Sprintf("file://%s", dir), dataSourceName)
	if err != nil {
		return fmt.Errorf("error creating migration instance: %w", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error running migrations: %w", err)
	}

	log.Println("Database migrations ran successfully!")
	return nil
}
