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
	"time"
)

var db *sql.DB

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
func ConnectDB() error {
	dbConfig := getDBConfigFromEnv()

	// Open a connection to the database
	err := openDB(dbConfig)
	if err != nil {
		log.Printf("Error connecting to the database: %v\n", err)
		return err
	}

	// Run database migrations
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

// GetDB returns the global db variable
func GetDB() *sql.DB {
	return db
}

// openDB opens a connection to the PostgreSQL database.
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

	// Ping the database to ensure the connection is valid
	if err := db.Ping(); err != nil {
		log.Printf("Error pinging database: %v\n", err)
		return err
	}

	return nil
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

// InsertUser inserts a new user into the database.
func InsertUser(username, password, email string) (int, error) {
	// Execute the insert query to insert the user into the database
	_, err := db.Exec("INSERT INTO users (username, password, email) VALUES ($1, $2, $3)", username, password, email)
	if err != nil {
		log.Printf("Error inserting user into database: %v\n", err)
		return 0, err
	}

	// Assuming you have a way to retrieve the userID after insertion
	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		log.Printf("Error retrieving userID: %v\n", err)
		return 0, err
	}

	log.Println("User inserted successfully")
	return userID, nil
}

func UserExistsByUsername(username string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	if err != nil {
		log.Printf("Error checking if user exists by username: %v\n", err)
		return false, err
	}
	return count > 0, nil
}

// InsertSession inserts a new session into the database.
func InsertSession(userID int, token, ipAddress string) error {
	now := time.Now()
	location, err := getLocationFromIPAddress(ipAddress)
	if err != nil {
		return err
	}
	_, err = db.Exec(
		"INSERT INTO user_sessions (user_id, session_token, ip_address, location, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, token, ipAddress, location, now, now,
	)
	if err != nil {
		return err
	}
	return nil
}

func getLocationFromIPAddress(ipAddress string) (string, error) {
	// Implement logic to obtain location information from the IP address.
	// This might involve using a geoip library or calling an external service.
	// Here's a simplified example using a placeholder value:
	return "New York, USA", nil
}

// GetActiveSessions retrieves active sessions for a user from the database.
func GetActiveSessions(userID int) ([]Session, error) {
	rows, err := db.Query(
		"SELECT id, ip_address, created_at, updated_at, location FROM user_sessions WHERE user_id = $1 AND is_active = true",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing connection %v\n", err)
		}
	}(rows)

	var sessions []Session
	for rows.Next() {
		var session Session
		if err := rows.Scan(
			&session.ID, &session.IPAddress, &session.CreatedAt, &session.UpdatedAt, &session.Location,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sessions, nil
}

// RevokeSession revokes a session in the database.
func RevokeSession(sessionID string) error {
	_, err := db.Exec(
		"UPDATE user_sessions SET is_active = false WHERE session_token = $1",
		sessionID,
	)
	if err != nil {
		return err
	}
	return nil
}

func CheckSession(tokenString string) (bool, error) {
	var session Session
	err := db.QueryRow("SELECT is_active FROM user_sessions WHERE session_token = $1 AND is_active = true", tokenString).Scan(&session.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		log.Printf("Error retrieving session: %v\n", err)
		return false, err
	}
	return session.IsActive, nil
}

func UserExistsByEmail(email string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		log.Printf("Error checking if user exists by email: %v\n", err)
		return false, err
	}
	return count > 0, nil
}

func UpdateSessionUpdatedAt(userID int) error {
	// Get current time
	currentTime := time.Now()

	// Perform update query to update the updated_at column
	_, err := db.Exec(
		"UPDATE user_sessions SET updated_at = $1 WHERE user_id = $2",
		currentTime, userID,
	)
	return err
}

func GetUserByEmail(identifier string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, username, password, email FROM users WHERE email = $1", identifier).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("Error retrieving user by username: %v\n", err)
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(username string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT username FROM users WHERE username = $1", username).Scan(&user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("Error retrieving user by username: %v\n", err)
		return nil, err
	}
	return &user, nil
}
