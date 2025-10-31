package routes

import (
	"authentication/controllers"
	"authentication/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
// It defines both public routes (signup, login) and protected routes (user endpoints)
// Protected routes require JWT authentication via the JWTAuthMiddleware
// Parameters:
//   - r: The Gin engine instance to configure routes on
func SetupRoutes(r *gin.Engine) {
	// Public routes - no authentication required
	r.POST("/signup", controllers.Signup()) // User registration endpoint
	r.POST("/login", controllers.Login())   // User login endpoint

	// Protected routes group - requires JWT authentication
	protected := r.Group("/")

	// Apply JWT authentication middleware to all routes in this group
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/users", controllers.GetUsers())    // Get all users (ADMIN only)
		protected.GET("/users/:id", controllers.GetUser()) // Get user by ID
	}
}
