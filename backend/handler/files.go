package handler

import (
	"Personal-Storage-Server/backend/database"
	"Personal-Storage-Server/backend/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"github.com/gin-gonic/gin"
)

func UploadFiles(c *gin.Context) {
	//Get the file from the user request (only one file at a time)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}
	//to avoid path traversal attaks
	safeFileName := filepath.Base(file.Filename)

	//look for the uploads directory
	uploadDir := os.Getenv("UPLOADS")
	//fetch the username from user_session Cookie to put the new file in uploads/username/newFile
	userSession, _ := c.Cookie("user_session")
	//safely join the paths together, this will be the final full path to save the new File
	savePath := filepath.Join(uploadDir, userSession, safeFileName)
	/*	ID        uint           `gorm:"primaryKey" json:"id"`

		FileName  string `json:"file_name"`
		FilePath  string `json:"file_path"`
		FileSize  int64  `json:"file_size"`
		FileType  string `json:"file_type"`

		OwnerID   uint   `json:"owner_id"`
		}*/
	newFile := models.File{
		FileName:      file.Filename,
		FilePath:      savePath,
		FileSize:      file.Size,
		FileType:      filepath.Ext(file.Filename),
		OwnerUsername: userSession,
	}

	if err := database.DB.Create(&newFile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File already exists"})
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer src.Close()

	os.MkdirAll(filepath.Dir(savePath), os.ModePerm)

	out, err := os.Create(savePath)
	if err != nil {
		fmt.Println("File Create Error:", err)
		c.JSON(500, gin.H{"error": "Failed to create save destination"})
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		fmt.Println("Upload Error:", err)
		c.JSON(500, gin.H{"error": "Failed to write file"})
		return
	}

	c.JSON(200, gin.H{"message": "Uploaded successfully to " + file.Filename})
}

func DeleteFiles(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	username, _ := c.Cookie("user_session")
	uploadDir := os.Getenv("UPLOADS")
	safeName:=filepath.Base(req.Name)
	// Construct the absolute path
	targetPath := filepath.Join(uploadDir, username, safeName)

	// Security strict check to ensure targetPath actually falls within the user's root folder
	// to prevent paths like "../../windows/system32"
	//userRoot := filepath.Join(uploadDir, username)
	//if !filepath.HasPrefix(targetPath, userRoot) {
	//	c.JSON(403, gin.H{"error": "Forbidden path"})
	//	return
	//}

	database.DB.Where("file_name = ?", safeName).Delete(&models.File{})
	if err := os.RemoveAll(targetPath); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete file"})
		return
	}

	c.JSON(200, gin.H{"message": "File deleted"})
}

func LoadUserFiles(c *gin.Context) {

	username, _ := c.Cookie("user_session")
	uploadDir := os.Getenv("UPLOADS")
	userPath := filepath.Join(uploadDir, username)

	_, err := os.ReadDir(userPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Cannot access storage at " + userPath,
			"details": err.Error(),
		})
		return
	}
	var Files [] models.File
	database.DB.Where("owner_username = ?", username).Find(&Files)

	var result []gin.H
	for _,file:=range Files{
		result = append(result, gin.H{
			"name":         file.FileName,
			"ext":     		file.FileType,
			"size":         file.FileSize,
			"path": 		file.FilePath,
		})
	}

	c.JSON(http.StatusOK, result)
}

func Raw(c *gin.Context) {
	// This takes the path from the URL and serves the actual file
	username, _ := c.Cookie("user_session")
	uploadDir := os.Getenv("UPLOADS")
	userPath := filepath.Join(uploadDir, username)
	userPath += "/"
	fileName := c.Query("name")
	//fmt.Println(userPath+fileName)
	fullPath := filepath.Join(userPath, fileName)
	// Security check: Ensure the path starts with your USB mount point
	// to prevent people from accessing your system files.
	c.File(fullPath)
}

func GetStorageInfo(c *gin.Context){
	var stat syscall.Statfs_t
		path := os.Getenv("UPLOADS")
		err := syscall.Statfs(path, &stat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read disk"})
			return
		}

		// Available blocks * size per block = available space in bytes
		availableBytes := stat.Bavail * uint64(stat.Bsize)
		totalBytes := stat.Blocks * uint64(stat.Bsize)
		freeBytes := stat.Bfree * uint64(stat.Bsize)

		totalGigaBytes:=float32(totalBytes)/ (1024 * 1024 * 1024);
		avaGigaBytes:=float32(availableBytes)/ (1024 * 1024 * 1024);
		freeGigaBytes:=float32(freeBytes)/ (1024 * 1024 * 1024);
		c.JSON(http.StatusOK, gin.H{
			"path":           path,
			"total_gb":       totalGigaBytes,
			"available_gb":   avaGigaBytes,
			"free_gb":        freeGigaBytes,
		})
}
