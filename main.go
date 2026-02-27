package main

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"github.com/gin-gonic/gin"
)

// This matches your frontend/embed.go logic
//go:embed frontend/*
var embedFrontend embed.FS

func main() {
	r := gin.Default()

	// 1. Extract the 'frontend' subfolder from the embedded files
	frontendFS, _ := fs.Sub(embedFrontend, "frontend")
	
	// 2. Load the HTML templates from the embedded file system
	tmpl := template.Must(template.ParseFS(frontendFS, "*.html"))
	r.SetHTMLTemplate(tmpl)

	// 3. Serve static files (CSS/JS) from the embedded system
	r.StaticFS("/static", http.FS(frontendFS))

	// 4. THE ROUTE: Serve the Introduction Page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Welcome to My Cloud",
		})
	})

	// 5. RUN: Listen on all interfaces so your other PC can see it
	// Using 0.0.0.0 ensures it is not locked to 'localhost'
	r.Run("0.0.0.0:8080")
}