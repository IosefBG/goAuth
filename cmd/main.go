package main

import (
	"backendGoAuth/internal/controllers"
	"backendGoAuth/internal/database"
	"backendGoAuth/internal/metrics"
	"backendGoAuth/internal/middlewares"
	"backendGoAuth/internal/repositories"
	"backendGoAuth/internal/services"
	"backendGoAuth/internal/utils"
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
	metrics.RegisterMetrics(reg)

	// Create Gin router
	router := setupRouter()

	// Start HTTP server
	startServer(router)
}

// setupRouter initializes and configures the Gin router
func setupRouter() *gin.Engine {
	router := gin.Default()

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Browser", "X-Device"},
		AllowCredentials: true,
	}

	router.Use(cors.New(config))

	// Apply middleware to track request duration
	router.Use(metrics.InstrumentHandler())

	// Get the database instance
	db := database.GetDB()

	// Instantiate repositories and services
	sessionRepo := repositories.NewSessionRepository(db)
	sessionService := services.NewSessionService(sessionRepo)
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, sessionService)

	// Initialize the session service in the utils package
	utils.SetSessionService(sessionRepo)

	// Initialize JWT middleware with the secret and JWT service
	jwtMiddleware := middlewares.NewJWTMiddleware(os.Getenv("JWT_SECRET"))

	// Instantiate controllers
	authController := controllers.NewAuthController(authService, sessionService)
	adminController := controllers.NewAdminController(userRepo)

	// Define routes
	api := router.Group("/api")
	{
		api.POST("/login", authController.Login)
		api.POST("/register", authController.Register)

		authGroup := api.Group("/auth", jwtMiddleware.MiddlewareFunc()) // Apply JWT middleware here
		{
			authGroup.POST("/logout", authController.RevokeCurrentSession)
			authGroup.POST("/revokeSession", authController.RevokeSession)
			authGroup.GET("/activeSessions", authController.GetActiveSessions)
			authGroup.GET("/secure", authController.SecureEndpoint)
		}

		adminGroup := api.Group("/admin", jwtMiddleware.MiddlewareFunc()) // Apply JWT middleware here
		{
			adminGroup.GET("/users", adminController.GetAllUsers)
		}
	}

	// Register Prometheus metrics endpoint
	router.GET("/metrics", metrics.MetricsHandler())

	return router
}

// startServer starts the HTTP server
func startServer(router *gin.Engine) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	fmt.Printf("Server is running on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Println("Error starting the server:", err)
	}
}
