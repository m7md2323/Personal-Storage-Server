package auth

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func DeviceAuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        token, err := c.Cookie("trusted_device")

       
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