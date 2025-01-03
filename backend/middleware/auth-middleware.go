package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is used to protect routes that require authentication
func AuthMiddleware(c *gin.Context) {
	// Here, you can add JWT token verification logic
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	c.Next() // Proceed with the request if token is valid
}
