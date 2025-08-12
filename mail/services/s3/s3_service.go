package s3service

import (
	"context"
	"os"

	"github.com/atomic-blend/backend/mail/utils/s3"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3ServiceWrapper wraps the existing S3 functionality
type S3ServiceWrapper struct {
	service *s3.Service
}

// NewS3Service creates a new S3 service wrapper
func NewS3Service() (S3ServiceInterface, error) {
	s3Service, err := s3.NewS3Service(os.Getenv("AWS_BUCKET"))
	if err != nil {
		return nil, err
	}

	return &S3ServiceWrapper{
		service: s3Service,
	}, nil
}

// GenerateUploadPayload generates a payload for uploading a file to S3
func (s *S3ServiceWrapper) GenerateUploadPayload(ctx context.Context, data []byte, s3Path, filename string, metadata map[string]string) (*awss3.PutObjectInput, error) {
	return s.service.GenerateUploadPayload(ctx, data, s3Path, filename, metadata)
}

// BulkUploadFiles uploads multiple files given a list of payloads
func (s *S3ServiceWrapper) BulkUploadFiles(ctx context.Context, payloads []*awss3.PutObjectInput) ([]string, error) {
	return s.service.BulkUploadFiles(ctx, payloads)
}

// BulkDeleteFiles deletes multiple files from S3
func (s *S3ServiceWrapper) BulkDeleteFiles(ctx context.Context, uploadedKeys []string) {
	s.service.BulkDeleteFiles(ctx, uploadedKeys)
}
