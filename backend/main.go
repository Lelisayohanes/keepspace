package main

import (
	"keepspace/db"
	"keepspace/handlers"
	"keepspace/middleware"
	"keepspace/storage"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Initialize database
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize MinIO
	if err := storage.InitMinio(); err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "KeepSpace API",
		})
	})

	// API Routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", handlers.Signup)
			auth.POST("/login", handlers.Login)
			auth.POST("/refresh", handlers.RefreshToken)
		}

		// Space routes (protected by JWT)
		spaces := v1.Group("/spaces")
		spaces.Use(middleware.AuthMiddleware())
		{
			spaces.GET("", handlers.ListSpaces)
			spaces.POST("", handlers.CreateSpace)
			spaces.DELETE("/:id", handlers.DeleteSpace)
		}

		// File routes (protected by API Key)
		files := v1.Group("/files")
		files.Use(middleware.APIKeyMiddleware())
		{
			files.GET("", handlers.ListFiles)
			files.POST("", handlers.UploadFile)
			files.GET("/download", handlers.DownloadFile)
			files.DELETE("", handlers.DeleteFile)
			files.GET("/presigned-url", handlers.GetPresignedURL)
		}
	}

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 KeepSpace Backend running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
