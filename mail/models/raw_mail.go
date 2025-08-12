package models

import (
	"bytes"
	"fmt"
	"io"

	ageencryption "github.com/atomic-blend/backend/mail/utils/age_encryption"
	"github.com/emersion/go-message"
)

// RawMail represents the collected content from an email
type RawMail struct {
	Headers        map[string]interface{} `json:"headers" binding:"required"`
	TextContent    string                 `json:"text_content"`
	HTMLContent    string                 `json:"html_content"`
	Attachments    []RawAttachment        `json:"attachments"`
	Rejected       bool                   `json:"rejected"`
	RewriteSubject bool                   `json:"rewrite_subject"`
	Greylisted     bool                   `json:"graylisted"`
}

// RawAttachment represents a file attachment
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

// ToMailEntity converts a RawMail to a Mail entity
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

// ToMessageEntity creates a proper email message using go-message library
func (m *RawMail) ToMessageEntity() (*message.Entity, error) {
	var buf bytes.Buffer
	var header message.Header

	// Set headers
	for key, value := range m.Headers {
		switch v := value.(type) {
		case string:
			header.Set(key, v)
		case []string:
			for _, item := range v {
				header.Add(key, item)
			}
		default:
			header.Set(key, fmt.Sprintf("%v", v))
		}
	}

	// Create multipart message if we have both HTML and text content
	if m.HTMLContent != "" && m.TextContent != "" {
		header.SetContentType("multipart/alternative", nil)
		w, err := message.CreateWriter(&buf, header)
		if err != nil {
			return nil, fmt.Errorf("failed to create multipart writer: %w", err)
		}

		// Add text part
		textHeader := message.Header{}
		textHeader.SetContentType("text/plain", map[string]string{"charset": "utf-8"})
		textHeader.Set("Content-Transfer-Encoding", "quoted-printable")
		textPart, err := w.CreatePart(textHeader)
		if err != nil {
			w.Close()
			return nil, fmt.Errorf("failed to create text part: %w", err)
		}
		io.WriteString(textPart, m.TextContent)
		textPart.Close()

		// Add HTML part
		htmlHeader := message.Header{}
		htmlHeader.SetContentType("text/html", map[string]string{"charset": "utf-8"})
		htmlHeader.Set("Content-Transfer-Encoding", "quoted-printable")
		htmlPart, err := w.CreatePart(htmlHeader)
		if err != nil {
			w.Close()
			return nil, fmt.Errorf("failed to create HTML part: %w", err)
		}
		io.WriteString(htmlPart, m.HTMLContent)
		htmlPart.Close()

		w.Close()
	} else if m.HTMLContent != "" {
		// HTML only
		header.SetContentType("text/html", map[string]string{"charset": "utf-8"})
		header.Set("Content-Transfer-Encoding", "quoted-printable")
		w, err := message.CreateWriter(&buf, header)
		if err != nil {
			return nil, fmt.Errorf("failed to create HTML writer: %w", err)
		}
		io.WriteString(w, m.HTMLContent)
		w.Close()
	} else if m.TextContent != "" {
		// Text only
		header.SetContentType("text/plain", map[string]string{"charset": "utf-8"})
		header.Set("Content-Transfer-Encoding", "quoted-printable")
		w, err := message.CreateWriter(&buf, header)
		if err != nil {
			return nil, fmt.Errorf("failed to create text writer: %w", err)
		}
		io.WriteString(w, m.TextContent)
		w.Close()
	}

	// Add attachments if any
	if len(m.Attachments) > 0 {
		// Create a new buffer for the multipart message with attachments
		var attachmentBuf bytes.Buffer
		attachmentHeader := message.Header{}
		attachmentHeader.SetContentType("multipart/mixed", nil)

		w, err := message.CreateWriter(&attachmentBuf, attachmentHeader)
		if err != nil {
			return nil, fmt.Errorf("failed to create attachment multipart writer: %w", err)
		}

		// Add the original content as first part
		contentHeader := message.Header{}
		if m.HTMLContent != "" {
			contentHeader.SetContentType("text/html", map[string]string{"charset": "utf-8"})
		} else {
			contentHeader.SetContentType("text/plain", map[string]string{"charset": "utf-8"})
		}
		contentHeader.Set("Content-Transfer-Encoding", "quoted-printable")

		contentPart, err := w.CreatePart(contentHeader)
		if err != nil {
			w.Close()
			return nil, fmt.Errorf("failed to create content part: %w", err)
		}

		if m.HTMLContent != "" {
			io.WriteString(contentPart, m.HTMLContent)
		} else {
			io.WriteString(contentPart, m.TextContent)
		}
		contentPart.Close()

		// Add attachments
		for _, attachment := range m.Attachments {
			attachmentPartHeader := message.Header{}
			attachmentPartHeader.SetContentType(attachment.ContentType, nil)
			attachmentPartHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", attachment.Filename))
			attachmentPartHeader.Set("Content-Transfer-Encoding", "base64")

			attachmentPart, err := w.CreatePart(attachmentPartHeader)
			if err != nil {
				w.Close()
				return nil, fmt.Errorf("failed to create attachment part: %w", err)
			}

			attachmentPart.Write(attachment.Data)
			attachmentPart.Close()
		}

		w.Close()

		// Replace the original buffer with the attachment buffer
		buf = attachmentBuf
	}

	// Create the final message entity
	msg, err := message.Read(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read message from buffer: %w", err)
	}

	return msg, nil
}
