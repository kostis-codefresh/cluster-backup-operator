package minio

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	client *minio.Client
}

func NewMinIOClient(endpoint, accessKey, secretKey string, useSSL bool) (*MinIOClient, error) {
	// Initialize MinIO client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating MinIO client: %w", err)
	}

	return &MinIOClient{client: client}, nil
}

// UploadFile uploads a file to MinIO.
func (c *MinIOClient) UploadFile(ctx context.Context, bucketName, objectName, filePath string) error {
	// Open the local file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Get file stat for size info
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	// Upload the file to MinIO
	_, err = c.client.PutObject(ctx, bucketName, objectName, file, fileInfo.Size(), minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("error uploading file to MinIO: %w", err)
	}

	return nil
}

// DownloadFile downloads a file from MinIO.
func (c *MinIOClient) DownloadFile(ctx context.Context, bucketName, objectName, filePath string) error {
	// Create the local file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// Download the file from MinIO
	reader, err := c.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("error downloading file from MinIO: %w", err)
	}
	defer reader.Close()

	// Copy the downloaded content to the local file
	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("error copying file content: %w", err)
	}

	return nil
}
