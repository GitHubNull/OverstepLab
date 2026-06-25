package main

import (
	"io/fs"
	"log"
	"mime"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oversteplab/oversteplab/database"
	"github.com/oversteplab/oversteplab/database/seed"
	"github.com/oversteplab/oversteplab/internal/config"
	"github.com/oversteplab/oversteplab/internal/vuln"
	"github.com/oversteplab/oversteplab/internal/web"
	"github.com/oversteplab/oversteplab/router"
)

func main() {
	cfg := config.Load()

	// Initialize vulnerability mode from env
	vuln.SetSecureMode(cfg.SafeMode)

	// Initialize database
	db := database.Init(cfg.DBPath)

	// Auto migrate
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Seed initial data
	if err := seed.Seed(db); err != nil {
		log.Printf("Warning: Failed to seed database: %v", err)
	}

	// Setup router
	r := router.Setup(cfg)

	// Serve embedded frontend
	distFS, err := fs.Sub(web.DistFS, "dist")
	if err == nil {
		// SPA fallback - serve index.html or static files for all non-API routes
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path

			// API routes return 404
			if strings.HasPrefix(path, "/api/") {
				c.JSON(404, gin.H{"code": 404, "message": "not found"})
				return
			}

			// Try to serve static file from embedded FS
			if path == "/" {
				path = "/index.html"
			}
			// Strip leading slash for fs.Open
			filePath := strings.TrimPrefix(path, "/")
			f, err := distFS.Open(filePath)
			if err == nil {
				defer f.Close()
				info, err := f.Stat()
				if err == nil && !info.IsDir() {
					ext := strings.ToLower(filepath.Ext(filePath))
					contentType := mime.TypeByExtension(ext)
					if contentType == "" {
						switch ext {
						case ".html":
							contentType = "text/html; charset=utf-8"
						case ".js":
							contentType = "application/javascript"
						case ".css":
							contentType = "text/css"
						case ".json":
							contentType = "application/json"
						case ".png":
							contentType = "image/png"
						case ".jpg", ".jpeg":
							contentType = "image/jpeg"
						case ".svg":
							contentType = "image/svg+xml"
						case ".woff2":
							contentType = "font/woff2"
						case ".woff":
							contentType = "font/woff"
						case ".ttf":
							contentType = "font/ttf"
						case ".eot":
							contentType = "application/vnd.ms-fontobject"
						default:
							contentType = "application/octet-stream"
						}
					}
					c.DataFromReader(200, info.Size(), contentType, f, nil)
					return
				}
			}

			// Fall back to index.html (SPA)
			f2, err := distFS.Open("index.html")
			if err != nil {
				c.String(500, "Frontend not found")
				return
			}
			defer f2.Close()
			info, _ := f2.Stat()
			c.DataFromReader(200, info.Size(), "text/html; charset=utf-8", f2, nil)
		})
	}

	log.Printf("Starting server on :%s (mode: %s)", cfg.Port, vuln.GetModeString())
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
