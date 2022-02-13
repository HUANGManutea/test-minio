package controllers

import (
	"net/http"
	"os"
	"path"
	"test-minio-server/internal/models"
	"test-minio-server/internal/s3"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
)

type StorageController struct {
	Service *s3.Service
}

func NewStorageController(service *s3.Service) *StorageController {
	return &StorageController{
		Service: service,
	}
}

func (controller *StorageController) StoreFile(c echo.Context) error {
	// retrieve data
	storeFileRequest := models.StoreFileRequest{}

	if err := c.Bind(&storeFileRequest); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusBadRequest)
	}

	// create temp file
	tempFileName := uuid.New().String()

	tempPath := path.Join("/tmp", tempFileName)
	dst, err := os.Create(tempPath)
	if err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	_, err = dst.Write([]byte(storeFileRequest.Data))
	if err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	dst.Close()

	err = controller.Service.StoreFile(storeFileRequest, tempPath)
	if err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	err = os.Remove(tempPath)
	if err != nil {
		c.Logger().Error("failed to delete temp file", "error", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	storeFileResponse := models.StoreFileResponse{
		S3Filename: storeFileRequest.S3Filename,
	}

	return c.JSON(http.StatusOK, storeFileResponse)
}

func (controller *StorageController) GetFile(c echo.Context) error {
	// retrieve data
	getFileRequest := models.GetFileRequest{}
	c.Logger().Infof("getFileRequest: %v", getFileRequest)

	if err := c.Bind(&getFileRequest); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusBadRequest)
	}

	data, err := controller.Service.GetFile(getFileRequest)
	if err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Blob(http.StatusOK, "text/plain", data)
}
