package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadFiles(c *gin.Context) {
	// 1. Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}

	// 2. Define the destination (Your USB Path)
	//dst := "/mnt/usb/server_storage/uploads/" + file.Filename
	uploadDir := os.Getenv("UPLOADS")
	userSession, _ := c.Cookie("user_session")
	savePath := filepath.Join(uploadDir, userSession, file.Filename)

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

	// Construct the absolute path
	targetPath := filepath.Join(uploadDir, username, req.Name)

	// Security strict check to ensure targetPath actually falls within the user's root folder
	// to prevent paths like "../../windows/system32"
	userRoot := filepath.Join(uploadDir, username)
	if !filepath.HasPrefix(targetPath, userRoot) {
		c.JSON(403, gin.H{"error": "Forbidden path"})
		return
	}

	if err := os.RemoveAll(targetPath); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete file"})
		return
	}

	c.JSON(200, gin.H{"message": "File deleted"})
}

func LoadUserFiles(c *gin.Context) {
	username, _ := c.Cookie("user_session")

	// Force the path to be /uploads/USERNAME

	///mnt/usb/server_storage/uploads
	uploadDir := os.Getenv("UPLOADS")
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
