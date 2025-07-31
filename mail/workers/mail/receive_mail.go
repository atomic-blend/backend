package mail

import (
	"io"
	"strings"

	"github.com/emersion/go-message"
	"github.com/rs/zerolog/log"
)

// MailContent represents the collected content from an email
type MailContent struct {
	Headers struct {
		From      string
		To        string
		Subject   string
		Date      string
		MessageID string
		Cc        string
	}
	TextContent string
	HTMLContent string
	Attachments []Attachment
}

// Attachment represents a file attachment
type Attachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

func receiveMail(mimeContent string) {
	// Create a reader from the MIME content string
	reader := strings.NewReader(mimeContent)

	//TODO: send the email to rspamd via HTTP for spam detection

	// Parse the MIME message
	entity, err := message.Read(reader)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse MIME message")
		return
	}

	// Initialize mail content structure
	mailContent := &MailContent{
		Attachments: make([]Attachment, 0),
	}

	// Extract basic email information
	header := entity.Header
	from := header.Get("From")
	to := header.Get("To")
	subject := header.Get("Subject")
	date := header.Get("Date")

	// Store headers
	mailContent.Headers.From = from
	mailContent.Headers.To = to
	mailContent.Headers.Subject = subject
	mailContent.Headers.Date = date

	log.Info().
		Str("from", from).
		Str("to", to).
		Str("subject", subject).
		Str("date", date).
		Msg("Received email")

	// Process the message body and collect all content
	processMessageBody(entity, mailContent)

	// TODO: Upload all collected content to MongoDB and S3
	// TODO: Create complete mail document with all content
	// TODO: Upload attachments to S3
	// TODO: Save mail document with S3 references to MongoDB
}

func processMessageBody(entity *message.Entity, mailContent *MailContent) {
	// Get the media type of the message
	mediaType, params, err := entity.Header.ContentType()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get content type")
		return
	}

	log.Info().
		Str("mediaType", mediaType).
		Interface("params", params).
		Msg("Message content type")

	// Handle multipart messages
	if strings.HasPrefix(mediaType, "multipart/") {
		processMultipartMessage(entity, mailContent)
	} else {
		// Handle single part message
		processMessagePart(entity, mailContent)
	}

	// TODO: Encrypt the mail content for mongodb
	// TODO: Encrypt the attachments for s3

	//TODO: upload the attachments to s3
	//TODO: save the mail document with s3 references to mongodb
	//TODO: upload the mail content to mongodb
}

func processMultipartMessage(entity *message.Entity, mailContent *MailContent) {
	// Create a multipart reader
	mr := entity.MultipartReader()
	if mr == nil {
		log.Error().Msg("Failed to create multipart reader")
		return
	}

	// Iterate through all parts
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error().Err(err).Msg("Failed to read multipart")
			break
		}

		// Process each part
		processMessagePart(part, mailContent)
	}
}

func processMessagePart(part *message.Entity, mailContent *MailContent) {
	// Read the part content
	body, err := io.ReadAll(part.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read part body")
		return
	}

	// Get content type for this part
	mediaType, params, err := part.Header.ContentType()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get part content type")
		return
	}

	// Get content disposition
	disposition, dispositionParams, err := part.Header.ContentDisposition()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get content disposition")
		return
	}

	// Extract filename from content disposition
	filename := ""
	if dispositionParams != nil {
		if name, exists := dispositionParams["filename"]; exists {
			filename = name
		}
	}

	// If filename is not in disposition params, try to get it from content type params
	if filename == "" && params != nil {
		if name, exists := params["name"]; exists {
			filename = name
		}
	}

	log.Info().
		Str("contentType", mediaType).
		Str("disposition", disposition).
		Str("filename", filename).
		Interface("dispositionParams", dispositionParams).
		Interface("contentParams", params).
		Int("bodySize", len(body)).
		Msg("Message part")

	// Handle different content types
	switch {
	case strings.HasPrefix(mediaType, "text/plain"):
		log.Info().Str("textContent", string(body)).Msg("Plain text content")
		mailContent.TextContent = string(body)
	case strings.HasPrefix(mediaType, "text/html"):
		log.Info().Str("htmlContent", string(body)).Msg("HTML content")
		mailContent.HTMLContent = string(body)
	default:
		log.Info().Str("contentType", mediaType).Str("filename", filename).Msg("Other content type")
		// TODO: Collect attachment data for later upload to S3
		attachment := Attachment{
			Filename:    filename,
			ContentType: mediaType,
			Data:        body,
		}
		mailContent.Attachments = append(mailContent.Attachments, attachment)
	}
}
