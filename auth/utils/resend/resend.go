package resend

import (
	"github.com/resend/resend-go/v2"
)

// EmailClient interface defines the methods needed from the Resend client
type EmailClient interface {
	Send(to []string, subject string, html string, text string) (string, error)
}

// Client implements EmailClient using the actual Resend API
type Client struct {
	client *resend.Client
}

// NewResendClient creates a new ResendClient
func NewResendClient(apiKey string) *Client {
	return &Client{
		client: resend.NewClient(apiKey),
	}
}

// Send sends an email using the Resend API
func (c *Client) Send(to []string, subject string, html string, text string) (string, error) {
	params := &resend.SendEmailRequest{
		From:    "noreply@brandonguigo.com",
		To:      to,
		Subject: subject,
		Html:    html,
		Text:    text,
	}

	sent, err := c.client.Emails.Send(params)
	if err != nil {
		return "", err
	}
	return sent.Id, nil
}
