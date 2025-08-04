package mail

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"strings"

	"connectrpc.com/connect"
	"filippo.io/age"
	userv1 "github.com/atomic-blend/backend/grpc/gen/user/v1"
	"github.com/atomic-blend/backend/mail/grpc/clients"
	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/mail/utils/rspamd"
	"github.com/emersion/go-message"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
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

func receiveMail(m *amqp.Delivery, payload ReceivedMailPayload) {
	// Create a reader from the MIME content string
	reader := strings.NewReader(payload.Content)

	// Send the email to rspamd via HTTP for spam detection
	client := rspamd.NewClient(nil) // Use default config

	mailContent := &MailContent{
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

	// Store headers
	mailContent.Headers.From = entity.Header.Get("From")
	mailContent.Headers.To = entity.Header.Get("To")
	mailContent.Headers.Subject = entity.Header.Get("Subject")
	mailContent.Headers.Date = entity.Header.Get("Date")

	log.Info().
		Str("from", mailContent.Headers.From).
		Str("to", mailContent.Headers.To).
		Str("subject", mailContent.Headers.Subject).
		Str("date", mailContent.Headers.Date).
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
	haveErrors := false

	for _, rcpt := range strings.Split(mailContent.Headers.To, ",") {
		log.Info().Str("rcpt", rcpt).Msg("Handling recepient")

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
			log.Error().Err(err).Msg("Failed to get user public key")
			haveErrors = true
			continue
		}

		log.Info().Str("rcpt", rcpt).Str("publicKey", rcptPublicKey.Msg.PublicKey).Msg("User public key")
		log.Info().Interface("encryptedMails", encryptedMails).Msg("Encrypted mails")

		userPublicKey := rcptPublicKey.Msg.PublicKey

		// encrypt the mail content for mongodb with user's public key
		log.Info().Str("rcpt", rcpt).Msg("Encrypting mail content")
		recipient, err := age.ParseX25519Recipient(userPublicKey)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse public key")
			haveErrors = true
			continue
		}

		out := &bytes.Buffer{}

		w, err := age.Encrypt(out, recipient)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create encrypted file")
			haveErrors = true
			continue
		}
		if _, err := io.WriteString(w, "Black lives matter."); err != nil {
			log.Error().Err(err).Msg("Failed to write to encrypted file")
			haveErrors = true
			continue
		}
		if err := w.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close encrypted file")
			haveErrors = true
			continue
		}

		// use base64 decode + strings.NewReader to get the encrypted content when decrypting
		// TODO: refactor the encryption / decryption logic into a dedicated util
		encryptedContent := base64.StdEncoding.EncodeToString(out.Bytes())

		log.Info().Str("rcpt", rcpt).Str("encryptedContent", encryptedContent).Msg("Encrypted mail content")

		//TODO: encrypt other fields of the email (headers, etc...)

		//TODO: encrypt the files in the attachments

		//TODO: upload the attachments to s3

		//TODO: save the mail document with s3 references to mongodb
	}

	if haveErrors {
		log.Error().Msg("Errors occurred while processing email")
		return
	}

	m.Ack(false)
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
