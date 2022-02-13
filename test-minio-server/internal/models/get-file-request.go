package models

type GetFileRequest struct {
	BucketName string `json:"bucketName"`
	S3Filename string `json:"s3Filename"`
}
