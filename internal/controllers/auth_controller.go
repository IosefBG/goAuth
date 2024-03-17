// In your controllers/auth_controller.go

package controllers

import (
	"backendGoAuth/internal/auth"
	"backendGoAuth/internal/database"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// AuthController handles authentication-related requests.
type AuthController struct {
	// Add any dependencies here
}

type RegistrationRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// Register handles the registration request.
func (controller *AuthController) Register(c *gin.Context) {
	// Parse the request body to extract the registration details
	var req RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Hash the password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Implement logic to create a new user with the provided details
	userID, err := createUser(req.Username, hashedPassword, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate a JWT token for the newly registered user
	token, err := auth.GenerateJWT(map[string]interface{}{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Insert session into the database
	ipAddress := c.ClientIP() // Get client IP address
	err = database.InsertSession(userID, token, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func createUser(username, password, email string) (int, error) {
	// Call InsertUser function from the database package
	userID, err := database.InsertUser(username, password, email)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// NewAuthController creates a new instance of AuthController.
func NewAuthController() *AuthController {
	return &AuthController{}
}

// Login handles the login request.
func (controller *AuthController) Login(c *gin.Context) {
	// Parse the request body to extract the user's credentials
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Authenticate the user based on the provided credentials
	userID, err := authenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate a JWT token with the user's ID
	token, err := auth.GenerateJWT(map[string]interface{}{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Insert session into the database
	ipAddress := c.ClientIP() // Get client IP address
	err = database.InsertSession(userID, token, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// RevokeSession revokes a session for the current user.
func (controller *AuthController) RevokeSession(c *gin.Context) {
	// Extract the session token from the request
	token, err := auth.ExtractToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid token"})
		return
	}

	// Revoke the session token
	err = revokeSessionToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke session token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session revoked"})
}

// AuthenticateUser authenticates the user based on the provided credentials.
func authenticateUser(username, password string) (int, error) {
	// Retrieve user from database by username
	user, err := database.GetUserByUsername(username)
	if err != nil {
		return 0, err
	}

	// Compare hashed passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, errors.New("authentication failed")
	}

	return user.ID, nil
}

// revokeSessionToken revokes a session token from the database.
func revokeSessionToken(token string) error {
	// Implement logic to revoke the session token from the database
	// Example: db.Exec("DELETE FROM sessions WHERE token = ?", token)
	return nil
}

func (controller *AuthController) SecureEndpoint(c *gin.Context) {
	// Access user ID from the context
	userID, err := auth.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Your secure endpoint logic here
	c.JSON(http.StatusOK, gin.H{"message": "Secure Endpoint", "user_id": userID})
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
