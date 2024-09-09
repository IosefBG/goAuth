// main.go

package main

import (
	"backendGoAuth/internal/controllers"
	"backendGoAuth/internal/database"
	prometheusmetrics "backendGoAuth/internal/metrics"
	"backendGoAuth/internal/middlewares"
	"backendGoAuth/internal/repositories"
	"backendGoAuth/internal/services"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file:", err)
	}

	// Connect to the database
	if err := database.ConnectDB(); err != nil {
		log.Println("Error connecting to the database:", err)
		return
	}

	defer func() {
		if err := database.GetDB().Close(); err != nil {
			log.Println("Error closing the database connection:", err)
		}
	}()

	// Initialize Prometheus metrics registry
	reg := prometheus.NewRegistry()

	// Register Prometheus metrics
	prometheusmetrics.RegisterMetrics(reg)

	// Create Gin router
	router := setupRouter()

	// Start HTTP server
	startServer(router)
}

// setupRouter initializes and configures the Gin router
func setupRouter() *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Browser", "X-Device"} // Include Authorization header
	// Enable CORS for all origins, methods, and headers
	router.Use(cors.New(config))

	// Apply middleware to track request duration
	router.Use(prometheusmetrics.InstrumentHandler())

	// Setup JWT middleware
	router.Use(middlewares.SetupJWTMiddleware(os.Getenv("JWT_SECRET")))

	// Get the database instance
	db := database.GetDB()

	// Instantiate AuthService
	sessionRepo := repositories.NewSessionRepository(db) // Corrected to use pointer
	sessionService := services.NewSessionService(sessionRepo)
	userRepo := repositories.NewUserRepository(db) // Initialize UserRepo
	authService := services.NewAuthService(userRepo, sessionService)

	// Instantiate middleware with session service
	sessionMiddleware := middlewares.NewSessionMiddleware(sessionService)

	// Instantiate AuthController with AuthService
	authController := controllers.NewAuthController(authService, sessionService)
	adminController := controllers.NewAdminController(userRepo)

	// Define routes
	api := router.Group("/api")
	{
		api.POST("/login", authController.Login)
		api.POST("/register", authController.Register)

		authGroup := api.Group("/auth", sessionMiddleware.UpdateSessionMiddleware())
		{
			authGroup.POST("/logout", authController.RevokeCurrentSession)
			authGroup.POST("/revokeSession", authController.RevokeSession)
			authGroup.GET("/activeSessions", authController.GetActiveSessions)
			authGroup.GET("/secure", authController.SecureEndpoint)
		}

		adminGroup := api.Group("/admin", sessionMiddleware.UpdateSessionMiddleware())
		{
			adminGroup.GET("/users", adminController.GetAllUsers)
		}
	}

	// Register Prometheus metrics endpoint
	router.GET("/metrics", prometheusmetrics.MetricsHandler())

	return router
}

// startServer starts the HTTP server
func startServer(router *gin.Engine) {
	port := os.Getenv("PORT")
	fmt.Printf("Server is running on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Println("Error starting the server:", err)
	}
}
