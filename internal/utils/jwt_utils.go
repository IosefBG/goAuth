package utils

//JwtUtils

import (
	"backendGoAuth/internal/repositories"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	errNoToken = errors.New("no token provided")
)

var (
	JwtSecret   string
	jwtDuration time.Duration
	sessionSvc  *repositories.SessionRepository
)

func init() {
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

func SetSessionService(repo *repositories.SessionRepository) {
	sessionSvc = repo
}

// ExtractToken extracts the JWT token from the cookies.
func ExtractToken(c *gin.Context) (string, error) {
	tokenString, err := c.Cookie("access_token")
	if err != nil {
		return "", errNoToken
	}

	return tokenString, nil
}

// GenerateJWT generates a new JWT token with the specified claims and token type.
func GenerateJWT(claims jwt.MapClaims, tokenType string, sessionID int) (string, error) {
	// Set the expiration time for the token
	if tokenType == "access" {
		claims["exp"] = time.Now().Add(jwtDuration).Unix()
	} else if tokenType == "refresh" {
		// Longer duration for refresh tokens
		claims["exp"] = time.Now().Add(7 * 24 * time.Hour).Unix() // 7 days for refresh token
	}

	// Include the session_id claim
	claims["session_id"] = sessionID

	// Create a new JWT token with the provided claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(JwtSecret))
	if err != nil {
		fmt.Println("Error signing token:", err) // Log error if signing fails
		return "", err
	}

	return tokenString, nil
}

// SetJWTTokenCookies sets JWT access and refresh tokens as HttpOnly cookies.
func SetJWTTokenCookies(c *gin.Context, accessToken string) {
	// Set access token as HttpOnly cookie
	c.SetCookie("access_token", accessToken, int(jwtDuration.Seconds()), "/", "", false, true) //if we use prodlike https secure should be true
	// Optionally set a refresh token cookie if you have one
	// c.SetCookie("refresh_token", refreshToken, 7*24*60*60, "/", "", true, true) // 7 days for refresh token
}

// ValidateToken validates a JWT token and returns the claims.
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token uses the correct signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Unexpected signing method:", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key as a byte slice
		return []byte(JwtSecret), nil
	})

	if err != nil {
		fmt.Println("Error validating token:", err) // Log the error for debugging
		return nil, err
	}

	if !token.Valid {
		fmt.Println("Token is invalid")
		return nil, errors.New("invalid token") // Token is invalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Invalid claims type")
		return nil, errors.New("invalid claims") // Claims are not of expected type
	}

	// Print the token claims for debugging
	fmt.Println("Token claims:", claims)
	fmt.Println("JwtSecret:", JwtSecret)

	//fixme something wrong happends here after user sucesfully logged in when trying to get these
	// Extract session ID from token (ensure `getSessionIDFromToken` works correctly)
	sessionID, err := getSessionIDFromToken(tokenString)
	if err != nil {
		fmt.Println("Session not found:", err)
		return nil, errors.New("session not found") // Session ID could not be retrieved
	}

	if sessionSvc == nil {
		fmt.Println("sessionSvc is nil")
	}

	// Check if the session is revoked
	session, err := sessionSvc.GetSessionByID(sessionID)
	if err != nil {
		fmt.Println("Error fetching session:", err)
		return nil, errors.New("error fetching session") // Error fetching session
	}

	println(session)

	if !session.IsActive {
		fmt.Println("Session is revoked")
		return nil, errors.New("session revoked") // Session is revoked
	}

	return claims, nil // Token is valid and claims are retrieved successfully
}

// getSessionIDFromToken extracts the session ID from the provided JWT token.
func getSessionIDFromToken(tokenString string) (int, error) {
	// Parse the token without verifying the signature to extract claims.
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return 0, fmt.Errorf("error parsing token: %v", err)
	}

	// Assert that the claims are of type jwt.MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	// Extract the session ID from the claims.
	sessionIDFloat, ok := claims["session_id"].(float64) // JWT numeric claims are typically float64
	if !ok {
		return 0, errors.New("session ID not found in token")
	}

	return int(sessionIDFloat), nil // Convert float64 to int
}
