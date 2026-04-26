package handlers

import (
	"fmt"
	"io"
	"keepspace/storage"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// UploadFile handles file uploads
func UploadFile(c *gin.Context) {
	// Get space ID from context (set by API key middleware)
	spaceIDStr, exists := c.Get("space_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Space not authenticated"})
		return
	}
	spaceID := spaceIDStr.(string)

	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	// Get optional path parameter
	path := c.PostForm("path")
	if path == "" {
		path = "/"
	}

	// Clean and normalize path
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	// Construct full file path
	var filePath string
	if path == "" {
		filePath = header.Filename
	} else {
		filePath = filepath.Join(path, header.Filename)
	}

	// Detect content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to MinIO
	err = storage.UploadFile(c.Request.Context(), spaceID, filePath, file, header.Size, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "File uploaded successfully",
		"filename": header.Filename,
		"path":     filePath,
		"size":     header.Size,
	})
}

// ListFiles handles listing files in a space
func ListFiles(c *gin.Context) {
	// Get space ID from context
	spaceIDStr, exists := c.Get("space_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Space not authenticated"})
		return
	}
	spaceID := spaceIDStr.(string)

	// Get optional path parameter
	path := c.Query("path")
	if path == "" {
		path = "/"
	}

	// List files from MinIO
	files, err := storage.ListFiles(c.Request.Context(), spaceID, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to list files: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"path":  path,
		"files": files,
		"count": len(files),
	})
}

// DownloadFile handles file downloads
func DownloadFile(c *gin.Context) {
	// Get space ID from context
	spaceIDStr, exists := c.Get("space_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Space not authenticated"})
		return
	}
	spaceID := spaceIDStr.(string)

	// Get file path from query or path parameter
	filePath := c.Query("path")
	if filePath == "" {
		filePath = c.Param("filepath")
	}

	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File path is required"})
		return
	}

	// Clean path
	filePath = strings.TrimPrefix(filePath, "/")

	// Get file from MinIO
	object, err := storage.GetFile(c.Request.Context(), spaceID, filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	defer object.Close()

	// Get object info for content type and size
	objectInfo, err := object.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
		return
	}

	// Set headers
	c.Header("Content-Type", objectInfo.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(filePath)))
	c.Header("Content-Length", fmt.Sprintf("%d", objectInfo.Size))

	// Stream file to response
	_, err = io.Copy(c.Writer, object)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stream file"})
		return
	}
}

// DeleteFile handles file deletion
func DeleteFile(c *gin.Context) {
	// Get space ID from context
	spaceIDStr, exists := c.Get("space_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Space not authenticated"})
		return
	}
	spaceID := spaceIDStr.(string)

	// Get file path from query
	filePath := c.Query("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File path is required"})
		return
	}

	// Clean path
	filePath = strings.TrimPrefix(filePath, "/")

	// Delete file from MinIO
	err := storage.DeleteFile(c.Request.Context(), spaceID, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete file: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted successfully",
		"path":    filePath,
	})
}

// GetPresignedURL generates a presigned URL for file download
func GetPresignedURL(c *gin.Context) {
	// Get space ID from context
	spaceIDStr, exists := c.Get("space_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Space not authenticated"})
		return
	}
	spaceID := spaceIDStr.(string)

	// Get file path from query
	filePath := c.Query("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File path is required"})
		return
	}

	// Clean path
	filePath = strings.TrimPrefix(filePath, "/")

	// Generate presigned URL
	url, err := storage.GetPresignedURL(c.Request.Context(), spaceID, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate URL: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":       url,
		"path":      filePath,
		"expires_in": "1 hour",
	})
}
