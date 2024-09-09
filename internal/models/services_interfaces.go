package models

type AuthServiceInterface interface {
	RegisterUser(req RegistrationRequest, ipAddress, browser, device string) (AuthResponse, error)
	AuthenticateUser(identifier, password, ipAddress, browser, device string) (AuthResponse, error)
	RevokeSession(sessionID int) error
}

type SessionServiceInterface interface {
	InsertSession(userID int, token, ipAddress, browser, device string) (int, error)
	GetActiveSessions(userID int) ([]SessionResponse, error)
	RevokeCurrentSessionToken(token string) error
}
