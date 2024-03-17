package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// DBConfig holds the configuration for the database connection.
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// getDBConfigFromEnv retrieves the database configuration from environment variables.
func getDBConfigFromEnv() DBConfig {
	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}

// ConnectDB connects to the PostgreSQL database and applies migrations.
func ConnectDB() (*sql.DB, error) {
	dbConfig := getDBConfigFromEnv()

	// Open a connection to the database
	db, err := openDB(dbConfig)
	if err != nil {
		log.Printf("Error connecting to the database: %v\n", err)
		return nil, err
	}

	// Run database migrations
	if err := runMigrations(dbConfig); err != nil {
		log.Printf("Error applying database migrations: %v\n", err)
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database connection: %v\n", closeErr)
		}
		return nil, err
	}

	log.Println("Database connection and migrations applied successfully")
	return db, nil
}

// openDB opens a connection to the PostgreSQL database.
func openDB(cfg DBConfig) (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Printf("Error opening database connection: %v\n", err)
		return nil, err
	}
	return db, nil
}

// runMigrations runs database migrations using the "migrate" library.
func runMigrations(cfg DBConfig) error {
	// Specify the migration directory
	migrationDir := "migrations"

	// Initialize the migration tool
	m, err := migrate.New(
		"file://"+migrationDir,
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name),
	)
	if err != nil {
		log.Printf("Error initializing migration tool: %v\n", err)
		return err
	}

	// Apply all pending migrations
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Printf("Error applying migrations: %v\n", err)
			return err
		}
	}

	log.Println("Migrations applied successfully")
	return nil
}
