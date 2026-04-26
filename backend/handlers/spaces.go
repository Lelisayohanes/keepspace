package handlers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"keepspace/db"
	"keepspace/models"
	"keepspace/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateSpaceRequest struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

type SpaceResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	APIKey    string `json:"api_key,omitempty"`
	CreatedAt string `json:"created_at"`
}

// generateAPIKey generates a random 64-character hex API key
func generateAPIKey() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 64 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "ks_live_" + hex.EncodeToString(bytes), nil
}

// hashAPIKey hashes an API key using SHA-256
func hashAPIKey(apiKey string) string {
	hash := sha256.Sum256([]byte(apiKey))
	return hex.EncodeToString(hash[:])
}

// CreateSpace creates a new space for the authenticated user
func CreateSpace(c *gin.Context) {
	var req CreateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Generate API key
	apiKey, err := generateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate API key"})
		return
	}

	// Hash the API key for storage
	apiKeyHash := hashAPIKey(apiKey)

	// Create space
	space := models.Space{
		OwnerID:    userID,
		Name:       req.Name,
		APIKeyHash: apiKeyHash,
	}

	if err := db.DB.Create(&space).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create space"})
		return
	}

	// Return space with plain API key (only time it's shown)
	c.JSON(http.StatusCreated, SpaceResponse{
		ID:        space.ID.String(),
		Name:      space.Name,
		APIKey:    apiKey,
		CreatedAt: space.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}

// ListSpaces returns all spaces for the authenticated user
func ListSpaces(c *gin.Context) {
	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Query spaces
	var spaces []models.Space
	if err := db.DB.Where("owner_id = ?", userID).Order("created_at DESC").Find(&spaces).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch spaces"})
		return
	}

	// Convert to response format
	response := make([]SpaceResponse, len(spaces))
	for i, space := range spaces {
		response[i] = SpaceResponse{
			ID:        space.ID.String(),
			Name:      space.Name,
			CreatedAt: space.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, response)
}

// DeleteSpace deletes a space (ownership verified)
func DeleteSpace(c *gin.Context) {
	spaceID := c.Param("id")

	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	spaceUUID, err := uuid.Parse(spaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid space ID"})
		return
	}

	// Find space and verify ownership
	var space models.Space
	if err := db.DB.Where("id = ? AND owner_id = ?", spaceUUID, userID).First(&space).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Space not found"})
		return
	}

	// Delete files from MinIO first
	if err := deleteSpaceFiles(c.Request.Context(), space.ID.String()); err != nil {
		// Log error but continue with database deletion
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete space files from storage"})
		return
	}

	// Delete space from database
	if err := db.DB.Delete(&space).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete space"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Space deleted successfully"})
}

// deleteSpaceFiles deletes all files in a space from MinIO
func deleteSpaceFiles(ctx context.Context, spaceID string) error {
	return storage.DeleteSpaceFiles(ctx, spaceID)
}
