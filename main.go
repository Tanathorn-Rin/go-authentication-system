package main

import (
	"authentication/config"
	"authentication/helpers"
	"authentication/routes"
	"log"

	"github.com/gin-gonic/gin"
)

// main is the entry point of the application
// It initializes the JWT key, sets up the Gin router with routes,
// and starts the HTTP server on port 8080
func main() {
	// Generate and set JWT key
	// This creates a random 32-byte key used for signing JWT tokens
	key := config.GenerateRandomKey()
	helpers.SetJWTKey(key)

	// Initialize Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Setup routes (public and protected endpoints)
	routes.SetupRoutes(r)

	// Configure and start the HTTP server
	port := "8080"
	log.Println("Server running on port " + port)
	r.Run(":" + port)
}
