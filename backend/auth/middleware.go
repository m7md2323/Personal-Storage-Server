package auth

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func DeviceAuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        token, err := c.Cookie("trusted_device")

        // Just check if the cookie exists. 
        // In the future, we will check if this token exists in a 'TrustedDevices' database table.
        if err != nil || token == "" {
            c.Redirect(http.StatusFound, "/verify")
            c.Abort()
            return
        }

        c.Next()
    }
}
func UserAuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        token, err := c.Cookie("user_session")

        if err != nil || token == "" {
            c.Redirect(http.StatusFound, "/")
            c.Abort()
            return
        }

        c.Next()
    }
}