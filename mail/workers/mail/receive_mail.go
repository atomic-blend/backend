package mail

import (
	"io"
	"strings"

	"github.com/emersion/go-message"
	"github.com/rs/zerolog/log"
)

func receiveMail(mimeContent string) {
	// Create a reader from the MIME content string
	reader := strings.NewReader(mimeContent)

	// Parse the MIME message
	entity, err := message.Read(reader)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse MIME message")
		return
	}

	// Extract basic email information
	header := entity.Header
	from := header.Get("From")
	to := header.Get("To")
	subject := header.Get("Subject")
	date := header.Get("Date")

	log.Info().
		Str("from", from).
		Str("to", to).
		Str("subject", subject).
		Str("date", date).
		Msg("Received email")

	// TODO: Create base mail document with headers and metadata
	// TODO: Initialize mail document with from, to, subject, date, messageId, etc.
	// TODO: Generate unique mail ID for this email

	// Process the message body
	processMessageBody(entity)
}

func processMessageBody(entity *message.Entity) {
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
		processMultipartMessage(entity)
	} else {
		// Handle single part message
		processSinglePartMessage(entity)
	}
}

func processMultipartMessage(entity *message.Entity) {
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
		processMessagePart(part)
	}
}

func processSinglePartMessage(entity *message.Entity) {
	// Read the body content
	body, err := io.ReadAll(entity.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read message body")
		return
	}

	// Get content type for this part
	mediaType, _, err := entity.Header.ContentType()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get part content type")
		return
	}

	log.Info().
		Str("contentType", mediaType).
		Int("bodySize", len(body)).
		Msg("Single part message")

	// Handle different content types
	switch {
	case strings.HasPrefix(mediaType, "text/plain"):
		log.Info().Str("textContent", string(body)).Msg("Plain text content")
		// TODO: Store plain text content in MongoDB
		// TODO: Create mail document with text content, headers, and metadata
		// TODO: Save to mail collection in MongoDB
	case strings.HasPrefix(mediaType, "text/html"):
		log.Info().Str("htmlContent", string(body)).Msg("HTML content")
		// TODO: Store HTML content in MongoDB
		// TODO: Create mail document with HTML content, headers, and metadata
		// TODO: Save to mail collection in MongoDB
	default:
		log.Info().Str("contentType", mediaType).Msg("Other content type")
	}
}

func processMessagePart(part *message.Entity) {
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
		log.Info().Str("textContent", string(body)).Msg("Plain text part")
		// TODO: Store plain text content in MongoDB
		// TODO: Create mail document with text content, headers, and metadata
		// TODO: Save to mail collection in MongoDB
	case strings.HasPrefix(mediaType, "text/html"):
		log.Info().Str("htmlContent", string(body)).Msg("HTML part")
		// TODO: Store HTML content in MongoDB
		// TODO: Create mail document with HTML content, headers, and metadata
		// TODO: Save to mail collection in MongoDB
	default:
		log.Info().Str("contentType", mediaType).Str("filename", filename).Msg("Other content type part")
		// TODO: Store other file types in S3 storage
		// TODO: Use original filename or generate unique filename for the file
		// TODO: Upload file bytes to S3 bucket with filename
		// TODO: Store S3 file reference and original filename in MongoDB mail document
	}
}
