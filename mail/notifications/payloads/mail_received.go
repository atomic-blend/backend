package payloads

// MailReceivedPayload represents the payload for mail received notifications.
type MailReceivedPayload struct {
	Type           string `json:"type"`
	MailFrom       string `json:"mail_from"`
	Subject        string `json:"subject"`
	ContentPreview string `json:"content_preview"`
}

// NewMailReceivedPayload creates a new MailReceivedPayload with the given from, subject and content preview.
func NewMailReceivedPayload(from string, subject string, contentPreview string) *MailReceivedPayload {
	return &MailReceivedPayload{
		Type:           "MAIL_RECEIVED",
		MailFrom:       from,
		Subject:        subject,
		ContentPreview: contentPreview,
	}
}

// GetType returns the type of the payload.
func (p *MailReceivedPayload) GetType() string {
	return p.Type
}

// GetData returns the ready to send data for the payload.
func (p *MailReceivedPayload) GetData() map[string]string {
	return map[string]string{
		"type":            p.Type,
		"mail_from":       p.MailFrom,
		"subject":         p.Subject,
		"content_preview": p.ContentPreview,
	}
}
