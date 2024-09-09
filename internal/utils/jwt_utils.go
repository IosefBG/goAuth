package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var (
	errNoToken = errors.New("no token provided")
)

var (
	JwtSecret   string
	jwtDuration time.Duration
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	JwtSecret = os.Getenv("JWT_SECRET")
	jwtDurationStr := os.Getenv("JWT_DURATION_HOURS")

	log.Printf("JWT_SECRET: %s", JwtSecret)
	log.Printf("JWT_DURATION_HOURS: %s", jwtDurationStr)

	jwtDuration, err = time.ParseDuration(jwtDurationStr + "h")
	if err != nil {
		log.Fatalf("Error parsing JWT duration: %s", err)
	}
	log.Printf("JWT duration: %s", jwtDuration)
}

// JWTMiddleware handles JWT token validation and extraction.
type JWTMiddleware struct {
	SigningKey []byte
}

//// SetupJWTMiddleware initializes JWT middleware with the provided secret key.
//func SetupJWTMiddleware(secretKey string) *JWTMiddleware {
//	return &JWTMiddleware{
//		SigningKey: []byte(secretKey),
//	}
//}

// MiddlewareFunc returns a Gin middleware function for JWT validation.
func (jwtMiddleware *JWTMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := ExtractToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, err := ValidateToken(token, jwtMiddleware.SigningKey)
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

// ExtractToken extracts the JWT token from the request header.
func ExtractToken(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return "", errNoToken
	}

	const bearerPrefix = "Bearer "
	if strings.HasPrefix(tokenString, bearerPrefix) {
		// Remove the 'Bearer ' prefix
		tokenString = tokenString[len(bearerPrefix):]
	}

	return tokenString, nil
}

// GenerateJWT generates a new JWT token with the provided claims.
func GenerateJWT(claims jwt.MapClaims) (string, error) {
	// Set the expiration time for the token
	claims["exp"] = time.Now().Add(jwtDuration).Unix()

	// Create a new JWT token with the provided claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(JwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims.
func ValidateToken(tokenString string, secretKey []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		fmt.Println("Error validating token:", err) // Log the error for debugging
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token") // Token is invalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims") // Claims are not of expected type
	}

	return claims, nil // Token is valid and claims are retrieved successfully
}
