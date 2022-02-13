package main

import (
	"context"
	"fmt"
	"test-minio-server/internal/config"
	"test-minio-server/internal/controllers"
	"test-minio-server/internal/s3"

	"github.com/labstack/echo/v4"
)

func main() {
	ctx := context.Background()
	e := echo.New()
	e.Debug = true

	// Cnnfig
	config, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	// StoreController
	storageService := s3.NewService(ctx, e.Logger, &config.Minio)
	storageController := controllers.NewStorageController(storageService)
	e.POST("/storeFile", storageController.StoreFile)
	e.POST("/getFile", storageController.GetFile)

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:1323")))
}
