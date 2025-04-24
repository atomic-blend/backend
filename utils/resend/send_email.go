package resend

import (
	"errors"
	"os"
)

// ClientFactory function type for creating email clients
type ClientFactory func(apiKey string) EmailClient

// DefaultClientFactory that creates real ResendClient instances
var DefaultClientFactory ClientFactory = func(apiKey string) EmailClient {
	return NewResendClient(apiKey)
}

// SendEmailRequest sends an email using the Resend API
func SendEmailRequest(to []string, subject string, html string, text string) (bool, error) {
	return SendEmailRequestWithFactory(to, subject, html, text, DefaultClientFactory)
}

// SendEmailRequestWithFactory sends an email using a provided client factory
// This makes testing easier by allowing dependency injection
func SendEmailRequestWithFactory(to []string, subject string, html string, text string, factory ClientFactory) (bool, error) {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return false, errors.New("RESEND_API_KEY is not set")
	}

	if len(to) == 0 {
		return false, errors.New("recipients list is empty")
	}

	client := factory(apiKey)

	_, err := client.Send(to, subject, html, text)
	if err != nil {
		return false, err
	}

	return true, nil
}
