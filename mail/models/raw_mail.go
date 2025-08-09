package models

import (
	"github.com/atomic-blend/backend/mail/utils/age_encryption"
)

// RawMail represents the collected content from an email
type RawMail struct {
	Headers        interface{}
	TextContent    string
	HTMLContent    string
	Attachments    []RawAttachment
	Rejected       bool
	RewriteSubject bool
	Greylisted     bool
}

// Attachment represents a file attachment
type RawAttachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

// Encrypt encrypts the content using the age encryption library
func (m *RawMail) Encrypt(publicKey string) (*RawMail, error) {
	encryptedMail := &RawMail{
		Attachments:    make([]RawAttachment, 0),
		Rejected:       m.Rejected,
		RewriteSubject: m.RewriteSubject,
		Greylisted:     m.Greylisted,
	}

	// encrypt all headers
	if m.Headers != nil {
		if headersMap, ok := m.Headers.(map[string]string); ok {
			encryptedHeaders := make(map[string]string)
			for key, value := range headersMap {
				encryptedValue, err := ageencryption.EncryptString(publicKey, value)
				if err != nil {
					return nil, err
				}
				encryptedHeaders[key] = encryptedValue
			}
			encryptedMail.Headers = encryptedHeaders
		} else {
			// If Headers is not a map[string]string, leave it as is (no encryption)
			encryptedMail.Headers = m.Headers
		}
	}

	encryptedTextContent, err := ageencryption.EncryptString(publicKey, m.TextContent)
	if err != nil {
		return nil, err
	}
	encryptedMail.TextContent = encryptedTextContent

	encryptedHTMLContent, err := ageencryption.EncryptString(publicKey, m.HTMLContent)
	if err != nil {
		return nil, err
	}
	encryptedMail.HTMLContent = encryptedHTMLContent

	for _, attachment := range m.Attachments {
		encryptedAttachment, err := ageencryption.EncryptBytes(publicKey, attachment.Data)
		if err != nil {
			return nil, err
		}
		encryptedFilename, err := ageencryption.EncryptString(publicKey, attachment.Filename)
		if err != nil {
			return nil, err
		}

		encryptedContentType, err := ageencryption.EncryptString(publicKey, attachment.ContentType)
		if err != nil {
			return nil, err
		}

		encryptedMail.Attachments = append(encryptedMail.Attachments, RawAttachment{
			Filename:    encryptedFilename,
			ContentType: encryptedContentType,
			Data:        encryptedAttachment,
		})
	}

	return encryptedMail, nil
}

func (m *RawMail) ToMailEntity() *Mail {
	return &Mail{
		Headers:        m.Headers,
		TextContent:    m.TextContent,
		HTMLContent:    m.HTMLContent,
		Attachments:    make([]MailAttachment, len(m.Attachments)),
		Rejected:       &m.Rejected,
		RewriteSubject: &m.RewriteSubject,
		Greylisted:     &m.Greylisted,
	}
}
