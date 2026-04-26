package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Mock data and types for the MVP Go Backend structure
// In a real scenario, this would connect to Postgres and MinIO
type Space struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API Routes
	v1 := r.Group("/api/v1")
	{
		// Auth
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", func(c *gin.Context) {
				c.JSON(http.StatusCreated, gin.H{"message": "User created"})
			})
			auth.POST("/login", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"access_token":  "mock_access_token",
					"refresh_token": "mock_refresh_token",
				})
			})
		}

		// Spaces (Protected by JWT)
		spaces := v1.Group("/spaces")
		{
			spaces.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, []Space{
					{ID: "1", Name: "Personal Documents", CreatedAt: time.Now()},
					{ID: "2", Name: "Project Assets", CreatedAt: time.Now()},
				})
			})
			spaces.POST("", func(c *gin.Context) {
				c.JSON(http.StatusCreated, Space{
					ID:        "3",
					Name:      "New Space",
					APIKey:    "ks_live_5173_fake_key_long_random_string",
					CreatedAt: time.Now(),
				})
			})
			spaces.DELETE("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Space deleted"})
			})
		}

		// Files (Protected by API Key)
		files := v1.Group("/files")
		{
			files.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"files": []string{"demo.txt", "image.png"}})
			})
			files.POST("", func(c *gin.Context) {
				c.JSON(http.StatusCreated, gin.H{"message": "File uploaded"})
			})
		}
	}

	log.Println("KeepSpace Backend running on :8080")
	// r.Run(":8080") // Uncomment to run
}