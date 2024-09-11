package controllers

//AuthController

import (
	"backendGoAuth/internal/models"
	"backendGoAuth/internal/services"
	"backendGoAuth/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"log"
	"net/http"
)

type AuthController struct {
	authService    *services.AuthService
	sessionService *services.SessionService
}

// NewAuthController creates a new instance of AuthController.
func NewAuthController(authService *services.AuthService, sessionService *services.SessionService) *AuthController {
	return &AuthController{
		authService:    authService,
		sessionService: sessionService,
	}
}

// Register handles the registration request.
func (controller *AuthController) Register(c *gin.Context) {
	var req models.RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	userAgent := c.GetHeader("User-Agent")
	ua := user_agent.New(userAgent)
	browser, _ := ua.Browser()
	device := ua.OS()

	ipAddress := c.ClientIP()
	authResponse, err := controller.authService.RegisterUser(req, ipAddress, browser, device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// Login handles the login request.
func (controller *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	userAgent := c.GetHeader("User-Agent")
	ua := user_agent.New(userAgent)
	browser, _ := ua.Browser()
	device := ua.OS()

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ipAddress := c.ClientIP()
	authResponse, err := controller.authService.AuthenticateUser(req.Identifier, req.Password, ipAddress, browser, device, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// RevokeCurrentSession revokes a session for the current user.
func (controller *AuthController) RevokeCurrentSession(c *gin.Context) {
	// Extract the session token from the request
	token, err := utils.ExtractToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid token"})
		return
	}

	// Revoke the session token
	err = controller.sessionService.RevokeSession(token) // Assuming token is session ID
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke session token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session revoked"})
}

func (controller *AuthController) SecureEndpoint(c *gin.Context) {
	// Access user ID from the context
	userID, err := controller.sessionService.GetUserIDFromTokenOrSource(c)
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
	userID, err := controller.sessionService.GetUserIDFromTokenOrSource(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("UserID from context: %v\n", userID)

	// Retrieve active sessions for the user from the database
	sessions, err := controller.sessionService.GetActiveSessions(userID)
	if err != nil {
		log.Printf("Error retrieving active sessions: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve active sessions"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func (ac *AuthController) RevokeSession(c *gin.Context) {
	// Extract sessionID from request (you may use query params, body, etc.)
	//todo check with the frontend later
	sessionID := c.Query("session_id") // Assuming sessionID is sent as a query parameter
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	// Call the RevokeSession method from the SessionService
	err := ac.sessionService.RevokeSession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke session"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Session revoked successfully"})
}
