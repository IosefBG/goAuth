// In your controllers/auth_controller.go

package controllers

import (
	"backendGoAuth/internal/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthController handles authentication-related requests.
type AuthController struct {
	// Add any dependencies here
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

	// Store the session token along with the user ID
	err = storeSessionToken(userID, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store session token"})
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

// authenticateUser authenticates the user based on the provided credentials.
func authenticateUser(username, password string) (int, error) {
	// Implement logic to authenticate the user
	// Example: user, err := getUserByUsername(username)
	// if err != nil || !checkPassword(user.Password, password) {
	//     return 0, errors.New("authentication failed")
	// }
	// return user.ID, nil

	// For simplicity, return a hardcoded user ID
	return 1, nil
}

// storeSessionToken stores the session token in a database.
func storeSessionToken(userID int, token string) error {
	// Implement logic to store the session token in a database
	// Example: db.Exec("INSERT INTO sessions (user_id, token) VALUES (?, ?)", userID, token)
	return nil
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
