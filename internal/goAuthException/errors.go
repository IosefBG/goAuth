package goAuthException

import (
	"net/http"
)

// Error codes
const (
	BadRequestCode          = 400
	UnauthorizedCode        = 401
	ForbiddenCode           = 403
	NotFoundCode            = 404
	Teapot                  = 418
	InternalServerErrorCode = 500
)

// Error messages
const (
	UsernameExistsMessage = "Username already exists"
	EmailExistsMessage    = "Email already exists"
	UsernameCheckError    = "Error checking username uniqueness"
	EmailCheckError       = "Error checking email uniqueness"
	HashingError          = "Error hashing password"
	UserCreationError     = "Error creating user"
	TokenGenerationError  = "Error generating JWT token"
	SessionInsertionError = "Error inserting session"
	InternalErrorMessage  = "Internal server error"
)

// CustomError represents an error with an associated error code.
type CustomError struct {
	Code    int    // Error code
	Message string // Error message
}

// Error returns the error message.
func (e *CustomError) Error() string {
	return e.Message
}

// NewCustomError creates a new CustomError with the given code and message.
func NewCustomError(code int, message string) *CustomError {
	return &CustomError{Code: code, Message: message}
}

// ErrorHandler handles errors and returns the appropriate HTTP response.
type ErrorHandler struct{}

// HandleError handles errors and returns the appropriate HTTP response.
func (eh *ErrorHandler) HandleError(err error) (int, interface{}) {
	switch e := err.(type) {
	case *CustomError:
		switch e.Code {
		case BadRequestCode:
			return http.StatusBadRequest, map[string]string{"error": e.Message}
		case UnauthorizedCode:
			return http.StatusUnauthorized, map[string]string{"error": e.Message}
		case ForbiddenCode:
			return http.StatusForbidden, map[string]string{"error": e.Message}
		case NotFoundCode:
			return http.StatusNotFound, map[string]string{"error": e.Message}
		case Teapot:
			return http.StatusTeapot, map[string]string{"error": e.Message}
		default:
			return http.StatusInternalServerError, map[string]string{"error": InternalErrorMessage}
		}
	default:
		return http.StatusInternalServerError, map[string]string{"error": InternalErrorMessage}
	}
}
