package main

import (
	"Personal-Storage-Server/backend/auth"
	"Personal-Storage-Server/backend/database"
	"Personal-Storage-Server/backend/models"
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// embed frontend files to make them executable
//go:embed frontend/*
var embedFrontend embed.FS

func main() {
	//Init database
	database.ConnectDatabase()

	//Create router
	r := gin.Default()

	//Setup Embedded Files
	frontendFS, _ := fs.Sub(embedFrontend, "frontend")
	tmpl := template.Must(template.ParseFS(frontendFS, "*.html"))
	r.SetHTMLTemplate(tmpl)

	//Static files Path
	staticFS, _ := fs.Sub(frontendFS, "static")
	r.StaticFS("/static", http.FS(staticFS))

	//Routes:
	//For untrusted devices
	r.GET("/verify", func(c *gin.Context) {
		c.HTML(http.StatusOK, "verify.html", nil)
	})

	//To reques a 10-digit code to verify that the user is trusted
	r.POST("/api/request-code", func(c *gin.Context) {
		//this function will print 10-digits code on the server's terminal
		auth.GenerateStorageCode()
		c.JSON(http.StatusOK, gin.H{"status": "printed"})
	})

	//This will allow the user to enter the 10-digits code and send it to the backend to check it
	r.POST("/api/submit-code", func(c *gin.Context) {
		var req struct {
			Code string `json:"code"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		//f the code is correct then set a 30-day trusted_device Cookie
		if req.Code == auth.CurrentAccessCode && auth.CurrentAccessCode != "" {

			c.SetCookie("trusted_device", req.Code, 2592000, "/", "", false, true)
			auth.CurrentAccessCode = "" //make the code unuseable
			c.JSON(http.StatusOK, gin.H{"status": "trusted"})
		} else {
			//if not the same code then abort
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong code"})
		}
	})

	//this is for the middleware
	private := r.Group("/")
	private.Use(auth.DeviceAuthRequired())
	{
		//User selection screen
		private.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		})

		// API to get the users
		private.GET("/api/users", func(c *gin.Context) {
			var users []models.User
			database.DB.Find(&users)
			c.JSON(http.StatusOK, users)
		})
	}

	r.Run("0.0.0.0:8080")
}
