package s3service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3ServiceInterface defines the interface for S3 operations
type S3ServiceInterface interface {
	GenerateUploadPayload(ctx context.Context, data []byte, s3Path, filename string, metadata map[string]string) (*s3.PutObjectInput, error)
	BulkUploadFiles(ctx context.Context, payloads []*s3.PutObjectInput) ([]string, error)
	BulkDeleteFiles(ctx context.Context, uploadedKeys []string)
}
