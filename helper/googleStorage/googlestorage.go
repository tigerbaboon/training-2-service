package helper

import (
	"app/internal/modules/googlestorage"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

var (
	storageClient *storage.Client
)

func UploadFileGCS(file *multipart.FileHeader) (*string, error) {
	ctx := context.Background()
	bucket := os.Getenv("GCS_BUCKET_NAME")

	blobFile, err := file.Open()
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext
	filePath := "group/" + filename

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	client := googlestorage.StorageConfig(ctx)

	obj := client.Bucket(bucket).Object(filePath)

	writer := obj.NewWriter(ctx)

	if _, err := io.Copy(writer, blobFile); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	acl := obj.ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return nil, err
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, filePath)

	return &publicURL, nil
}

// func UploadFileGCS(file *multipart.FileHeader) (*string, error) {

// 	ctx := context.Background()
// 	bucket := os.Getenv("GCS_BUCKET_NAME")
// 	var err error
// 	storageClient, err = storage.NewClient(ctx, googlestorage.StorageConfig()...)
// 	if err != nil {
// 		fmt.Printf("Failed to create GCS client: %v\n", err)
// 		return nil, err
// 	}

// 	blobFile, err := file.Open()
// 	if err != nil {
// 		fmt.Printf("Failed to open file: %v\n", err)
// 		return nil, err
// 	}
// 	defer blobFile.Close() // Ensure file is closed after use

// 	ext := filepath.Ext(file.Filename)
// 	filename := uuid.New().String() + ext
// 	filePath := "group/" + filename

// 	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
// 	defer cancel()

// 	obj := storageClient.Bucket(bucket).Object(filePath)
// 	writer := obj.NewWriter(ctx)

// 	if _, err := io.Copy(writer, blobFile); err != nil {
// 		fmt.Printf("Failed to upload file to GCS: %v\n", err)
// 		return nil, err
// 	}
// 	if err := writer.Close(); err != nil {
// 		fmt.Printf("Failed to close GCS writer: %v\n", err)
// 		return nil, err
// 	}

// 	acl := obj.ACL()
// 	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
// 		fmt.Printf("Failed to set ACL for public access: %v\n", err)
// 		return nil, err
// 	}

// 	// Construct the public URL and return it
// 	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, filePath)
// 	fmt.Printf("Upload successful, public URL: %s\n", publicURL)

// 	return &publicURL, nil
// }
