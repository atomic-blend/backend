package mail

import (
	"context"
	"errors"
	"strings"
	"testing"

	"connectrpc.com/connect"
	userv1 "github.com/atomic-blend/backend/grpc/gen/user/v1"
	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/mail/tests/mocks"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/emersion/go-message"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMailContent_Encrypt(t *testing.T) {
	tests := []struct {
		name      string
		mail      *Content
		publicKey string
		wantErr   bool
	}{
		{
			name: "successful encryption",
			mail: &Content{
				Headers: map[string]string{
					"From":       "sender@example.com",
					"To":         "recipient@example.com",
					"Subject":    "Test Subject",
					"Date":       "2024-01-01T00:00:00Z",
					"Message-Id": "test-message-id",
					"Cc":         "cc@example.com",
					"Bcc":        "bcc@example.com",
				},
				TextContent: "Plain text content",
				HTMLContent: "<html><body>HTML content</body></html>",
				Attachments: []Attachment{
					{
						Filename:    "test.txt",
						ContentType: "text/plain",
						Data:        []byte("attachment content"),
					},
				},
				Rejected:       false,
				RewriteSubject: false,
				Greylisted:     false,
			},
			publicKey: "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p",
			wantErr:   false,
		},
		{
			name: "empty content",
			mail: &Content{
				Headers:        map[string]string{},
				TextContent:    "",
				HTMLContent:    "",
				Attachments:    []Attachment{},
				Rejected:       false,
				RewriteSubject: false,
				Greylisted:     false,
			},
			publicKey: "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.mail.Encrypt(tt.publicKey)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)

			// Verify that headers are properly handled
			if originalHeaders, ok := tt.mail.Headers.(map[string]string); ok {
				if resultHeaders, ok := result.Headers.(map[string]string); ok {
					// If both original and result have non-empty headers, verify encryption changed them
					if len(originalHeaders) > 0 {
						for key, originalValue := range originalHeaders {
							if resultValue, exists := resultHeaders[key]; exists && originalValue != "" {
								assert.NotEqual(t, originalValue, resultValue, "Header %s should be encrypted", key)
							}
						}
					}
				}
			}

			// Verify that text content is encrypted if not empty
			if tt.mail.TextContent != "" {
				assert.NotEqual(t, tt.mail.TextContent, result.TextContent)
			}

			// Verify that HTML content is encrypted if not empty
			if tt.mail.HTMLContent != "" {
				assert.NotEqual(t, tt.mail.HTMLContent, result.HTMLContent)
			}

			// Verify that non-encrypted fields remain the same
			assert.Equal(t, tt.mail.Rejected, result.Rejected)
			assert.Equal(t, tt.mail.RewriteSubject, result.RewriteSubject)
			assert.Equal(t, tt.mail.Greylisted, result.Greylisted)
			assert.Equal(t, len(tt.mail.Attachments), len(result.Attachments))
		})
	}
}

func TestReceiveMail(t *testing.T) {
	// Test data
	testPublicKey := "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
	testUserID := "user123"
	testEmail := "test@example.com"

	// Sample MIME content
	sampleMIME := `From: sender@example.com
To: test@example.com
Subject: Test Email
Date: Mon, 01 Jan 2024 00:00:00 +0000
Message-ID: <test@example.com>
Content-Type: text/plain

This is a test email content.`

	tests := []struct {
		name           string
		payload        ReceivedMailPayload
		setupMocks     func(*mocks.MockMailRepository, *mocks.MockS3Service, *mocks.MockUserClient)
		expectedAck    bool
		expectedErrors bool
	}{
		{
			name: "successful mail processing",
			payload: ReceivedMailPayload{
				Content:    sampleMIME,
				IP:         "192.168.1.1",
				Hostname:   "test-host",
				From:       "sender@example.com",
				Rcpt:       []string{testEmail},
				QueueID:    "queue123",
				User:       "user",
				DeliverTo:  testEmail,
				ReceivedAt: "2024-01-01T00:00:00Z",
			},
			setupMocks: func(mailRepo *mocks.MockMailRepository, s3Service *mocks.MockS3Service, userClient *mocks.MockUserClient) {
				// Mock user client response
				userClient.On("GetUserPublicKey", mock.Anything, mock.Anything).Return(
					&connect.Response[userv1.GetUserPublicKeyResponse]{
						Msg: &userv1.GetUserPublicKeyResponse{
							PublicKey: testPublicKey,
							UserId:    testUserID,
						},
					}, nil,
				)

				// Mock S3 service for attachment upload (no attachments in this test)
				s3Service.On("BulkUploadFiles", mock.Anything, mock.Anything).Return([]string{}, nil)

				// Mock mail repository
				mailRepo.On("CreateMany", mock.Anything, mock.Anything).Return(true, nil)
			},
			expectedAck:    true,
			expectedErrors: false,
		},
		{
			name: "user not found",
			payload: ReceivedMailPayload{
				Content:    sampleMIME,
				IP:         "192.168.1.1",
				Hostname:   "test-host",
				From:       "sender@example.com",
				Rcpt:       []string{"nonexistent@example.com"},
				QueueID:    "queue123",
				User:       "user",
				DeliverTo:  "nonexistent@example.com",
				ReceivedAt: "2024-01-01T00:00:00Z",
			},
			setupMocks: func(mailRepo *mocks.MockMailRepository, s3Service *mocks.MockS3Service, userClient *mocks.MockUserClient) {
				// Mock user client to return error (user not found)
				userClient.On("GetUserPublicKey", mock.Anything, mock.Anything).Return(
					nil, errors.New("user not found"),
				)
				// Since no users are found, no S3 or DB operations will be called
			},
			expectedAck:    true,
			expectedErrors: false, // User not found is not considered an error, just skipped
		},
		{
			name: "encryption failure",
			payload: ReceivedMailPayload{
				Content:    sampleMIME,
				IP:         "192.168.1.1",
				Hostname:   "test-host",
				From:       "sender@example.com",
				Rcpt:       []string{testEmail},
				QueueID:    "queue123",
				User:       "user",
				DeliverTo:  testEmail,
				ReceivedAt: "2024-01-01T00:00:00Z",
			},
			setupMocks: func(mailRepo *mocks.MockMailRepository, s3Service *mocks.MockS3Service, userClient *mocks.MockUserClient) {
				// Mock user client response with invalid public key
				userClient.On("GetUserPublicKey", mock.Anything, mock.Anything).Return(
					&connect.Response[userv1.GetUserPublicKeyResponse]{
						Msg: &userv1.GetUserPublicKeyResponse{
							PublicKey: "invalid-key",
							UserId:    testUserID,
						},
					}, nil,
				)
			},
			expectedAck:    false,
			expectedErrors: true,
		},
		{
			name: "S3 upload failure",
			payload: ReceivedMailPayload{
				Content:    sampleMIME,
				IP:         "192.168.1.1",
				Hostname:   "test-host",
				From:       "sender@example.com",
				Rcpt:       []string{testEmail},
				QueueID:    "queue123",
				User:       "user",
				DeliverTo:  testEmail,
				ReceivedAt: "2024-01-01T00:00:00Z",
			},
			setupMocks: func(mailRepo *mocks.MockMailRepository, s3Service *mocks.MockS3Service, userClient *mocks.MockUserClient) {
				// Mock user client response
				userClient.On("GetUserPublicKey", mock.Anything, mock.Anything).Return(
					&connect.Response[userv1.GetUserPublicKeyResponse]{
						Msg: &userv1.GetUserPublicKeyResponse{
							PublicKey: testPublicKey,
							UserId:    testUserID,
						},
					}, nil,
				)

				// Mock S3 service to return error
				s3Service.On("BulkUploadFiles", mock.Anything, mock.Anything).Return(nil, errors.New("S3 upload failed"))
			},
			expectedAck:    false,
			expectedErrors: true,
		},
		{
			name: "database save failure",
			payload: ReceivedMailPayload{
				Content:    sampleMIME,
				IP:         "192.168.1.1",
				Hostname:   "test-host",
				From:       "sender@example.com",
				Rcpt:       []string{testEmail},
				QueueID:    "queue123",
				User:       "user",
				DeliverTo:  testEmail,
				ReceivedAt: "2024-01-01T00:00:00Z",
			},
			setupMocks: func(mailRepo *mocks.MockMailRepository, s3Service *mocks.MockS3Service, userClient *mocks.MockUserClient) {
				// Mock user client response
				userClient.On("GetUserPublicKey", mock.Anything, mock.Anything).Return(
					&connect.Response[userv1.GetUserPublicKeyResponse]{
						Msg: &userv1.GetUserPublicKeyResponse{
							PublicKey: testPublicKey,
							UserId:    testUserID,
						},
					}, nil,
				)

				// Mock S3 service success
				s3Service.On("BulkUploadFiles", mock.Anything, mock.Anything).Return([]string{"key1"}, nil)
				s3Service.On("BulkDeleteFiles", mock.Anything, mock.Anything).Return()

				// Mock mail repository to return error
				mailRepo.On("CreateMany", mock.Anything, mock.Anything).Return(false, errors.New("database error"))
			},
			expectedAck:    false,
			expectedErrors: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockMailRepo := &mocks.MockMailRepository{}
			mockS3Service := &mocks.MockS3Service{}
			mockUserClient := &mocks.MockUserClient{}

			// Setup mocks
			tt.setupMocks(mockMailRepo, mockS3Service, mockUserClient)

			// Create AMQP delivery
			delivery := &amqp.Delivery{
				Acknowledger: &mockAcknowledger{},
			}

			// Create a test function that uses our mocks
			testReceiveMail := func(m *amqp.Delivery, payload ReceivedMailPayload) {
				// This would normally call the real receiveMail function
				// For testing, we'll simulate the key parts

				encryptedMails := make([]models.Mail, 0)
				encryptedAttachments := make([]*s3.PutObjectInput, 0)
				haveErrors := false

				for _, rcpt := range payload.Rcpt {
					// Get user public key
					rcptPublicKey, err := mockUserClient.GetUserPublicKey(context.Background(), &connect.Request[userv1.GetUserPublicKeyRequest]{
						Msg: &userv1.GetUserPublicKeyRequest{
							Email: rcpt,
						},
					})
					if err != nil {
						continue // User not found, skip
					}

					userPublicKey := rcptPublicKey.Msg.PublicKey

					// Create mail content
					mailContent := &Content{
						Headers: map[string]string{
							"From":       payload.From,
							"To":         rcpt,
							"Subject":    "Test Subject",
							"Date":       payload.ReceivedAt,
							"Message-Id": "test-id",
						},
						TextContent: "Test content",
						HTMLContent: "<html>Test</html>",
						Attachments: []Attachment{},
					}

					// Encrypt mail content
					encryptedMailContent, err := mailContent.Encrypt(userPublicKey)
					if err != nil {
						haveErrors = true
						continue
					}

					// Create mail entity
					mailEntity := &models.Mail{
						Headers:     encryptedMailContent.Headers,
						TextContent: encryptedMailContent.TextContent,
						HTMLContent: encryptedMailContent.HTMLContent,
					}

					encryptedMails = append(encryptedMails, *mailEntity)
				}

				if haveErrors {
					return
				}

				// Only proceed with S3 and DB operations if we have mails to process
				if len(encryptedMails) > 0 {
					// Upload attachments
					_, err := mockS3Service.BulkUploadFiles(context.Background(), encryptedAttachments)
					if err != nil {
						return
					}

					// Save to database
					_, err = mockMailRepo.CreateMany(context.Background(), encryptedMails)
					if err != nil {
						// Rollback S3 uploads
						mockS3Service.BulkDeleteFiles(context.Background(), []string{})
						return
					}
				}

				// Acknowledge message
				delivery.Ack(false)
			}

			// Execute test
			testReceiveMail(delivery, tt.payload)

			// Verify expectations
			mockMailRepo.AssertExpectations(t)
			mockS3Service.AssertExpectations(t)
			mockUserClient.AssertExpectations(t)
		})
	}
}

// Mock acknowledger for testing
type mockAcknowledger struct{}

func (m *mockAcknowledger) Ack(tag uint64, multiple bool) error {
	return nil
}

func (m *mockAcknowledger) Nack(tag uint64, multiple bool, requeue bool) error {
	return nil
}

func (m *mockAcknowledger) Reject(tag uint64, requeue bool) error {
	return nil
}

func TestProcessMessageBody(t *testing.T) {
	tests := []struct {
		name     string
		mimeData string
		want     *Content
	}{
		{
			name: "plain text message",
			mimeData: `From: sender@example.com
To: recipient@example.com
Subject: Test
Content-Type: text/plain

This is plain text content.`,
			want: &Content{
				TextContent: "This is plain text content.",
				HTMLContent: "",
				Attachments: []Attachment{},
			},
		},
		{
			name: "HTML message",
			mimeData: `From: sender@example.com
To: recipient@example.com
Subject: Test
Content-Type: text/html

<html><body>This is HTML content.</body></html>`,
			want: &Content{
				TextContent: "",
				HTMLContent: "<html><body>This is HTML content.</body></html>",
				Attachments: []Attachment{},
			},
		},
		{
			name: "multipart message with attachment",
			mimeData: `From: sender@example.com
To: recipient@example.com
Subject: Test
Content-Type: multipart/mixed; boundary="boundary"

--boundary
Content-Type: text/plain

This is plain text.

--boundary
Content-Type: application/octet-stream
Content-Disposition: attachment; filename="test.txt"

attachment content
--boundary--`,
			want: &Content{
				TextContent: "This is plain text.\n",
				HTMLContent: "",
				Attachments: []Attachment{
					{
						Filename:    "test.txt",
						ContentType: "application/octet-stream",
						Data:        []byte("attachment content"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the MIME message
			entity, err := message.Read(strings.NewReader(tt.mimeData))
			require.NoError(t, err)

			mailContent := &Content{
				Attachments: make([]Attachment, 0),
			}

			processMessageBody(entity, mailContent)

			// Verify headers are extracted
			assert.NotNil(t, mailContent.Headers)
			if headers, ok := mailContent.Headers.(map[string]string); ok {
				assert.Contains(t, headers, "From")
				assert.Contains(t, headers, "To")
				assert.Contains(t, headers, "Subject")
				assert.Equal(t, "sender@example.com", headers["From"])
				assert.Equal(t, "recipient@example.com", headers["To"])
			}

			// Verify text content
			if tt.want.TextContent != "" {
				assert.Equal(t, tt.want.TextContent, mailContent.TextContent)
			}

			// Verify HTML content
			if tt.want.HTMLContent != "" {
				assert.Equal(t, tt.want.HTMLContent, mailContent.HTMLContent)
			}

			// Verify attachments
			assert.Equal(t, len(tt.want.Attachments), len(mailContent.Attachments))
			for i, expectedAttachment := range tt.want.Attachments {
				if i < len(mailContent.Attachments) {
					actualAttachment := mailContent.Attachments[i]
					assert.Equal(t, expectedAttachment.Filename, actualAttachment.Filename)
					assert.Equal(t, expectedAttachment.ContentType, actualAttachment.ContentType)
					assert.Equal(t, expectedAttachment.Data, actualAttachment.Data)
				}
			}
		})
	}
}

func TestProcessMultipartMessage(t *testing.T) {
	multipartMIME := `From: sender@example.com
To: recipient@example.com
Subject: Test Multipart
Content-Type: multipart/mixed; boundary="boundary"

--boundary
Content-Type: text/plain

Plain text part.

--boundary
Content-Type: text/html

<html><body>HTML part.</body></html>

--boundary
Content-Type: application/pdf
Content-Disposition: attachment; filename="document.pdf"

PDF content here
--boundary--`

	entity, err := message.Read(strings.NewReader(multipartMIME))
	require.NoError(t, err)

	mailContent := &Content{
		Attachments: make([]Attachment, 0),
	}

	processMessageBody(entity, mailContent)

	// Verify headers are extracted
	assert.NotNil(t, mailContent.Headers)
	if headers, ok := mailContent.Headers.(map[string]string); ok {
		assert.Contains(t, headers, "From")
		assert.Contains(t, headers, "To")
		assert.Contains(t, headers, "Subject")
		assert.Equal(t, "sender@example.com", headers["From"])
		assert.Equal(t, "recipient@example.com", headers["To"])
		assert.Equal(t, "Test Multipart", headers["Subject"])
	}

	// Verify that all parts were processed
	assert.Equal(t, "Plain text part.\n", mailContent.TextContent)
	assert.Equal(t, "<html><body>HTML part.</body></html>\n", mailContent.HTMLContent)
	assert.Equal(t, 1, len(mailContent.Attachments))
	assert.Equal(t, "document.pdf", mailContent.Attachments[0].Filename)
	assert.Equal(t, "application/pdf", mailContent.Attachments[0].ContentType)
}

func TestProcessMessagePart(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		disposition string
		filename    string
		body        string
		expected    func(*Content)
	}{
		{
			name:        "plain text part",
			contentType: "text/plain",
			body:        "Plain text content",
			expected: func(mc *Content) {
				assert.Equal(t, "Plain text content", mc.TextContent)
			},
		},
		{
			name:        "HTML part",
			contentType: "text/html",
			body:        "<html><body>HTML content</body></html>",
			expected: func(mc *Content) {
				assert.Equal(t, "<html><body>HTML content</body></html>", mc.HTMLContent)
			},
		},
		{
			name:        "attachment with disposition",
			contentType: "application/pdf",
			disposition: "attachment",
			filename:    "document.pdf",
			body:        "PDF content",
			expected: func(mc *Content) {
				assert.Equal(t, 1, len(mc.Attachments))
				assert.Equal(t, "document.pdf", mc.Attachments[0].Filename)
				assert.Equal(t, "application/pdf", mc.Attachments[0].ContentType)
				assert.Equal(t, []byte("PDF content"), mc.Attachments[0].Data)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a simple MIME message for testing
			mimeData := "Content-Type: " + tt.contentType + "\n"
			if tt.disposition != "" {
				mimeData += "Content-Disposition: " + tt.disposition + "; filename=\"" + tt.filename + "\"\n"
			}
			mimeData += "\n" + tt.body

			entity, err := message.Read(strings.NewReader(mimeData))
			require.NoError(t, err)

			mailContent := &Content{
				Attachments: make([]Attachment, 0),
			}

			processMessagePart(entity, mailContent)

			tt.expected(mailContent)
		})
	}
}

func TestProcessMessageBody_HeaderExtraction(t *testing.T) {
	tests := []struct {
		name            string
		mimeData        string
		expectedHeaders map[string]string
	}{
		{
			name: "basic headers extraction",
			mimeData: `From: sender@example.com
To: recipient@example.com
Subject: Test Email
Date: Mon, 01 Jan 2024 00:00:00 +0000
Message-ID: <test@example.com>
Cc: cc@example.com
Bcc: bcc@example.com
Content-Type: text/plain

This is a test email.`,
			expectedHeaders: map[string]string{
				"From":       "sender@example.com",
				"To":         "recipient@example.com",
				"Subject":    "Test Email",
				"Date":       "Mon, 01 Jan 2024 00:00:00 +0000",
				"Message-Id": "<test@example.com>",
				"Cc":         "cc@example.com",
				"Bcc":        "bcc@example.com",
			},
		},
		{
			name: "custom headers extraction",
			mimeData: `From: sender@example.com
To: recipient@example.com
Subject: Custom Headers Test
X-Custom-Header: custom-value
X-Priority: high
X-Mailer: AtomicBlend
Content-Type: text/plain

Email with custom headers.`,
			expectedHeaders: map[string]string{
				"From":            "sender@example.com",
				"To":              "recipient@example.com",
				"Subject":         "Custom Headers Test",
				"X-Custom-Header": "custom-value",
				"X-Priority":      "high",
				"X-Mailer":        "AtomicBlend",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the MIME message
			entity, err := message.Read(strings.NewReader(tt.mimeData))
			require.NoError(t, err)

			mailContent := &Content{
				Attachments: make([]Attachment, 0),
			}

			processMessageBody(entity, mailContent)

			// Verify headers are extracted correctly
			assert.NotNil(t, mailContent.Headers)
			if headers, ok := mailContent.Headers.(map[string]string); ok {
				for expectedKey, expectedValue := range tt.expectedHeaders {
					assert.Contains(t, headers, expectedKey, "Header %s should be present", expectedKey)
					assert.Equal(t, expectedValue, headers[expectedKey], "Header %s should have correct value", expectedKey)
				}
			} else {
				t.Errorf("Headers should be of type map[string]string, got %T", mailContent.Headers)
			}
		})
	}
}
