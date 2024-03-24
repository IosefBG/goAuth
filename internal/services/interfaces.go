package services

import "backendGoAuth/internal/models"

type AuthServiceInterface interface {
	RegisterUser(req models.RegistrationRequest, ipAddress, browser, device string) (models.AuthResponse, error)
	AuthenticateUser(identifier, password, ipAddress, browser, device string) (models.AuthResponse, error)
	RevokeSession(sessionID int) error
}

type SessionServiceInterface interface {
	InsertSession(userID int, token, ipAddress, browser, device string) error
	GetActiveSessions(userID int) ([]models.SessionResponse, error)
	RevokeCurrentSessionToken(token string) error
}
