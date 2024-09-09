package services

import (
	"backendGoAuth/internal/entities"
	"backendGoAuth/internal/goAuthException"
	"backendGoAuth/internal/models"
	"backendGoAuth/internal/repositories"
	"backendGoAuth/internal/utils"
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

	// Generate a JWT token for the newly registered user
	token, err := utils.GenerateJWT(map[string]interface{}{"user_id": user.ID})
	if err != nil {
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.Teapot, goAuthException.TokenGenerationError)
	}

	// Insert session into the database
	sessionID, err := svc.SessionService.InsertSession(user.ID, token, ipAddress, browser, device)
	if err != nil {
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.Teapot, goAuthException.SessionInsertionError)
	}

	return models.AuthResponse{
		Token: token,
		User:  user,
		Session: models.SessionResponse{
			ID: sessionID,
		},
	}, nil
}

// AuthenticateUser authenticates a user based on the provided credentials.
// AuthenticateUser authenticates a user based on the provided identifier (username or email) and password.
func (svc *AuthService) AuthenticateUser(identifier, password, ipAddress, browser, device string) (models.AuthResponse, error) {
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

	// Generate a JWT token with the user's ID
	token, err := utils.GenerateJWT(map[string]interface{}{"user_id": user.ID})
	if err != nil {
		return models.AuthResponse{}, goAuthException.NewCustomError(goAuthException.Teapot, goAuthException.TokenGenerationError)
	}

	// Insert session into the database
	sessionID, err := svc.SessionService.InsertSession(user.ID, token, ipAddress, browser, device)
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
		Session: models.SessionResponse{
			ID: sessionID,
		},
	}

	return authResponse, nil
}

//func (svc *AuthService) GetActiveSessions(userID int) ([]database.Session, error) {
//	return database.GetActiveSessions(userID)
//}

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
