package main

import (
	"Personal-Storage-Server/backend/auth"
	"Personal-Storage-Server/backend/database"
	"Personal-Storage-Server/backend/models"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"

	"path/filepath"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

		//if the code is correct then set a 30-day trusted_device Cookie
		if req.Code == auth.CurrentAccessCode && auth.CurrentAccessCode != "" {

			c.SetCookie("trusted_device", req.Code, 2592000, "/", "", false, true)
			auth.CurrentAccessCode = "" //make the code unuseable
			c.JSON(http.StatusOK, gin.H{"status": "trusted"})
		} else {
			//if not the same code then abort
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong code"})
		}
	})

	//this is for the middleware and its the first thing its checked before anything else
	private := r.Group("/")
	//if the user has the trusted_device cookie auth.DeviceAuthRequired() will execute c.Next()
	//if not it will redirect the user to the verify page
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
		//for creating a new user
		private.GET("/create_user", func(c *gin.Context) {
			c.HTML(http.StatusOK, "create_user.html", nil)
		})
		//for getting the information of the new user and store it the database
		private.POST("/api/create_user", func(c *gin.Context) {
			var newUser struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}
			
			if err := c.ShouldBindJSON(&newUser); err != nil {
				c.JSON(400, gin.H{"error": "Invalid input"})
				return
			}

			// 1. Create the folder on the Storage
			userRoot := os.Getenv("UPLOADS")
			if userRoot == "" {
				userRoot = "D:/Personal-Storage-Server/server_storage/uploads"
			}
			userPath := filepath.Join(userRoot, newUser.Username)

			err := os.MkdirAll(userPath, 0755)
			if err != nil {
				c.JSON(500, gin.H{"error": "USB Write Error: " + err.Error()})
				return
			}

			// 2. SAVE TO SQLITE DATABASE
			// Ensure 'db' is your global *sql.DB connection variable
			query := "INSERT INTO users (username, password) VALUES (?, ?)"
			database.DB.Exec(query, newUser.Username, newUser.Password)

			c.JSON(200, gin.H{"message": "User created in DB and folder initialized"})
		})
		//for getting the information of the new user and store it the database
		private.POST("/api/users", func(c *gin.Context) {
			//a temp struct to bind the new user data
			var input struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}
			//if binding went wrong abort
			if err := c.BindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
				return
			}

			//hashing the password
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

			//preparing the data to store it in the database
			newUser := models.User{
				Username: input.Username,
				Password: string(hashedPassword),
			}

			//check if the user allready exisit.
			if err := database.DB.Create(&newUser).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Username already exists"})
				return
			}
			//if everything went well, user is now created.
			c.JSON(http.StatusOK, gin.H{"message": "User created"})
		})

		private.POST("/api/login", func(c *gin.Context) {
			var loginReq struct {
				Username string
				Password string
			}
			if err := c.ShouldBindJSON(&loginReq); err != nil {
				return
			}

			var user models.User
			// Look for the user in the USB database
			if err := database.DB.Where("username = ?", loginReq.Username).First(&user).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				return
			}

			// Compare the plain text password with the hashed password
			if user.Password != loginReq.Password {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong password"})
				return
			}
			// 2. Set a Secure Cookie with the username
			c.SetCookie("user_session", loginReq.Username, 3600*24, "/", "", false, true)

			// 3. Ensure their private folder exists
			uploadDir := os.Getenv("UPLOADS")
			if uploadDir == "" {
				uploadDir = "D:/Personal-Storage-Server/server_storage/uploads"
			}
			userPath := filepath.Join(uploadDir, loginReq.Username)
			os.MkdirAll(userPath, 0777)

			c.Status(200)
		})

		private.Use(auth.UserAuthRequired())
		{
			private.POST("/api/logout", func(c *gin.Context) {
				c.SetCookie("user_session", "", -1, "/", "", false, true)
			})
			private.GET("/my_files", func(c *gin.Context) {
				c.HTML(http.StatusOK, "my_files.html", nil)
			})
			private.GET("/my_files/photos", func(c *gin.Context) {
				c.HTML(http.StatusOK, "photos.html", nil)
			})

			private.POST("/api/upload", func(c *gin.Context) {
				// 1. Get the file from the request
				file, err := c.FormFile("file")
				if err != nil {
					c.JSON(400, gin.H{"error": "No file uploaded"})
					return
				}

				// 2. Define the destination (Your USB Path)
				//dst := "/mnt/usb/server_storage/uploads/" + file.Filename
				uploadDir := os.Getenv("UPLOADS")
				if uploadDir == "" {
					uploadDir = "D:/Personal-Storage-Server/server_storage/uploads"
				}
				userSession, _ := c.Cookie("user_session")
				savePath := filepath.Join(uploadDir, userSession, file.Filename)

				if err := c.SaveUploadedFile(file, savePath); err != nil {
					fmt.Println("Upload Error:", err) // CHECK YOUR TERMINAL FOR THIS
					c.JSON(500, gin.H{"error": "Failed to save file"})
					return
				}

				c.JSON(200, gin.H{"message": "Uploaded successfully to " + file.Filename})
			})
			// Add this inside your private route group in main.go
			private.GET("/api/raw", func(c *gin.Context) {
				// This takes the path from the URL and serves the actual file
				username, _ := c.Cookie("user_session")
				uploadDir := os.Getenv("UPLOADS")
				if uploadDir == "" {
					uploadDir = "D:/Personal-Storage-Server/server_storage/uploads"
				}
				userPath := filepath.Join(uploadDir, username)
				userPath+="/"
				fileName := c.Query("name")
				//fmt.Println(userPath+fileName)
				fullPath := filepath.Join(userPath, fileName)
				// Security check: Ensure the path starts with your USB mount point
				// to prevent people from accessing your system files.
				c.File(fullPath)
			})
			private.GET("/api/explorer", func(c *gin.Context) {
				username, _ := c.Cookie("user_session")

				// Force the path to be /uploads/USERNAME

				///mnt/usb/server_storage/uploads
				uploadDir := os.Getenv("UPLOADS")
				if uploadDir == "" {
					uploadDir = "D:/Personal-Storage-Server/server_storage/uploads"
				}
				userPath := filepath.Join(uploadDir, username)

				// We ignore the 'path' query for now to force it to show the absolute usbPath
				entries, err := os.ReadDir(userPath)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   "Cannot access USB at " + userPath,
						"details": err.Error(),
					})
					return
				}

				var result []gin.H
				for _, entry := range entries {
					info, _ := entry.Info()

					sizeStr := "---"
					if !entry.IsDir() {
						sizeStr = formatBytes(info.Size())
					}

					result = append(result, gin.H{
						"name":         entry.Name(),
						"isDir":        entry.IsDir(),
						"size":         sizeStr,
						"lastModified": info.ModTime().Format("Jan 02, 2006"),
					})
				}
				c.JSON(http.StatusOK, result)
			})
		}

	}

	r.Run("0.0.0.0:8080")
}
func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
