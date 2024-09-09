package repositories

import (
	"backendGoAuth/internal/entities"
	"database/sql"
	"errors"
	"log"
	"time"
)

type SessionRepository struct {
	DB *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

func (r *SessionRepository) InsertSession(session entities.Session) (int, error) {
	var sessionID int
	err := r.DB.QueryRow(
		"INSERT INTO user_sessions (user_id, session_token, ip_address, location, created_at, updated_at, device_connected, browser_used) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		session.UserID, session.Token, session.IPAddress, session.Location, session.CreatedAt, session.UpdatedAt, session.DeviceConnected, session.BrowserUsed,
	).Scan(&sessionID)
	if err != nil {
		return 0, err
	}
	return sessionID, nil
}

func (r *SessionRepository) GetActiveSessions(userID int) (*sql.Rows, error) {
	rows, err := r.DB.Query(
		"SELECT id, ip_address, created_at, updated_at, location, device_connected, browser_used FROM user_sessions WHERE user_id = $1 AND is_active = true",
		userID,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *SessionRepository) RevokeCurrentSession(sessionID string) error {
	_, err := r.DB.Exec(
		"UPDATE user_sessions SET is_active = false WHERE session_token = $1",
		sessionID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *SessionRepository) CheckSession(tokenString string) (bool, error) {
	var session entities.Session
	err := r.DB.QueryRow("SELECT is_active FROM user_sessions WHERE session_token = $1 AND is_active = true", tokenString).Scan(&session.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		log.Printf("Error retrieving session: %v\n", err)
		return false, err
	}
	return session.IsActive, nil
}

func (r *SessionRepository) UpdateSessionUpdatedAt(userID int) error {
	currentTime := time.Now()
	_, err := r.DB.Exec(
		"UPDATE user_sessions SET updated_at = $1 WHERE user_id = $2",
		currentTime, userID,
	)
	return err
}

// fixme check if this actually works, instead of id should be session token?
func (r *SessionRepository) RevokeSession(sessionID string) error {
	_, err := r.DB.Exec(
		"UPDATE user_sessions SET is_active = false WHERE id = $1",
		sessionID,
	)
	if err != nil {
		return err
	}
	return nil
}
