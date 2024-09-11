// Package services AuthService
package services

import (
	"backendGoAuth/internal/entities"
	"backendGoAuth/internal/goAuthException"
	"backendGoAuth/internal/models"
	"backendGoAuth/internal/repositories"
	"backendGoAuth/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

// AuthService provides authentication-related services.
type AuthService struct {
	UserRepo       *repositories.UserRepository
	SessionService *SessionService // Corrected reference
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(userRepo *repositories.UserRepository, sessionService *SessionService) *AuthService {
	return &AuthService{
		UserRepo:       userRepo, // Initialize UserRepo here
		SessionService: sessionService,
	}
}

// RegisterUser registers a new user with the provided details.
func (svc *AuthService) RegisterUser(req models.RegistrationRequest, ipAddress, browser, device string) (models.AuthResponse, error) {
	// Check if username already exists
	exists, err := svc.UserRepo.UserExistsByUsername(req.Username)
	if err != nil {
		return models.AuthResponse{}, err
	}
	if exists {
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.BadRequestCode, goAuthException.UsernameExistsMessage)
	}

	// Check if email already exists
	existsEmail, err := svc.UserRepo.UserExistsByEmail(req.Email)
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

	// Insert session into the database
	session, err := svc.SessionService.InsertSession(user.ID, ipAddress, browser, device)
	if err != nil {
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.Teapot, goAuthException.SessionInsertionError)
	}

	// Generate JWT with session ID included in the claims
	//todo also here should set cookie
	_, err = utils.GenerateJWT(map[string]interface{}{
		"user_id":    user.ID,
		"session_id": session.ID, // Include session ID in the token claims
	}, "access", session.ID)

	return models.AuthResponse{
		User: user,
	}, nil
}

// AuthenticateUser authenticates a user based on the provided credentials.
func (svc *AuthService) AuthenticateUser(identifier, password, ipAddress, browser, device string, c *gin.Context) (models.AuthResponse, error) {
	var user *entities.User
	var err error

	if strings.Contains(identifier, "@") {
		user, err = svc.UserRepo.GetUserByEmail(identifier)
		if err != nil {
			return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.NotFoundCode, goAuthException.UsernameCheckError)
		}
	} else {
		user, err = svc.UserRepo.GetUserByUsername(identifier)
		if err != nil {
			return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.NotFoundCode, goAuthException.EmailCheckError)
		}
	}

	if user == nil {
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.NotFoundCode, "User doesn't exist")
	}

	// Compare hashed passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password comparison failed for user: %s\n", identifier)
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.UnauthorizedCode, "Invalid credentials")
	}

	session, err := svc.SessionService.InsertSession(user.ID, ipAddress, browser, device)
	if err != nil {
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.InternalServerErrorCode, goAuthException.SessionInsertionError)
	}

	// Generate short-lived JWT token with the user's ID
	accessToken, err := utils.GenerateJWT(map[string]interface{}{
		"user_id":    user.ID,
		"session_id": session.ID, // Include session ID in the token claims
	}, "access", session.ID)
	if err != nil {
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.Teapot, goAuthException.TokenGenerationError)
	}

	// Store session in the database with the refresh token, if applicable
	//refreshToken, err := utils.GenerateJWT(map[string]interface{}{
	//	"user_id":    user.ID,
	//	"session_id": session.ID, // Include session ID in the token claims
	//}, "refresh", session.ID)
	//if err != nil {
	//	return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.Teapot, goAuthException.TokenGenerationError)
	//}

	authResponse := models.AuthResponse{
		User: models.UserData{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}

	utils.SetJWTTokenCookies(c, accessToken)

	return authResponse, nil
}

// createUser creates a new user in the database.
func (svc *AuthService) createUser(username, password, email string) (models.UserData, error) {
	userID, err := svc.UserRepo.InsertUser(username, password, email)
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

func (svc *AuthService) RevokeSession(sessionID string) error {
	// Call the relevant method to revoke the session in your service layer
	err := svc.SessionService.RevokeSession(sessionID)
	if err != nil {
		// Handle any errors that occur during session revocation
		return err
	}
	return nil
}
