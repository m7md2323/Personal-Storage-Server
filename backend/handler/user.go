package handler

import (
	"Personal-Storage-Server/backend/auth"
	"Personal-Storage-Server/backend/database"
	"Personal-Storage-Server/backend/models"
	"net/http"
	"os"
	"path/filepath"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SubmitCode(c *gin.Context){
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
}
func RequestCode(c *gin.Context) {
		//this function will print 10-digits code on the server's terminal
		auth.GenerateStorageCode()
		c.JSON(http.StatusOK, gin.H{"status": "printed"})
}		
func CreateUser(c *gin.Context) {
	var userInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)

	newUser:=models.User{
		Username: userInput.Username,
		Password: string(hashedPassword),
	}
	
		
	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username already exists"})
		return
	}
	// 1. Create the folder on the Storage
	userRoot := os.Getenv("UPLOADS")
	userPath := filepath.Join(userRoot, newUser.Username)

	err := os.MkdirAll(userPath, 0755)
	if err != nil {
		c.JSON(500, gin.H{"error": "USB Write Error: " + err.Error()})
		return
	}

	// 2. SAVE TO SQLITE DATABASE
	// Ensure 'db' is your global *sql.DB connection variable
	//query := "INSERT INTO users (username, password) VALUES (?, ?)"
	//database.DB.Exec(query, newUser.Username, newUser.Password)

	c.JSON(200, gin.H{"message": "User created in DB and folder initialized"})
}
func GetUsers(c *gin.Context) {
	var users []models.User
	database.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}
func Login(c *gin.Context) {
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
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))!=nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong password"})
		return
	}
	// 2. Set a Secure Cookie with the username
	c.SetCookie("user_session", loginReq.Username, 3600*24, "/", "", false, true)

	// 3. Ensure their private folder exists
	uploadDir := os.Getenv("UPLOADS")
	userPath := filepath.Join(uploadDir, loginReq.Username)
	os.MkdirAll(userPath, 0777)

	c.Status(200)
}

func LogOut(c *gin.Context) {
	//when user logs out remove his Cookie user_session
	c.SetCookie("user_session", "", -1, "/", "", false, true)
}
