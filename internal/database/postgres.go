package database

import (
	"database/sql"
	"errors"
	"fmt"
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

// ConnectDB connects to the PostgreSQL database and applies migrations.
func ConnectDB() (*sql.DB, error) {
	dbConfig := getDBConfigFromEnv()

	// Open a connection to the database
	db, err := openDB(dbConfig)
	if err != nil {
		return nil, err
	}

	// Run database migrations
	if err := runMigrations(dbConfig); err != nil {
		// Close the database connection in case of an error during migrations
		if closeErr := db.Close(); closeErr != nil {
			return nil, closeErr
		}
		return nil, err
	}

	return db, nil
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

// openDB opens a connection to the PostgreSQL database.
func openDB(cfg DBConfig) (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)
	return sql.Open("postgres", connectionString)
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
		return err
	}

	// Apply all pending migrations
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	return nil
}
