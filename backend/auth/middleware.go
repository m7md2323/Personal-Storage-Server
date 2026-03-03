package auth

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("trust_token")

		if err != nil || token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		if token != "123-456-7890" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}

		c.Next()
	}
}