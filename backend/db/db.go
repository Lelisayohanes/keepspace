package db

import (
	"fmt"
	"log"
	"os"

	"keepspace/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection and runs migrations
func InitDB() error {
	// Get database URL from environment variable
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Default for local development
		databaseURL = "postgres://keepspace:devpassword@localhost:5432/keepspace_dev?sslmode=disable"
	}

	// Open database connection
	var err error
	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("✓ Connected to PostgreSQL database")

	// Run auto-migrations
	err = DB.AutoMigrate(&models.User{}, &models.Space{})
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✓ Database migrations completed")

	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
