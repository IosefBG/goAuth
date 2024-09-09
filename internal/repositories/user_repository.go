package repositories

import (
	"backendGoAuth/internal/entities"
	"database/sql"
	"errors"
	"log"
)

// UserRepository is the concrete struct for interacting with user-related data.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

// GetAllUsers retrieves all active users from the database.
func (r *UserRepository) GetAllUsers() ([]entities.User, error) {
	query := "SELECT id, username, email, is_blocked, login_attempts, last_login, created_at, updated_at, is_active FROM users"
	rows, err := r.db.Query(query)
	if err != nil {
		log.Println("Error querying all users:", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Error closing rows:", err)
		}
	}(rows)

	var users []entities.User
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.IsBlocked, &user.LoginAttempts, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt, &user.IsActive); err != nil {
			log.Println("Error scanning user:", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// EditUser updates a user's details in the database.
func (r *UserRepository) EditUser(user entities.User) error {
	query := "UPDATE users SET username = $1, email = $2, is_blocked = $3, login_attempts = $4, updated_at = $5 WHERE id = $6"
	_, err := r.db.Exec(query, user.Username, user.Email, user.IsBlocked, user.LoginAttempts, user.UpdatedAt, user.ID)
	if err != nil {
		log.Println("Error updating user:", err)
		return err
	}
	return nil
}

// DeleteUser performs a soft delete on a user by setting is_active to FALSE.
func (r *UserRepository) DeleteUser(userID int) error {
	query := "UPDATE users SET is_active = FALSE WHERE id = $1"
	_, err := r.db.Exec(query, userID)
	if err != nil {
		log.Println("Error deleting (soft delete) user:", err)
		return err
	}
	return nil
}

// InsertUser adds a new user to the database and returns the new user's ID.
func (r *UserRepository) InsertUser(username, password, email string) (int, error) {
	_, err := r.db.Exec("INSERT INTO users (username, password, email) VALUES ($1, $2, $3)", username, password, email)
	if err != nil {
		log.Printf("Error inserting user into database: %v\n", err)
		return 0, err
	}

	var userID int
	err = r.db.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		log.Printf("Error retrieving userID: %v\n", err)
		return 0, err
	}

	log.Println("User inserted successfully")
	return userID, nil
}

// UserExistsByUsername checks if a user exists by their username.
func (r *UserRepository) UserExistsByUsername(username string) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	if err != nil {
		log.Printf("Error checking if user exists by username: %v\n", err)
		return false, err
	}
	return count > 0, nil
}

// UserExistsByEmail checks if a user exists by their email.
func (r *UserRepository) UserExistsByEmail(email string) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		log.Printf("Error checking if user exists by email: %v\n", err)
		return false, err
	}
	return count > 0, nil
}

// GetUserByEmail retrieves a user by their email.
func (r *UserRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.QueryRow("SELECT id, username, password, email FROM users WHERE email = $1", email).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("Error retrieving user by email: %v\n", err)
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by their username.
func (r *UserRepository) GetUserByUsername(username string) (*entities.User, error) {
	var user entities.User
	err := r.db.QueryRow("SELECT id, username, password, email FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("Error retrieving user by username: %v\n", err)
		return nil, err
	}
	return &user, nil
}

// RevokeUser revokes (soft deletes) a user by setting is_active to FALSE.
func (r *UserRepository) RevokeUser(userID int) error {
	_, err := r.db.Exec(
		"UPDATE users SET is_active = false WHERE id = $1",
		userID,
	)
	if err != nil {
		log.Printf("Error revoking user: %v\n", err)
	}
	return err
}
