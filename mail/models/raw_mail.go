package models

import (
	"fmt"

	ageencryption "github.com/atomic-blend/backend/mail/utils/age_encryption"
)

// RawMail represents the collected content from an email
type RawMail struct {
	Headers        map[string]interface{}
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
			encryptedHeaders := make(map[string]interface{})
		for key, value := range m.Headers {
			switch v := value.(type) {
			case string:
				encryptedValue, err := ageencryption.EncryptString(publicKey, v)
				if err != nil {
						return nil, err
					}
					encryptedHeaders[key] = encryptedValue
				case []string:
					// For slice of strings, encrypt each item and store as slice
					encryptedSlice := make([]string, len(v))
					for i, item := range v {
						encryptedItem, err := ageencryption.EncryptString(publicKey, item)
						if err != nil {
							return nil, err
						}
						encryptedSlice[i] = encryptedItem
					}
					encryptedHeaders[key] = encryptedSlice
				default:
					// For other types, convert to string representation and encrypt
					valueStr := fmt.Sprintf("%v", v)
					encryptedValue, err := ageencryption.EncryptString(publicKey, valueStr)
					if err != nil {
						return nil, err
					}
					encryptedHeaders[key] = encryptedValue
				}
			}
			encryptedMail.Headers = encryptedHeaders
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
