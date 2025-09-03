package storage

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageService struct {
	client     *minio.Client
	bucketName string
	publicURL  string
	region     string
}

type UploadResult struct {
	URL       string
	Key       string
	Size      int64
	MediaType string
}

func NewStorageService(config StorageConfig) (*StorageService, error) {
	// Initialize minio client
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
		Region: config.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	service := &StorageService{
		client:     client,
		bucketName: config.BucketName,
		publicURL:  config.PublicURL,
		region:     config.Region,
	}

	// Ensure bucket exists
	if err := service.ensureBucketExists(); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	return service, nil
}

func (s *StorageService) ensureBucketExists() error {
	ctx := context.Background()

	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{
			Region: s.region,
		})
		if err != nil {
			return err
		}
	}

	// Set bucket policy to allow public read access to uploaded files
	policy := fmt.Sprintf(`{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": "*",
                "Action": ["s3:GetObject"],
                "Resource": ["arn:aws:s3:::%s/*"]
            }
        ]
    }`, s.bucketName)

	return s.client.SetBucketPolicy(ctx, s.bucketName, policy)
}
