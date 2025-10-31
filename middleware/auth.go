package middleware

import (
	"authentication/helpers"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware is a middleware function that validates JWT tokens
// It extracts the token from the Authorization header, validates it,
// and stores the claims in the request context for use by handlers
// Returns: gin.HandlerFunc that can be used as middleware
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header from the request
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix from the token string
		authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		authHeader = strings.TrimSpace(authHeader)

		// Validate the JWT token and extract claims
		claims, err := helpers.ValidateToken(authHeader)
		if err != nil {
			log.Printf("Token parsing error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store claims in context for use by handler functions
		c.Set("claims", claims)
		// Continue to the next handler
		c.Next()
	}
}
