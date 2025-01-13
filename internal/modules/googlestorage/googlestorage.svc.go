package googlestorage

import (
	"app/config"
	"app/internal/modules/log"
	"context"
	"os"

	"cloud.google.com/go/storage"
	"github.com/joho/godotenv"
)

type GoogleStorageService struct {
}

func newService(conf *config.Config) *GoogleStorageService {
	godotenv.Load()

	return &GoogleStorageService{}
}

func initCloudStorage() {

	log.Info("GCS_BUCKET_NAME: %v", os.Getenv("GCS_BUCKET_NAME"))

}

func StorageConfig(ctx context.Context) *storage.Client {

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Error("Failed to create client: %v", err)
		panic(err)
	}

	return client
}

// func StorageConfig() []option.ClientOption {
// 	// ดึง credentials JSON จาก environment variable ผ่าน viper
// 	credentialsJSON := os.Getenv("GCS_CREDENTIALS_JSON")
// 	if credentialsJSON == "" {
// 		log.Info("GCS_CREDENTIALS_JSON environment variable is not set")
// 		return nil
// 	}

// 	// ส่งข้อมูล credentialsJSON เป็น byte array เข้าไปใน option.WithCredentialsJSON
// 	return []option.ClientOption{option.WithCredentialsJSON([]byte(credentialsJSON))}
// }
