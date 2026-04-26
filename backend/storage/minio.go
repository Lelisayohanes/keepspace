package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

const BucketName = "spaces"

// InitMinio initializes the MinIO client and ensures the bucket exists
func InitMinio() error {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "minioadmin"
	}

	secretKey := os.Getenv("MINIO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadmin"
	}

	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	// Initialize MinIO client
	var err error
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	log.Println("✓ Connected to MinIO")

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := MinioClient.BucketExists(ctx, BucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = MinioClient.MakeBucket(ctx, BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("✓ Created bucket: %s", BucketName)
	} else {
		log.Printf("✓ Bucket '%s' already exists", BucketName)
	}

	return nil
}

// UploadFile uploads a file to MinIO
func UploadFile(ctx context.Context, spaceID, filePath string, reader io.Reader, size int64, contentType string) error {
	objectName := fmt.Sprintf("spaces/%s/%s", spaceID, filePath)

	_, err := MinioClient.PutObject(ctx, BucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

// FileInfo represents file metadata
type FileInfo struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	Size         int64  `json:"size"`
	LastModified string `json:"last_modified"`
	IsFolder     bool   `json:"is_folder"`
}

// ListFiles lists files in a space with optional path prefix
func ListFiles(ctx context.Context, spaceID, pathPrefix string) ([]FileInfo, error) {
	prefix := fmt.Sprintf("spaces/%s/", spaceID)
	if pathPrefix != "" && pathPrefix != "/" {
		pathPrefix = strings.TrimPrefix(pathPrefix, "/")
		pathPrefix = strings.TrimSuffix(pathPrefix, "/")
		prefix = fmt.Sprintf("spaces/%s/%s/", spaceID, pathPrefix)
	}

	var files []FileInfo
	seenFolders := make(map[string]bool)

	// List objects with prefix
	objectCh := MinioClient.ListObjects(ctx, BucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: false,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("error listing objects: %w", object.Err)
		}

		// Remove the prefix to get relative path
		relativePath := strings.TrimPrefix(object.Key, prefix)
		if relativePath == "" {
			continue
		}

		// Check if it's a folder (ends with /)
		if strings.HasSuffix(object.Key, "/") {
			folderName := strings.TrimSuffix(relativePath, "/")
			if !seenFolders[folderName] {
				files = append(files, FileInfo{
					Name:     folderName,
					Path:     object.Key,
					IsFolder: true,
				})
				seenFolders[folderName] = true
			}
		} else {
			// It's a file
			fileName := filepath.Base(relativePath)
			files = append(files, FileInfo{
				Name:         fileName,
				Path:         object.Key,
				Size:         object.Size,
				LastModified: object.LastModified.Format("2006-01-02 15:04:05"),
				IsFolder:     false,
			})
		}
	}

	return files, nil
}

// GetFile retrieves a file from MinIO
func GetFile(ctx context.Context, spaceID, filePath string) (*minio.Object, error) {
	objectName := fmt.Sprintf("spaces/%s/%s", spaceID, filePath)

	object, err := MinioClient.GetObject(ctx, BucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return object, nil
}

// DeleteFile deletes a single file from MinIO
func DeleteFile(ctx context.Context, spaceID, filePath string) error {
	objectName := fmt.Sprintf("spaces/%s/%s", spaceID, filePath)

	err := MinioClient.RemoveObject(ctx, BucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// DeleteSpaceFiles deletes all files in a space
func DeleteSpaceFiles(ctx context.Context, spaceID string) error {
	prefix := fmt.Sprintf("spaces/%s/", spaceID)

	objectCh := MinioClient.ListObjects(ctx, BucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	errorCh := MinioClient.RemoveObjects(ctx, BucketName, objectCh, minio.RemoveObjectsOptions{})

	for err := range errorCh {
		if err.Err != nil {
			return fmt.Errorf("failed to delete object %s: %w", err.ObjectName, err.Err)
		}
	}

	return nil
}

// GetPresignedURL generates a presigned URL for downloading a file
func GetPresignedURL(ctx context.Context, spaceID, filePath string) (string, error) {
	objectName := fmt.Sprintf("spaces/%s/%s", spaceID, filePath)

	// Generate presigned URL valid for 1 hour
	url, err := MinioClient.PresignedGetObject(ctx, BucketName, objectName, 3600, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}
