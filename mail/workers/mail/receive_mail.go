package mail

import (
	"context"
	"io"
	"os"
	"strings"

	"connectrpc.com/connect"
	userv1 "github.com/atomic-blend/backend/grpc/gen/user/v1"
	"github.com/atomic-blend/backend/mail/grpc/clients"
	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/mail/repositories"
	ageencryption "github.com/atomic-blend/backend/mail/utils/age_encryption"
	"github.com/atomic-blend/backend/mail/utils/db"
	"github.com/atomic-blend/backend/mail/utils/rspamd"
	"github.com/atomic-blend/backend/mail/utils/s3"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/emersion/go-message"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Content represents the collected content from an email
type Content struct {
	Headers        interface{}
	TextContent    string
	HTMLContent    string
	Attachments    []Attachment
	Rejected       bool
	RewriteSubject bool
	Greylisted     bool
}

// Attachment represents a file attachment
type Attachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

// Encrypt encrypts the content using the age encryption library
func (m *Content) Encrypt(publicKey string) (*Content, error) {
	encryptedMail := &Content{
		Attachments:    make([]Attachment, 0),
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

		encryptedMail.Attachments = append(encryptedMail.Attachments, Attachment{
			Filename:    encryptedFilename,
			ContentType: encryptedContentType,
			Data:        encryptedAttachment,
		})
	}

	return encryptedMail, nil
}

func receiveMail(m *amqp.Delivery, payload ReceivedMailPayload) {
	mailRepository := repositories.NewMailRepository(db.Database)
	s3Service, err := s3.NewS3Service(os.Getenv("AWS_BUCKET"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create S3 service")
		return
	}

	// Create a reader from the MIME content string
	reader := strings.NewReader(payload.Content)

	// Send the email to rspamd via HTTP for spam detection
	client := rspamd.NewClient(nil) // Use default config

	mailContent := &Content{
		Attachments:    make([]Attachment, 0),
		Rejected:       false,
		RewriteSubject: false,
		Greylisted:     false,
	}

	checkRequest := &rspamd.CheckRequest{
		Message:   []byte(payload.Content),
		IP:        payload.IP,
		Helo:      payload.Hostname,
		Hostname:  payload.Hostname,
		From:      payload.From,
		Rcpt:      payload.Rcpt,
		QueueID:   payload.QueueID,
		User:      payload.User,
		DeliverTo: payload.DeliverTo,
	}

	checkResponse, err := client.CheckMessage(checkRequest)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check message with Rspamd")
		// Continue processing even if Rspamd check fails
	} else {
		log.Info().
			Str("action", checkResponse.Action).
			Float64("score", checkResponse.Score).
			Float64("required_score", checkResponse.RequiredScore).
			Bool("is_spam", checkResponse.IsSpam()).
			Msg("Rspamd check completed")

		// Log triggered symbols if any
		if len(checkResponse.Symbols) > 0 {
			log.Info().Interface("symbols", checkResponse.Symbols).Msg("Rspamd triggered symbols")
		}

		switch checkResponse.Action {
		case "reject":
			log.Info().Msg("Rejecting email")
			mailContent.Rejected = true
		case "soft reject":
			log.Info().Msg("Soft rejecting email")
			mailContent.Rejected = true
		case "no action":
			log.Info().Msg("No action taken")
		case "add header":
			log.Info().Msg("Adding spam header")
			mailContent.RewriteSubject = true
		case "rewrite subject":
			log.Info().Msg("Rewrite subject")
			// mark the email subject as needing a rewrite (only when sending, ignored on receiving)
			mailContent.RewriteSubject = true
		case "greylist":
			log.Info().Msg("Greylisting email")
			mailContent.Greylisted = true
		default:
			log.Info().Msg("No action taken")
		}
	}

	// Parse the MIME message
	entity, err := message.Read(reader)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse MIME message")
		return
	}

	log.Info().
		Str("from", payload.From).
		Interface("to", payload.Rcpt).
		Str("date", payload.ReceivedAt).
		Str("client_ip", payload.IP).
		Str("hostname", payload.Hostname).
		Str("queue_id", payload.QueueID).
		Str("deliver_to", payload.DeliverTo).
		Str("received_at", payload.ReceivedAt).
		Interface("recipients", payload.Rcpt).
		Bool("rejected", mailContent.Rejected).
		Bool("rewrite_subject", mailContent.RewriteSubject).
		Bool("greylisted", mailContent.Greylisted).
		Msg("Received email")

	// Process the message body and collect all content
	processMessageBody(entity, mailContent)

	encryptedMails := make([]models.Mail, 0)
	encryptedAttachments := make([]*awss3.PutObjectInput, 0)
	haveErrors := false

	for _, rcpt := range payload.Rcpt {
		log.Info().Str("rcpt", rcpt).Msg("Handling recepient")

		mailEntity := &models.Mail{}

		// get the user public key from the auth service via grpc
		log.Info().Str("rcpt", rcpt).Msg("Instantiating user client")
		userClient, err := clients.NewUserClient()
		if err != nil {
			log.Error().Err(err).Msg("Failed to create user client")
			haveErrors = true
			continue
		}

		log.Info().Str("rcpt", rcpt).Msg("Getting user public key")
		rcptPublicKey, err := userClient.GetUserPublicKey(context.Background(), &connect.Request[userv1.GetUserPublicKeyRequest]{
			Msg: &userv1.GetUserPublicKeyRequest{
				Email: rcpt,
			},
		})
		if err != nil {
			log.Info().Str("rcpt", rcpt).Msg("User not found, skipping")
			continue
		}

		userID, err := primitive.ObjectIDFromHex(rcptPublicKey.Msg.UserId)
		if err != nil {
			log.Error().Err(err).Msg("Failed to convert user ID to ObjectID")
			haveErrors = true
			continue
		}
		mailEntity.UserID = userID

		log.Info().Str("rcpt", rcpt).Str("publicKey", rcptPublicKey.Msg.PublicKey).Msg("User public key")
		log.Info().Interface("encryptedMails", encryptedMails).Msg("Encrypted mails")

		userPublicKey := rcptPublicKey.Msg.PublicKey

		// encrypt the mail content for mongodb with user's public key
		log.Info().Str("rcpt", rcpt).Msg("Encrypting mail content")
		encryptedMailContent, err := mailContent.Encrypt(userPublicKey)
		if err != nil {
			log.Error().Err(err).Msg("Failed to encrypt mail content")
			haveErrors = true
			continue
		}

		log.Info().Str("rcpt", rcpt).Interface("encryptedContent", encryptedMailContent).Msg("Encrypted mail content")

		// upload the attachments to s3 and store the references in the mail entity
		for _, attachment := range encryptedMailContent.Attachments {
			uniqueFileID := uuid.New().String()
			payload, err := s3Service.GenerateUploadPayload(context.Background(), attachment.Data, "mail/attachments/"+rcptPublicKey.Msg.UserId, uniqueFileID, map[string]string{})
			if err != nil {
				log.Error().Err(err).Msg("Failed to upload attachment to S3")
				haveErrors = true
				continue
			}
			mailEntity.Attachments = append(mailEntity.Attachments, models.MailAttachment{
				StoragePath: *payload.Key,
				Filename:    attachment.Filename,
				ContentType: attachment.ContentType,
				StorageType: "s3",
				Size:        int64(len(attachment.Data)),
			})
			encryptedAttachments = append(encryptedAttachments, payload)
		}

		// set the mail entity fields
		mailEntity.Headers = encryptedMailContent.Headers

		mailEntity.TextContent = encryptedMailContent.TextContent
		mailEntity.HTMLContent = encryptedMailContent.HTMLContent
		mailEntity.Rejected = boolPtr(encryptedMailContent.Rejected)
		mailEntity.RewriteSubject = boolPtr(encryptedMailContent.RewriteSubject)
		mailEntity.Greylisted = boolPtr(encryptedMailContent.Greylisted)

		encryptedMails = append(encryptedMails, *mailEntity)
	}

	if haveErrors {
		log.Error().Msg("Errors occurred while processing email")
		return
	}

	// upload the attachments to s3 in bulk
	uploadedKeys, err := s3Service.BulkUploadFiles(context.Background(), encryptedAttachments)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upload attachments to S3")
		return
	}

	// save the mail documents with s3 references to mongodb
	_, err = mailRepository.CreateMany(context.Background(), encryptedMails)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save mail documents to MongoDB")
		s3Service.BulkDeleteFiles(context.Background(), uploadedKeys)
		return
	}

	m.Ack(false)
}

func processMessageBody(entity *message.Entity, mailContent *Content) {
	// Extract all headers first
	headers := make(map[string]string)
	for field := entity.Header.Fields(); field.Next(); {
		key := field.Key()
		value, _ := field.Text()
		headers[key] = value
	}
	mailContent.Headers = headers

	// Get the media type of the message
	mediaType, params, err := entity.Header.ContentType()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get content type")
		return
	}

	log.Info().
		Str("mediaType", mediaType).
		Interface("params", params).
		Interface("headers", headers).
		Msg("Message content type and headers")

	// Handle multipart messages
	if strings.HasPrefix(mediaType, "multipart/") {
		processMultipartMessage(entity, mailContent)
	} else {
		// Handle single part message
		processMessagePart(entity, mailContent)
	}
}

func processMultipartMessage(entity *message.Entity, mailContent *Content) {
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

func processMessagePart(part *message.Entity, mailContent *Content) {
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

	// Get content disposition (optional - some parts may not have it)
	disposition := ""
	var dispositionParams map[string]string
	if contentDisposition := part.Header.Get("Content-Disposition"); contentDisposition != "" {
		disposition, dispositionParams, err = part.Header.ContentDisposition()
		if err != nil {
			log.Warn().Err(err).Msg("Failed to parse content disposition, continuing without it")
			disposition = ""
			dispositionParams = nil
		}
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

func boolPtr(b bool) *bool {
	if b {
		return &b
	}
	return nil
}