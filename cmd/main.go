package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"

	"backendGoAuth/internal/controllers"
	"backendGoAuth/internal/database"
	"backendGoAuth/internal/metrics"
	"backendGoAuth/internal/middlewares"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file:", err)
	}

	// Connect to the database
	db, err := database.ConnectDB()
	if err != nil {
		log.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	// Initialize JWT middleware
	jwtMiddleware := middlewares.SetupJWTMiddleware()

	// Create Prometheus metrics registry
	reg := prometheus.NewRegistry()

	// Register Prometheus metrics
	prometheusmetrics.RegisterMetrics(reg)

	// Create Gin router
	router := gin.Default()

	// Initialize controllers
	authController := controllers.NewAuthController()

	// Define routes
	router.POST("/login", authController.Login)
	authGroup := router.Group("/auth")
	authGroup.Use(jwtMiddleware)
	{
		authGroup.GET("/secure", authController.SecureEndpoint)
	}

	// Register Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})))

	// Start HTTP server
	port := 8080
	fmt.Printf("Server is running on port %d...\n", port)
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Println("Error starting the server:", err)
	}
}
