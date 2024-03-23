package services

import (
	"backendGoAuth/internal/auth"
	"backendGoAuth/internal/database"
	"backendGoAuth/internal/goAuthException"
	"backendGoAuth/internal/models"
	"golang.org/x/crypto/bcrypt"
	"strings"
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
func (svc *AuthService) RegisterUser(req models.RegistrationRequest, ipAddress string) (models.AuthResponse, error) {
	// Check if username already exists
	exists, err := database.UserExistsByUsername(req.Username)
	if err != nil {
		return models.AuthResponse{}, err
	}
	if exists {
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.BadRequestCode, goAuthException.UsernameExistsMessage)
	}

	// Check if email already exists
	existsEmail, err := database.UserExistsByEmail(req.Email)
	if err != nil {
		return models.AuthResponse{}, err
	}
	if existsEmail {
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.BadRequestCode, goAuthException.EmailExistsMessage)
	}

	// Hash the password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.Teapot, goAuthException.HashingError)
	}

	// Create a new user
	user, err := svc.createUser(req.Username, hashedPassword, req.Email)
	if err != nil {
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.Teapot, goAuthException.UserCreationError)
	}

	// Generate a JWT token for the newly registered user
	token, err := auth.GenerateJWT(map[string]interface{}{"user_id": user.ID})
	if err != nil {
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.Teapot, goAuthException.TokenGenerationError)
	}

	// Insert session into the database
	err = database.InsertSession(user.ID, token, ipAddress)
	if err != nil {
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.Teapot, goAuthException.SessionInsertionError)
	}

	return models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// AuthenticateUser authenticates a user based on the provided credentials.
// AuthenticateUser authenticates a user based on the provided identifier (username or email) and password.
func (svc *AuthService) AuthenticateUser(identifier, password, ipAddress string) (models.AuthResponse, error) {
	var user *database.User
	var err error

	if strings.Contains(identifier, "@") {
		user, err = database.GetUserByEmail(identifier)
	} else {
		user, err = database.GetUserByUsername(identifier)
	}
	if err != nil {
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.NotFoundCode, goAuthException.UsernameCheckError)
	}

	// Compare hashed passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.UnauthorizedCode, "Authentication failed")
	}

	// Generate a JWT token with the user's ID
	token, err := auth.GenerateJWT(map[string]interface{}{"user_id": user.ID})
	if err != nil {
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.Teapot, goAuthException.TokenGenerationError)
	}

	// Insert session into the database
	err = database.InsertSession(user.ID, token, ipAddress)
	if err != nil {
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.InternalServerErrorCode, goAuthException.SessionInsertionError)
	}

	authResponse := models.AuthResponse{
		Token: token,
		User: models.UserData{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}

	return authResponse, nil
}

func (svc *AuthService) GetActiveSessions(userID int) ([]database.Session, error) {
	return database.GetActiveSessions(userID)
}

func (svc *AuthService) RevokeSessionToken(token string) error {
	return database.RevokeSession(token)
}

// createUser creates a new user in the database.
func (svc *AuthService) createUser(username, password, email string) (models.UserData, error) {
	userID, err := database.InsertUser(username, password, email)
	if err != nil {
		return models.UserData{}, goAuthException.NewCustomError(goAuthException.InternalServerErrorCode, goAuthException.UserCreationError)
	}

	return models.UserData{
		ID:       userID,
		Username: username,
		Email:    email,
	}, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", goAuthException.NewCustomError(goAuthException.InternalServerErrorCode, goAuthException.SessionInsertionError)
	}
	return string(hashedPassword), nil
}
