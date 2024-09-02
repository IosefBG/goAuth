package main

import (
	"backendGoAuth/internal/controllers"
	"backendGoAuth/internal/database"
	"backendGoAuth/internal/metrics"
	"backendGoAuth/internal/middlewares"
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

	//router.Use(middlewares.UpdateSessionMiddleware())

	// Apply middleware to track request duration
	router.Use(prometheusmetrics.InstrumentHandler())

	// Setup JWT middleware
	router.Use(middlewares.SetupJWTMiddleware(os.Getenv("JWT_SECRET")))

	// Instantiate AuthService
	sessionService := services.NewSessionService()
	authService := services.NewAuthService(sessionService)

	// Instantiate AuthController with AuthService
	authController := controllers.NewAuthController(authService, sessionService)

	// Define routes
	api := router.Group("/api")
	{
		api.POST("/login", authController.Login)
		api.POST("/register", authController.Register)
		api.POST("/logout", middlewares.UpdateSessionMiddleware(), authController.RevokeCurrentSession)
		api.POST("/revokeSession", middlewares.UpdateSessionMiddleware(), authController.RevokeSession)
		api.GET("/activeSessions", middlewares.UpdateSessionMiddleware(), authController.GetActiveSessions)

		authGroup := api.Group("/auth", middlewares.UpdateSessionMiddleware())
		{
			authGroup.GET("/secure", authController.SecureEndpoint)
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
