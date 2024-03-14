package main

import (
	"backendGoAuth/internal/auth"
	"backendGoAuth/internal/controllers"
	"backendGoAuth/internal/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	// Connect to the database
	db, err := database.ConnectDB()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	// Initialize JWT middleware
	jwtMiddleware := auth.SetupJWTMiddleware("mysecretkey")

	// Create a Gin router
	router := gin.Default()

	// Set up routes using controllers
	authController := controllers.NewAuthController()
	router.POST("/login", authController.Login)
	authGroup := router.Group("/auth")
	authGroup.Use(jwtMiddleware.MiddlewareFunc())
	{
		authGroup.GET("/secure", authController.SecureEndpoint)
	}

	// Run the server
	port := 8080
	fmt.Printf("Server is running on port %d...\n", port)
	err = router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
