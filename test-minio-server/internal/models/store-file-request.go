package models

type StoreFileRequest struct {
	BucketName string `json:"bucketName"`
	Location   string `json:"location"`
	Data       string `json:"data"`
	S3Filename string `json:"s3Filename"`
}
