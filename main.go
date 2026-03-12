package main

import (
	"Personal-Storage-Server/backend/auth"
	"Personal-Storage-Server/backend/database"
	"Personal-Storage-Server/backend/handler"
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// embed frontend files to make them executable
//
//go:embed frontend/*
var embedFrontend embed.FS

func main() {
	//Load Env variables
	godotenv.Load()

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

	r.MaxMultipartMemory = 10 << 20
	//Routes:
	//For untrusted devices
	r.GET("/verify", handler.VerifyPage)

	//To reques a 10-digit code to verify that the user is trusted
	r.POST("/api/request-code", handler.RequestCode)

	//This will allow the user to enter the 10-digits code and send it to the backend to check it
	r.POST("/api/submit-code", handler.SubmitCode)

	//this is for the middleware and its the first thing its checked before anything else
	private := r.Group("/")
	//if the user has the trusted_device cookie auth.DeviceAuthRequired() will execute c.Next()
	//if not it will redirect the user to the verify page
	private.Use(auth.DeviceAuthRequired())
	{
		//User selection screen
		private.GET("/", handler.IndexPage)

		// API to get the users
		private.GET("/api/users", handler.GetUsers)
		//for creating a new user
		private.GET("/create_user", handler.CreateUserPage)
		//for getting the information of the new user and store it the database
		private.POST("/api/create_user", handler.CreateUser)
		//API for login
		private.POST("/api/login", handler.Login)

		//Middleware to make sure the user is logged in
		private.Use(auth.UserAuthRequired())
		{
			private.POST("/api/logout", handler.LogOut)
			private.GET("/my_files", handler.MyFilesPage)
			private.GET("/my_files/photos", handler.PhotosPage)
			//to handle uploading files
			private.POST("/api/upload", handler.UploadFiles)
			//to handle deleteing files
			private.DELETE("/api/delete", handler.DeleteFiles)

			private.GET("/api/get_storage_info", handler.GetStorageInfo)

			private.GET("/api/raw", handler.Raw)

			private.GET("/api/explorer", handler.LoadUserFiles)
		}
	}

	r.Run("0.0.0.0:8080")
}
