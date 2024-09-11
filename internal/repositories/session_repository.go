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
	if db == nil {
		log.Fatal("Database connection is nil in NewSessionRepository")
	}
	return &SessionRepository{DB: db}
}

func (r *SessionRepository) InsertSession(session entities.Session) (int, error) {
	var sessionID int
	err := r.DB.QueryRow(
		"INSERT INTO user_sessions (user_id, ip_address, location, created_at, updated_at, device_connected, browser_used) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		session.UserID, session.IPAddress, session.Location, session.CreatedAt, session.UpdatedAt, session.DeviceConnected, session.BrowserUsed,
	).Scan(&sessionID)
	if err != nil {
		return 0, err
	}
	return sessionID, nil
}

func (r *SessionRepository) GetActiveSessions(userID int) (*sql.Rows, error) {
	rows, err := r.DB.Query(
		"SELECT id, ip_address, created_at, updated_at, location, device_connected, browser_used, is_active FROM user_sessions WHERE user_id = $1 AND is_active = true",
		userID,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *SessionRepository) RevokeCurrentSession(sessionID int) error {
	_, err := r.DB.Exec(
		"UPDATE user_sessions SET is_active = false WHERE id = $1",
		sessionID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *SessionRepository) CheckSession(sessionId int) (bool, error) {
	var session entities.Session
	err := r.DB.QueryRow(`
    SELECT id, user_id, ip_address, is_active, created_at, updated_at, location, device_connected, browser_used
    FROM user_sessions WHERE id = $1
`, sessionId).Scan(
		&session.ID,
		&session.UserID,
		&session.IPAddress,
		&session.IsActive,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.Location,
		&session.DeviceConnected,
		&session.BrowserUsed,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		log.Printf("Error retrieving session: %v\n", err)
		return false, err
	}
	return session.IsActive, nil
}

func (r *SessionRepository) GetSessionByID(sessionID int) (*entities.Session, error) {
	if r.DB == nil {
		return nil, errors.New("database connection is nil")
	}

	var session entities.Session
	// Query the session from the database
	err := r.DB.QueryRow("SELECT id, is_active FROM user_sessions WHERE id = $1", sessionID).Scan(&session.ID, &session.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return nil if no session is found
			return nil, nil
		}
		// Log and return the error if there is any other issue
		log.Printf("Error retrieving session with ID %d: %v", sessionID, err)
		return nil, err
	}
	// Return the session object
	return &session, nil
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
