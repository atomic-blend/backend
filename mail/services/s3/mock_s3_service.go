// Package s3service defines s3 operations available
package s3service

import (
	"context"

	s3interfaces "github.com/atomic-blend/backend/mail/services/s3/interfaces"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/mock"
)

// MockS3Service provides a mock implementation of S3Service
type MockS3Service struct {
	mock.Mock
}

// Ensure MockS3Service implements the interface
var _ s3interfaces.S3ServiceInterface = (*MockS3Service)(nil)

// BulkDeleteFiles deletes multiple files from S3
func (m *MockS3Service) BulkDeleteFiles(ctx context.Context, uploadedKeys []string) {
	m.Called(ctx, uploadedKeys)
}

// GenerateUploadPayload generates a payload for uploading a file to S3
func (m *MockS3Service) GenerateUploadPayload(ctx context.Context, data []byte, s3Path, filename string, metadata map[string]string) (*s3.PutObjectInput, error) {
	args := m.Called(ctx, data, s3Path, filename, metadata)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.PutObjectInput), args.Error(1)
}

// BulkUploadFiles uploads multiple files given a list of payloads
func (m *MockS3Service) BulkUploadFiles(ctx context.Context, payloads []*s3.PutObjectInput) ([]string, error) {
	args := m.Called(ctx, payloads)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// UploadFile uploads a file to S3 with the given path and filename
func (m *MockS3Service) UploadFile(ctx context.Context, data []byte, s3Path, filename string, metadata map[string]string) error {
	args := m.Called(ctx, data, s3Path, filename, metadata)
	return args.Error(0)
}

// DeleteFile deletes a file from S3
func (m *MockS3Service) DeleteFile(ctx context.Context, s3Path, filename, s3Key *string) error {
	args := m.Called(ctx, s3Path, filename, s3Key)
	return args.Error(0)
}

// FileExists checks if a file exists in S3
func (m *MockS3Service) FileExists(ctx context.Context, s3Path, filename string) (bool, error) {
	args := m.Called(ctx, s3Path, filename)
	return args.Bool(0), args.Error(1)
}

// GeneratePreSignedDownloadURL generates a presigned download URL
func (m *MockS3Service) GeneratePreSignedDownloadURL(ctx context.Context, s3Path, filename string, expirationSeconds int64) (string, error) {
	args := m.Called(ctx, s3Path, filename, expirationSeconds)
	return args.String(0), args.Error(1)
}
