package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

var db *sql.DB

func getDBConfigFromEnv() DBConfig {
	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}

func ConnectDB() error {
	dbConfig := getDBConfigFromEnv()

	err := openDB(dbConfig)
	if err != nil {
		log.Printf("Error connecting to the database: %v\n", err)
		return err
	}

	if err := runMigrations(dbConfig); err != nil {
		log.Printf("Error applying database migrations: %v\n", err)
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database connection: %v\n", closeErr)
		}
		return err
	}

	log.Println("Database connection and migrations applied successfully")
	return nil
}

func GetDB() *sql.DB {
	if db == nil {
		log.Fatalf("Database connection is not initialized")
	}
	return db
}

func openDB(cfg DBConfig) error {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)
	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Printf("Error opening database connection: %v\n", err)
		return err
	}

	if err := db.Ping(); err != nil {
		log.Printf("Error pinging database: %v\n", err)
		return err
	}

	return nil
}

func runMigrations(cfg DBConfig) error {
	migrationDir := "migrations"

	m, err := migrate.New(
		"file://"+migrationDir,
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name),
	)
	if err != nil {
		log.Printf("Error initializing migration tool: %v\n", err)
		return err
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Printf("Error applying migrations: %v\n", err)
			return err
		}
	}

	log.Println("Migrations applied successfully")
	return nil
}
