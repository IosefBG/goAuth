// internal/controllers/auth_controller.go

package controllers

import (
	"backendGoAuth/internal/auth"
	//"backendGoAuth/internal/database/postgres"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// AuthController handles authentication-related requests.
type AuthController struct {
	DB            *sql.DB
	JWTMiddleware *auth.JWTMiddleware
}

// AuthController creates a new instance of AuthController.
func NewAuthController(db *sql.DB, jwtMiddleware *auth.JWTMiddleware) *AuthController {
	return &AuthController{
		DB:            db,
		JWTMiddleware: jwtMiddleware,
	}
}

// Login handles the login request.
func (controller *AuthController) Login(c *gin.Context) {
	// Handle user authentication and generate a JWT token
	// (You will need to implement this logic based on your requirements)

	// For example purposes, we'll create a dummy token
	duration := time.Hour
	token, err := auth.GenerateJWT("mysecretkey", duration, map[string]interface{}{"user_id": 1})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// SecureEndpoint is an example of a secure endpoint that requires JWT authentication.
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
