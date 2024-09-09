package middlewares

import (
	"backendGoAuth/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SessionMiddleware is a struct to hold the session service dependency.
type SessionMiddleware struct {
	sessionService *services.SessionService
}

// NewSessionMiddleware initializes the SessionMiddleware with a given session service.
func NewSessionMiddleware(sessionService *services.SessionService) *SessionMiddleware {
	return &SessionMiddleware{
		sessionService: sessionService,
	}
}

// UpdateSessionMiddleware creates a middleware function that updates the session's 'updated_at' timestamp.
func (sm *SessionMiddleware) UpdateSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from the token or source
		userID, err := sm.sessionService.GetUserIDFromTokenOrSource(c)
		if err != nil {
			// Handle error
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Update the session's 'updated_at' timestamp for the user's session
		err = sm.sessionService.UpdateSessionUpdatedAt(userID)
		if err != nil {
			// Handle error
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
			return
		}

		// Continue handling the request
		c.Next()
	}
}
