package routes

import (
	"authentication/controllers"
	"authentication/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/signup", controllers.Signup())
	r.POST("/login", controllers.Login())

	protected := r.Group("/")

	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/users", controllers.GetUsers())
		protected.GET("/users/:id", controllers.GetUser())
	}
}
