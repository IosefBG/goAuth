package services

import (
	"backendGoAuth/internal/auth"
	"backendGoAuth/internal/database"
	"backendGoAuth/internal/models"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// AuthService provides authentication-related services.
type AuthService struct {
	// Add any dependencies here
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService() *AuthService {
	return &AuthService{}
}

// RegisterUser registers a new user with the provided details.
func (svc *AuthService) RegisterUser(req models.RegistrationRequest, ipAddress string) (string, error) {
	// Hash the password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return "", err
	}

	// Create a new user
	userID, err := svc.createUser(req.Username, hashedPassword, req.Email)
	if err != nil {
		return "", err
	}

	// Generate a JWT token for the newly registered user
	token, err := auth.GenerateJWT(map[string]interface{}{"user_id": userID})
	if err != nil {
		return "", err
	}

	// Insert session into the database
	err = database.InsertSession(userID, token, ipAddress)
	if err != nil {
		return "", err
	}

	return token, nil
}

// AuthenticateUser authenticates a user based on the provided credentials.
func (svc *AuthService) AuthenticateUser(username, password, ipAddress string) (int, string, error) {
	user, err := database.GetUserByUsername(username)
	if err != nil {
		return 0, "", err
	}
	log.Println("user", user.ID)

	// Compare hashed passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, "", errors.New("authentication failed")
	}

	// Generate a JWT token with the user's ID
	token, err := auth.GenerateJWT(map[string]interface{}{"user_id": user.ID})
	if err != nil {
		return 0, "", err
	}

	// Insert session into the database
	err = database.InsertSession(user.ID, token, ipAddress)
	if err != nil {
		return 0, "", err
	}

	return user.ID, token, nil
}

func (svc *AuthService) GetActiveSessions(userID int) ([]database.Session, error) {
	return database.GetActiveSessions(userID)
}

func (svc *AuthService) RevokeSessionToken(token string) error {
	return database.RevokeSession(token)
}

// createUser creates a new user in the database.
func (svc *AuthService) createUser(username, password, email string) (int, error) {
	// Call InsertUser function from the database package
	userID, err := database.InsertUser(username, password, email)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
