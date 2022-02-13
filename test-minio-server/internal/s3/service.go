package s3

import (
	"context"
	"log"

	"test-minio-server/internal/config"
	"test-minio-server/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Service struct {
	Ctx    context.Context
	Logger echo.Logger
	Client *minio.Client
	Minio  *config.MinioConfig
}

func NewService(ctx context.Context, logger echo.Logger, minioConfig *config.MinioConfig) *Service {

	minioClient, err := minio.New(minioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.AccessKeyID, minioConfig.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		logger.Error(err)
	}

	return &Service{
		Ctx:    ctx,
		Logger: logger,
		Client: minioClient,
		Minio:  minioConfig,
	}
}

func (service *Service) StoreFile(storeFileRequest models.StoreFileRequest, localFileAbsolutePath string) error {

	err := service.Client.MakeBucket(service.Ctx, storeFileRequest.BucketName, minio.MakeBucketOptions{Region: storeFileRequest.Location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := service.Client.BucketExists(service.Ctx, storeFileRequest.BucketName)
		if errBucketExists == nil && exists {
			service.Logger.Errorf("We already own %s\n", storeFileRequest.BucketName)
			return nil
		} else {
			service.Logger.Error(err)
			return err
		}
	}

	service.Logger.Infof("Successfully created %s\n", storeFileRequest.BucketName)

	// Upload the file
	objectName := storeFileRequest.S3Filename
	filePath := localFileAbsolutePath
	contentType := "text/plain"

	// Upload the file with FPutObject
	info, err := service.Client.FPutObject(service.Ctx, storeFileRequest.BucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		service.Logger.Error(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	return nil
}

func (service *Service) GetFile(getFileRequest models.GetFileRequest) ([]byte, error) {
	object, err := service.Client.GetObject(service.Ctx, getFileRequest.BucketName, getFileRequest.S3Filename, minio.GetObjectOptions{})
	if err != nil {
		service.Logger.Error(err)
		return nil, err
	}
	info, err := object.Stat()
	if err != nil {
		service.Logger.Error(err)
		return nil, err
	}

	data := make([]byte, info.Size)

	object.Read(data)

	return data, nil
}
