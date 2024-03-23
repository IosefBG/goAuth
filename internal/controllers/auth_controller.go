// In your controllers/auth_controller.go

package controllers

import (
	"backendGoAuth/internal/auth"
	"backendGoAuth/internal/models"
	"backendGoAuth/internal/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// AuthController handles authentication-related requests.
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController creates a new instance of AuthController.
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register handles the registration request.
func (controller *AuthController) Register(c *gin.Context) {
	var req models.RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ipAddress := c.ClientIP()
	authResponse, err := controller.authService.RegisterUser(req, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// Login handles the login request.
func (controller *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ipAddress := c.ClientIP() // Get client IP address
	authResponse, err := controller.authService.AuthenticateUser(req.Identifier, req.Password, ipAddress)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, authResponse)
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
	err = controller.authService.RevokeSessionToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke session token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session revoked"})
}

func (controller *AuthController) SecureEndpoint(c *gin.Context) {
	// Access user ID from the context
	userID, err := auth.GetUserIDFromTokenOrSource(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Your secure endpoint logic here
	c.JSON(http.StatusOK, gin.H{"message": "Secure Endpoint", "user_id": userID})
}

// GetActiveSessions retrieves active sessions for a user.
func (controller *AuthController) GetActiveSessions(c *gin.Context) {
	// Extract the user ID from the context or request parameters
	userID, err := auth.GetUserIDFromTokenOrSource(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("UserID from context: %v\n", userID)

	// Retrieve active sessions for the user from the database
	sessions, err := controller.authService.GetActiveSessions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve active sessions"})
		return
	}

	// Return the active sessions as JSON response
	c.JSON(http.StatusOK, sessions)
}
