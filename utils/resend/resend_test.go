package resend

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockEmailClient implements the EmailClient interface for testing
type MockEmailClient struct {
	shouldFail bool
}

// Send mocks sending an email
func (m *MockEmailClient) Send(to []string, subject string, html string, text string) (string, error) {
	if m.shouldFail {
		return "", errors.New("mock send error")
	}
	return "mock-email-id", nil
}

// Create a mock factory for testing
func mockFactory(client EmailClient) ClientFactory {
	return func(apiKey string) EmailClient {
		return client
	}
}

func TestSendEmailRequest(t *testing.T) {
	// Save original API key to restore after tests
	originalAPIKey := os.Getenv("RESEND_API_KEY")
	defer os.Setenv("RESEND_API_KEY", originalAPIKey)

	tests := []struct {
		name        string
		to          []string
		subject     string
		html        string
		text        string
		mockClient  EmailClient
		expectError bool
		expectSent  bool
	}{
		{
			name:        "Successful email send",
			to:          []string{"test@example.com"},
			subject:     "Test Subject",
			html:        "<p>Test HTML</p>",
			text:        "Test plain text",
			mockClient:  &MockEmailClient{shouldFail: false},
			expectError: false,
			expectSent:  true,
		},
		{
			name:        "Failed email send",
			to:          []string{"test@example.com"},
			subject:     "Test Subject",
			html:        "<p>Test HTML</p>",
			text:        "Test plain text",
			mockClient:  &MockEmailClient{shouldFail: true},
			expectError: true,
			expectSent:  false,
		},
		{
			name:        "Empty recipients",
			to:          []string{},
			subject:     "Test Subject",
			html:        "<p>Test HTML</p>",
			text:        "Test plain text",
			mockClient:  &MockEmailClient{shouldFail: false},
			expectError: true,
			expectSent:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set a fake API key
			os.Setenv("RESEND_API_KEY", "test-api-key")

			// Create a factory that returns our mock client
			factory := mockFactory(tt.mockClient)

			// Call the function with our mock factory
			sent, err := SendEmailRequestWithFactory(tt.to, tt.subject, tt.html, tt.text, factory)

			// Check results
			if tt.expectError {
				assert.Error(t, err)
				assert.False(t, sent)
			} else {
				assert.NoError(t, err)
				assert.True(t, sent)
			}
		})
	}

	// Test missing API key scenario
	t.Run("Missing API Key", func(t *testing.T) {
		// Clear the API key
		os.Setenv("RESEND_API_KEY", "")

		// Use a mock factory that would succeed if called (but it shouldn't be called)
		factory := mockFactory(&MockEmailClient{shouldFail: false})

		// Call the function
		sent, err := SendEmailRequestWithFactory([]string{"test@example.com"}, "Test Subject", "<p>Test HTML</p>", "Test plain text", factory)

		// Should fail with a proper error message
		assert.False(t, sent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "RESEND_API_KEY is not set")
	})
}
