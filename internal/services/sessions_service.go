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

func (s *SessionService) InsertSession(userID int, token, ipAddress, browser, device string) (int, error) {
	now := time.Now()
	location, err := getLocationFromIPAddress(ipAddress)
	if err != nil {
		return 0, err
	}

	session := entities.Session{
		UserID:          userID,
		Token:           token,
		IPAddress:       ipAddress,
		Location:        location,
		CreatedAt:       now,
		UpdatedAt:       now,
		DeviceConnected: device,
		BrowserUsed:     browser,
	}

	sessionID, err := s.SessionRepo.InsertSession(session)
	if err != nil {
		return 0, err
	}
	return sessionID, nil
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
			&session.Token,
			&session.IPAddress,
			&session.CreatedAt,
			&session.UpdatedAt,
			&session.Location,
			&session.DeviceConnected,
			&session.BrowserUsed,
			&session.IsActive, // Ensure you have this field if you want to include it
		); err != nil {
			return nil, fmt.Errorf("failed to scan session row: %w", err)
		}

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

		sessionResponses = append(sessionResponses, sessionResponse)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over session rows: %w", err)
	}

	return sessionResponses, nil
}

func (s *SessionService) RevokeCurrentSessionToken(token string) error {
	return s.SessionRepo.RevokeCurrentSession(token)
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

func (svc *SessionService) CheckSession(tokenString string) (bool, error) {
	// Call the CheckSession method from SessionRepository
	return svc.SessionRepo.CheckSession(tokenString)
}

func (svc *SessionService) GetUserIDFromTokenOrSource(c *gin.Context) (int, error) {
	// Extract the JWT token from the request headers
	tokenString, err := utils.ExtractToken(c)
	if err != nil {
		return 0, err
	}

	// Check the session using the CheckSession method
	sessionValid, err := svc.CheckSession(tokenString)
	if err != nil {
		return 0, err
	}
	if !sessionValid {
		return 0, errors.New("session is invalid or expired")
	}

	// Validate the token and retrieve the claims
	claims, err := utils.ValidateToken(tokenString, []byte(utils.JwtSecret))
	if err != nil {
		return 0, err
	}

	// Extract the user ID from the claims
	userIDFloat64, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user ID not found in claims or not of the expected type")
	}

	return int(userIDFloat64), nil
}
