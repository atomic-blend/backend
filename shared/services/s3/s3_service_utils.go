package s3service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

// Service provides methods for interacting with AWS S3
type Service struct {
	client *s3.Client
	bucket string
}

// BulkDeleteFiles deletes multiple files from S3
func (s *Service) BulkDeleteFiles(ctx context.Context, uploadedKeys []string) {
	for _, key := range uploadedKeys {
		s.DeleteFile(ctx, nil, nil, &key)
	}
}

// NewS3Service creates a new S3 service instance
func NewS3Service() (*Service, error) {
	// Load AWS configuration
	bucket := os.Getenv("AWS_BUCKET")
	endpoint := os.Getenv("AWS_ENDPOINT")
	region := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithBaseEndpoint(endpoint), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// check if AWS_USE_PATH_STYLE_ENDPOINT is set
	usePathStyleEndpoint := false
	if os.Getenv("AWS_USE_PATH_STYLE_ENDPOINT") == "true" {
		log.Info().Msg("using path style endpoint")
		usePathStyleEndpoint = true
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = usePathStyleEndpoint
	})

	return &Service{
		client: client,
		bucket: bucket,
	}, nil
}

// GenerateUploadPayload generates a payload for uploading a file to S3
// have the same signature as UploadFile, but return the payload instead of uploading for bulk upload
func (s *Service) GenerateUploadPayload(ctx context.Context, data []byte, s3Path, filename string, metadata map[string]string) (*s3.PutObjectInput, error) {
	// Construct the full S3 key
	s3Key := filepath.Join(s3Path, filename)

	// Create a reader from the byte array
	reader := bytes.NewReader(data)

	// Prepare the upload input
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s3Key),
		Body:   reader,
	}

	// Add metadata if provided
	if metadata != nil {
		input.Metadata = metadata
	}

	return input, nil
}

// BulkUploadFiles uploads multiple files given a list of payloads.
// Auto rollback if there are errors
func (s *Service) BulkUploadFiles(ctx context.Context, payloads []*s3.PutObjectInput) ([]string, error) {
	// Upload the files
	var uploadedKeys []string
	var haveErrors bool

	for _, payload := range payloads {
		_, err := s.client.PutObject(ctx, payload)
		if err != nil {
			log.Error().Err(err).
				Str("bucket", s.bucket).
				Str("key", *payload.Key).
				Msg("failed to upload file to S3")
			haveErrors = true
			break
		}
		uploadedKeys = append(uploadedKeys, *payload.Key)
	}

	// rollback if there are errors
	if haveErrors {
		for _, key := range uploadedKeys {
			s.DeleteFile(ctx, nil, nil, &key)
		}
		return nil, fmt.Errorf("failed to upload some files to S3")
	}

	return uploadedKeys, nil
}

// UploadFile uploads a file to S3 with the given path and filename
// metadata is optional - if nil, no metadata will be added
func (s *Service) UploadFile(ctx context.Context, data []byte, s3Path, filename string, metadata map[string]string) error {
	// Construct the full S3 key
	s3Key := filepath.Join(s3Path, filename)

	// Create a reader from the byte array
	reader := bytes.NewReader(data)

	// Prepare the upload input
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s3Key),
		Body:   reader,
	}

	// Add metadata if provided
	if metadata != nil {
		input.Metadata = metadata
	}

	// Upload the file
	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		log.Error().Err(err).
			Str("bucket", s.bucket).
			Str("key", s3Key).
			Msg("failed to upload file to S3")
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Log success with appropriate metadata info
	logEntry := log.Info().
		Str("bucket", s.bucket).
		Str("key", s3Key).
		Int("size", len(data))

	if metadata != nil {
		logEntry = logEntry.Interface("metadata", metadata)
	}

	logEntry.Msg("successfully uploaded file to S3")

	return nil
}

// DeleteFile deletes a file from S3
func (s *Service) DeleteFile(ctx context.Context, s3Path, filename, s3Key *string) error {
	// Construct the full S3 key
	if s3Key == nil {
		s3Key = aws.String(filepath.Join(*s3Path, *filename))
	}

	// Prepare the delete input
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    s3Key,
	}

	log.Info().
		Str("bucket", s.bucket).
		Str("key", *s3Key).
		Msg("deleting file from S3")

	// Delete the file
	_, err := s.client.DeleteObject(ctx, input)
	if err != nil {
		log.Error().Err(err).
			Str("bucket", s.bucket).
			Str("key", *s3Key).
			Msg("failed to delete file from S3")
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	log.Info().
		Str("bucket", s.bucket).
		Str("key", *s3Key).
		Msg("successfully deleted file from S3")

	return nil
}

// FileExists checks if a file exists in S3
func (s *Service) FileExists(ctx context.Context, s3Path, filename string) (bool, error) {
	// Construct the full S3 key
	s3Key := filepath.Join(s3Path, filename)

	// Prepare the head object input
	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s3Key),
	}

	// Check if the file exists
	_, err := s.client.HeadObject(ctx, input)
	if err != nil {
		var responseError *awshttp.ResponseError
		if errors.As(err, &responseError) && responseError.ResponseError.HTTPStatusCode() == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// GeneratePreSignedDownloadURL generates a pre-signed URL for downloading a file from S3
// The URL will be valid for the specified duration (in seconds)
func (s *Service) GeneratePreSignedDownloadURL(ctx context.Context, s3Path, filename string, expirationSeconds int64) (string, error) {
	// Construct the full S3 key
	s3Key := filepath.Join(s3Path, filename)

	// Create a presign client
	presignClient := s3.NewPresignClient(s.client)

	// Prepare the presign input
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s3Key),
	}

	// Generate the presigned URL
	presignResult, err := presignClient.PresignGetObject(ctx, input, s3.WithPresignExpires(time.Duration(expirationSeconds)*time.Second))
	if err != nil {
		log.Error().Err(err).
			Str("bucket", s.bucket).
			Str("key", s3Key).
			Int64("expiration_seconds", expirationSeconds).
			Msg("failed to generate presigned download URL")
		return "", fmt.Errorf("failed to generate presigned download URL: %w", err)
	}

	log.Info().
		Str("bucket", s.bucket).
		Str("key", s3Key).
		Int64("expiration_seconds", expirationSeconds).
		Str("url", presignResult.URL).
		Msg("successfully generated presigned download URL")

	return presignResult.URL, nil
}
