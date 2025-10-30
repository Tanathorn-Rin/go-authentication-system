package main

import (
	"authentication/config"
	"authentication/helpers"
	"authentication/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Generate and set JWT key
	key := config.GenerateRandomKey()
	helpers.SetJWTKey(key)

	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	port := "8080"
	log.Println("Server running on port " + port)
	r.Run(":" + port)
}
