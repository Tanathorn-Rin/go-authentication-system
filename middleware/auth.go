package middleware

import (
	"authentication/helpers"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}
		// Remove "Bearer " prefix
		authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		authHeader = strings.TrimSpace(authHeader)

		claims, err := helpers.ValidateToken(authHeader)
		if err != nil {
			log.Printf("Token parsing error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
