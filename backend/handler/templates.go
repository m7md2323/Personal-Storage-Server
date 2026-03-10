package handler
import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func VerifyPage(c *gin.Context) {
	c.HTML(http.StatusOK, "verify.html", nil)
}

func IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
func CreateUserPage(c *gin.Context) {
	c.HTML(http.StatusOK, "create_user.html", nil)
}
func MyFilesPage(c *gin.Context) {
	c.HTML(http.StatusOK, "my_files.html", nil)
}

func PhotosPage(c *gin.Context) {
	c.HTML(http.StatusOK, "photos.html", nil)
}