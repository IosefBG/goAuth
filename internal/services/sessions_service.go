package services

import (
	"backendGoAuth/internal/entities"
	"backendGoAuth/internal/models"
	"backendGoAuth/internal/repositories"
	"backendGoAuth/internal/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type SessionService struct {
	SessionRepo *repositories.SessionRepository
}

// NewSessionService creates a new instance of SessionService.
func NewSessionService(sessionRepo *repositories.SessionRepository) *SessionService {
	return &SessionService{
		SessionRepo: sessionRepo,
	}
}

// getLocationFromIPAddress gets the location from the IP address
func getLocationFromIPAddress(ipAddress string) (string, error) {
	// todo add some way of taking user location maybe?
	return "Location", nil // Placeholder
}

func (s *SessionService) InsertSession(userID int, ipAddress, browser, device string) (*entities.Session, error) {
	now := time.Now()
	location, err := getLocationFromIPAddress(ipAddress)
	if err != nil {
		return nil, err
	}

	session := entities.Session{
		UserID:          userID,
		IPAddress:       ipAddress,
		Location:        location,
		CreatedAt:       now,
		UpdatedAt:       now,
		DeviceConnected: device,
		BrowserUsed:     browser,
		IsActive:        true,
	}

	// Insert the session and get the session ID
	sessionID, err := s.SessionRepo.InsertSession(session)
	if err != nil {
		return nil, err
	}

	// Retrieve the full session with the generated ID
	session.ID = sessionID
	return &session, nil
}

func (s *SessionService) GetActiveSessions(userID int) ([]models.SessionResponse, error) {
	// Retrieve active sessions from the repository
	rows, err := s.SessionRepo.GetActiveSessions(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve active sessions: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error closing rows: %v\n", closeErr)
		}
	}()

	var sessionResponses []models.SessionResponse

	for rows.Next() {
		var session entities.Session
		if err := rows.Scan(
			&session.ID,
			&session.IPAddress,
			&session.CreatedAt,
			&session.UpdatedAt,
			&session.Location,
			&session.DeviceConnected,
			&session.BrowserUsed,
			&session.IsActive, // Ensure this is included as it is in the SELECT statement
		); err != nil {
			return nil, fmt.Errorf("failed to scan session row: %w", err)
		}

		sessionResponse := models.SessionResponse{
			ID:              session.ID,
			UserID:          session.UserID,
			IPAddress:       session.IPAddress,
			IsActive:        session.IsActive,
			CreatedAt:       session.CreatedAt,
			UpdatedAt:       session.UpdatedAt,
			Location:        session.Location,
			DeviceConnected: session.DeviceConnected,
			BrowserUsed:     session.BrowserUsed,
		}

		sessionResponses = append(sessionResponses, sessionResponse)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over session rows: %w", err)
	}

	return sessionResponses, nil
}

func (s *SessionService) RevokeCurrentSessionToken(tokenString string) error {
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return err
	}

	return s.SessionRepo.RevokeCurrentSession(int(claims["session_id"].(float64)))
}

func (s *SessionService) RevokeSession(sessionID string) error {
	// Mark the session as inactive using the repository
	err := s.SessionRepo.RevokeSession(sessionID)
	if err != nil {
		// Handle any errors
		return err
	}
	return nil
}

func (s *SessionService) UpdateSessionUpdatedAt(userID int) error {
	err := s.SessionRepo.UpdateSessionUpdatedAt(userID)
	if err != nil {
		// Handle any errors
		return err
	}
	return nil
}

func (svc *SessionService) GetUserIDFromTokenOrSource(c *gin.Context) (int, error) {
	// Extract the JWT token from the request headers
	tokenString, err := utils.ExtractToken(c)
	if err != nil {
		return 0, err
	}

	// Validate the token and retrieve the claims
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}

	// Check the session using the CheckSession method
	sessionValid, err := svc.SessionRepo.CheckSession(int(claims["session_id"].(float64)))
	if err != nil {
		return 0, err
	}
	if !sessionValid {
		return 0, errors.New("session is invalid or expired")
	}

	// Extract the user ID from the claims
	userIDFloat64, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user ID not found in claims or not of the expected type")
	}

	return int(userIDFloat64), nil
}

func (svc *SessionService) GetSessionByID(sessionID int) (*entities.Session, error) {
	// Fetch the session by ID using the repository
	session, err := svc.SessionRepo.GetSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("error fetching session: %v", err)
	}

	if session == nil {
		// Return an error if the session is not found
		return nil, errors.New("session not found")
	}

	// Return the session if found
	return session, nil
}
