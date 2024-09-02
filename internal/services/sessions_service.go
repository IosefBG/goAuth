package services

import (
	"backendGoAuth/internal/database"
	"backendGoAuth/internal/models"
	"fmt"
	"log"
	"time"
)

type SessionService struct {
	// Add any dependencies here
}

// NewSessionService creates a new instance of SessionService.
func NewSessionService() *SessionService {
	return &SessionService{}
}

// getLocationFromIPAddress gets the location from the IP address
func getLocationFromIPAddress(ipAddress string) (string, error) {
	return "Location", nil // Placeholder
}

func (s *SessionService) InsertSession(userID int, token, ipAddress, browser, device string) (int, error) {
	now := time.Now()
	location, err := getLocationFromIPAddress(ipAddress)
	if err != nil {
		return 0, err
	}

	session := database.Session{
		UserID:          userID,
		Token:           token,
		IPAddress:       ipAddress,
		Location:        location,
		CreatedAt:       now,
		UpdatedAt:       now,
		DeviceConnected: device,
		BrowserUsed:     browser,
	}

	sessionID, err := database.InsertSession(session)
	if err != nil {
		return 0, err
	}
	return sessionID, nil
}

func (s *SessionService) GetActiveSessions(userID int) ([]models.SessionResponse, error) {
	// Retrieve active sessions from the database
	rows, err := database.GetActiveSessionsQuery(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve active sessions: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error closing rows: %v\n", closeErr)
		}
	}()

	var sessions []models.SessionResponse
	for rows.Next() {
		var session database.Session
		if err := rows.Scan(
			&session.ID, &session.IPAddress, &session.CreatedAt, &session.UpdatedAt, &session.Location, &session.DeviceConnected, &session.BrowserUsed,
		); err != nil {
			return nil, fmt.Errorf("failed to scan session row: %w", err)
		}

		// Map database.Session to models.SessionResponse
		sessionResponse := models.SessionResponse{
			ID:              session.ID,
			UserID:          session.UserID,
			Token:           session.Token,
			IPAddress:       session.IPAddress,
			IsActive:        session.IsActive,
			CreatedAt:       session.CreatedAt,
			UpdatedAt:       session.UpdatedAt,
			Location:        session.Location,
			DeviceConnected: session.DeviceConnected,
			BrowserUsed:     session.BrowserUsed,
		}

		sessions = append(sessions, sessionResponse)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over session rows: %w", err)
	}
	return sessions, nil
}

func (s *SessionService) RevokeCurrentSessionToken(token string) error {
	return database.RevokeCurrentSession(token)
}

func (s *SessionService) RevokeSession(sessionID int) error {
	// Your logic to revoke the session using the sessionID
	// This might involve updating the database to mark the session as inactive or deleting it
	// For example, assuming you have a function in your database package to update session status:
	err := database.RevokeSession(sessionID) // Marking session as inactive
	if err != nil {
		// Handle any errors
		return err
	}
	return nil
}
