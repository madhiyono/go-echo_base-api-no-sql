package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *StorageService) UploadProfilePhoto(userID primitive.ObjectID, file multipart.File, fileSize int64, filename string) (*UploadResult, error) {
	ctx := context.Background()

	// Generate unique key for the file
	timestamp := time.Now().Unix()
	ext := filepath.Ext(filename)
	key := fmt.Sprintf("profiles/%s_%d%s", userID.Hex(), timestamp, ext)

	// Get file info
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read file header: %w", err)
	}
	_, err = file.Seek(0, 0) // Reset file pointer
	if err != nil {
		return nil, fmt.Errorf("failed to reset file pointer: %w", err)
	}

	contentType := http.DetectContentType(buffer)

	// Upload file
	info, err := s.client.PutObject(ctx, s.bucketName, key, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Generate public URL
	url := fmt.Sprintf("%s/%s/%s", s.publicURL, s.bucketName, key)

	return &UploadResult{
		URL:       url,
		Key:       key,
		Size:      info.Size,
		MediaType: contentType,
	}, nil
}

func (s *StorageService) DeleteProfilePhoto(key string) error {
	ctx := context.Background()
	return s.client.RemoveObject(ctx, s.bucketName, key, minio.RemoveObjectOptions{})
}

func (s *StorageService) GetPublicURL(key string) string {
	return fmt.Sprintf("%s/%s/%s", s.publicURL, s.bucketName, key)
}
