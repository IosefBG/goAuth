package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

var (
	errNoToken = errors.New("no token provided")
)

type JWTMiddleware struct {
	SigningKey []byte
}

func SetupJWTMiddleware(secretKey string) *JWTMiddleware {
	return &JWTMiddleware{
		SigningKey: []byte(secretKey),
	}
}

func (jwtMiddleware *JWTMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, err := validateToken(token, jwtMiddleware.SigningKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Pass the user ID from claims to the context
		c.Set("user_id", claims["user_id"])

		c.Next()
	}
}

// extractToken extracts the JWT token from the request header.
func extractToken(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return "", errNoToken
	}

	return tokenString, nil
}

// GenerateJWT generates a new JWT token.
func GenerateJWT(secretKey string, duration time.Duration, claims map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(duration).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims.
func validateToken(tokenString string, secretKey []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GetUserIDFromContext retrieves the user ID from the Gin context.
func GetUserIDFromContext(c *gin.Context) (int, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("user ID not found in context")
	}

	return userID.(int), nil
}
