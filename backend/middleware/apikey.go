package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"keepspace/db"
	"keepspace/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// hashAPIKey hashes an API key using SHA-256
func hashAPIKey(apiKey string) string {
	hash := sha256.Sum256([]byte(apiKey))
	return hex.EncodeToString(hash[:])
}

// APIKeyMiddleware validates API keys for file operations
func APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "X-API-Key header required"})
			c.Abort()
			return
		}

		// Hash the provided API key
		apiKeyHash := hashAPIKey(apiKey)

		// Find space by API key hash
		var space models.Space
		if err := db.DB.Where("api_key_hash = ?", apiKeyHash).First(&space).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		// Store space info in context
		c.Set("space_id", space.ID.String())
		c.Set("space_owner_id", space.OwnerID.String())
		c.Set("space_name", space.Name)

		c.Next()
	}
}

// GetSpaceIDFromContext retrieves space ID from context
func GetSpaceIDFromContext(c *gin.Context) (uuid.UUID, error) {
	spaceIDStr, exists := c.Get("space_id")
	if !exists {
		return uuid.Nil, nil
	}
	return uuid.Parse(spaceIDStr.(string))
}
