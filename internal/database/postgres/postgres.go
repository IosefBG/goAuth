// internal/database/postgres/postgres.go

package postgres

import (
	"backendGoAuth/internal/config"
	"database/sql"
	"fmt"
)

// ConnectDB connects to the PostgreSQL database.
func ConnectDB(cfg config.AppConfig) (*sql.DB, error) {
	// Build the database connection string
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	// Open a connection to the database
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
