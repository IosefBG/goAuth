package middlewares

import (
	"backendGoAuth/internal/auth"
	"github.com/gin-gonic/gin"
	"os"
)

func SetupJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth.SetupJWTMiddleware(os.Getenv("JWT_SECRET"))
	}
}
