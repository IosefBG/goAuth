package middlewares

import (
	"backendGoAuth/internal/auth"
	"backendGoAuth/internal/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from the token or source
		userID, err := auth.GetUserIDFromTokenOrSource(c)
		if err != nil {
			// Handle error
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Update the updated_at column for the user's session
		err = database.UpdateSessionUpdatedAt(userID)
		if err != nil {
			// Handle error
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
			return
		}

		// Continue handling the request
		c.Next()
	}
}
