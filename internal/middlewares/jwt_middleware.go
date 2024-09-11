package middlewares

import (
	"backendGoAuth/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// JWTMiddleware holds the JWT service and secret key.
type JWTMiddleware struct {
	JwtSecret string
}

// NewJWTMiddleware creates a new instance of JWTMiddleware with the provided secret key and JwtService.
func NewJWTMiddleware(jwtSecret string) *JWTMiddleware {
	return &JWTMiddleware{
		JwtSecret: jwtSecret,
	}
}

// MiddlewareFunc returns a Gin middleware function for JWT validation.
func (jwtMiddleware *JWTMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from request cookies
		tokenString, err := utils.ExtractToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		// Validate token and retrieve claims
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract user ID from claims and set it in the context
		userID, ok := claims["user_id"].(float64) // JWT numeric claims are typically float64
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}
		c.Set("user_id", int(userID)) // Convert to int and set it in context

		c.Next()
	}
}
