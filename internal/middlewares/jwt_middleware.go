package middlewares

import (
	"backendGoAuth/internal/auth"
	"github.com/gin-gonic/gin"
)

func SetupJWTMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth.SetupJWTMiddleware(jwtSecret)
	}
}
