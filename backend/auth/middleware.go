package auth

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func DeviceAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("trusted_device")

		// If no cookie, or wrong cookie, redirect to the verify page
		if err != nil || token == "" {
            // Redirect instead of just throwing an error
			c.Redirect(http.StatusFound, "/verify")
			c.Abort()
			return
		}

        // We will replace this with real database validation later
		if token != "123-456-7890" {
			c.Redirect(http.StatusFound, "/verify")
			c.Abort()
			return
		}

        // Device is trusted! Let them see the Netflix screen.
		c.Next()
	}
}