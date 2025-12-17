package mail

import (
	"context"
	"io"
	"os"
	"strings"

	"connectrpc.com/connect"
	"github.com/appleboy/go-fcm"
	authv1 "github.com/atomic-blend/backend/grpc/gen/auth/v1"
	calendarv1 "github.com/atomic-blend/backend/grpc/gen/calendar/v1"
	userv1 "github.com/atomic-blend/backend/grpc/gen/user/v1"
	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/mail/notifications/payloads"
	"github.com/atomic-blend/backend/mail/repositories"
	icalparser "github.com/atomic-blend/backend/mail/utils/ical_parser"
	calendarclient "github.com/atomic-blend/backend/shared/grpc/calendar"
	userclient "github.com/atomic-blend/backend/shared/grpc/user"
	ageencryptionservice "github.com/atomic-blend/backend/shared/services/age_encryption"
	rspamdservice "github.com/atomic-blend/backend/shared/services/rspamd"
	rspamdclient "github.com/atomic-blend/backend/shared/services/rspamd/client"
	s3service "github.com/atomic-blend/backend/shared/services/s3"
	"github.com/atomic-blend/backend/shared/utils/db"
	fcmutils "github.com/atomic-blend/backend/shared/utils/fcm_utils"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/emersion/go-message"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	htmlcharset "golang.org/x/net/html/charset"
)

func receiveMail(m *amqp.Delivery, payload ReceivedMailPayload) {
	mailRepository := repositories.NewMailRepository(db.Database)
	s3Service, err := s3service.NewS3Service()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create S3 service")
		return
	}

	// Create a reader from the MIME content string
	reader := strings.NewReader(payload.Content)

	// Send the email to rspamd via HTTP for spam detection
	rspamdService := rspamdservice.NewRspamdService()

	// Create the FCM client to send notifications to the user
	firebaseProjectID := os.Getenv("FIREBASE_PROJECT_ID")
	googleApplicationCredentials := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	fcmClient, err := fcm.NewClient(
		context.TODO(),
		fcm.WithProjectID(firebaseProjectID),
		fcm.WithCredentialsFile(googleApplicationCredentials),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create FCM client")
	}

	mailContent := &models.RawMail{
		Attachments:    make([]models.RawAttachment, 0),
		Rejected:       false,
		RewriteSubject: false,
		Greylisted:     false,
	}

	checkRequest := &rspamdclient.CheckRequest{
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

	checkResponse, err := rspamdService.CheckMessage(checkRequest)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check message with Rspamd")
		// Continue processing even if Rspamd check fails
	} else {
		log.Debug().
			Str("action", checkResponse.Action).
			Float64("score", checkResponse.Score).
			Float64("required_score", checkResponse.RequiredScore).
			Bool("is_spam", checkResponse.IsSpam()).
			Msg("Rspamd check completed")

		// Log triggered symbols if any
		if len(checkResponse.Symbols) > 0 {
			log.Debug().Interface("symbols", checkResponse.Symbols).Msg("Rspamd triggered symbols")
		}

		switch checkResponse.Action {
		case "reject":
			log.Debug().Msg("Rejecting email")
			mailContent.Rejected = true
		case "soft reject":
			log.Debug().Msg("Soft rejecting email")
			mailContent.Rejected = true
		case "no action":
			log.Debug().Msg("No action taken")
		case "add header":
			log.Debug().Msg("Adding spam header")
			mailContent.RewriteSubject = true
		case "rewrite subject":
			log.Debug().Msg("Rewrite subject")
			// mark the email subject as needing a rewrite (only when sending, ignored on receiving)
			mailContent.RewriteSubject = true
		case "greylist":
			log.Debug().Msg("Greylisting email")
			mailContent.Greylisted = true
		default:
			log.Debug().Msg("No action taken")
		}
	}

	// Register a charset reader so charsets like iso-8859-1 are handled
	message.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return htmlcharset.NewReaderLabel(charset, input)
	}

	// Parse the MIME message
	entity, err := message.Read(reader)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse MIME message")
		return
	}

	log.Debug().
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

	// Create a new calendar client
	calendarService, err := calendarclient.NewCalendarClient()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create calendar client")
		return
	}

	encryptedMails := map[string]models.Mail{}
	encryptedNotifications := make(map[string]payloads.MailReceivedPayload, 0)
	encryptedAttachments := make([]*awss3.PutObjectInput, 0)
	calendarPayloads := map[string]*calendarv1.Calendar{}
	haveErrors := false

	for _, rcpt := range payload.Rcpt {
		log.Debug().Str("rcpt", rcpt).Msg("Handling recepient")

		mailEntity := &models.Mail{}

		// get the user public key from the auth service via grpc
		log.Debug().Str("rcpt", rcpt).Msg("Instantiating user client")
		userClient, err := userclient.NewUserClient()
		if err != nil {
			log.Error().Err(err).Msg("Failed to create user client")
			haveErrors = true
			continue
		}

		log.Debug().Str("rcpt", rcpt).Msg("Getting user public key")
		rcptPublicKey, err := userClient.GetUserPublicKey(context.Background(), &connect.Request[userv1.GetUserPublicKeyRequest]{
			Msg: &userv1.GetUserPublicKeyRequest{
				Email: rcpt,
			},
		})
		if err != nil {
			log.Debug().Str("rcpt", rcpt).Msg("User not found, skipping")
			continue
		}

		userID, err := primitive.ObjectIDFromHex(rcptPublicKey.Msg.UserId)
		if err != nil {
			log.Error().Err(err).Msg("Failed to convert user ID to ObjectID")
			haveErrors = true
			continue
		}
		mailEntity.UserID = userID

		log.Debug().Str("rcpt", rcpt).Str("publicKey", rcptPublicKey.Msg.PublicKey).Msg("User public key")
		log.Debug().Interface("encryptedMails", encryptedMails).Msg("Encrypted mails")

		userPublicKey := rcptPublicKey.Msg.PublicKey

		// encrypt the mail content for mongodb with user's public key
		log.Debug().Str("rcpt", rcpt).Msg("Encrypting mail content")
		encryptedMailContent, err := mailContent.Encrypt(userPublicKey)
		if err != nil {
			log.Error().Err(err).Msg("Failed to encrypt mail content")
			haveErrors = true
			continue
		}

		log.Debug().Str("rcpt", rcpt).Interface("encryptedContent", encryptedMailContent).Msg("Encrypted mail content")

		var calendarAttachment *models.RawAttachment

		// capture the calendar attachment only if there's a single attachment and it's a calendar file
		// for now, we consider that a calendar event is a single email with an ics attachment
		log.Debug().Str("rcpt", rcpt).Int("attachmentCount", len(mailContent.Attachments)).Msg("Checking for calendar attachment")
		log.Debug().Str("rcpt", rcpt).Interface("attachments", mailContent.Attachments).Msg("Listing attachments")
		if len(mailContent.Attachments) == 1 && (mailContent.Attachments[0].ContentType == "text/calendar" || strings.HasSuffix(mailContent.Attachments[0].Filename, ".ics")) {
			calendarAttachment = &mailContent.Attachments[0]
			log.Debug().Str("rcpt", rcpt).Str("filename", calendarAttachment.Filename).Msg("Found calendar attachment")
		}

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

		// NOTE: If calendar parsing fails, the mail will still be saved without the calendar event.
		// This is not treated as a fatal error, but a warning is logged for visibility.  ad
		if calendarAttachment != nil {
			log.Debug().Str("rcpt", rcpt).Str("filename", calendarAttachment.Filename).Msg("Parsing calendar attachment")

			calendars, err := icalparser.ParseICal(calendarAttachment.Data)
			if err != nil {
				log.Error().Err(err).Msg("Failed to parse iCal data")
				haveErrors = true
				continue
			}

			log.Debug().Str("rcpt", rcpt).Int("calendarCount", len(calendars)).Msg("Parsed iCal calendars")

			// handle only the first calendar and first event for now
			if len(calendars) == 0 || len(calendars[0].Events()) == 0 {
				log.Debug().Str("rcpt", rcpt).Msg("No calendar events found")
				continue
			}

			// convert the calendar to calendar payload
			newCalendarPayload, err := icalparser.ToCalendarPayload(calendars[0])
			if err != nil {
				log.Error().Err(err).Msg("Failed to convert calendar to payload")
			}

			calendarPayloads[userID.Hex()] = newCalendarPayload
			log.Debug().Str("rcpt", rcpt).Msg("Converted calendar to payload")
		}

		// set the mail entity fields
		mailEntity.Headers = encryptedMailContent.Headers

		mailEntity.TextContent = encryptedMailContent.TextContent
		mailEntity.HTMLContent = encryptedMailContent.HTMLContent
		mailEntity.Rejected = boolPtr(encryptedMailContent.Rejected)
		mailEntity.RewriteSubject = boolPtr(encryptedMailContent.RewriteSubject)
		mailEntity.Greylisted = boolPtr(encryptedMailContent.Greylisted)

		encryptedMails[userID.Hex()] = *mailEntity

		// encrypt the notification content for mongodb with user's public key
		log.Debug().Str("rcpt", rcpt).Msg("Encrypting notification content")
		ageService := ageencryptionservice.NewAgeEncryptionService()
		contentPreview := truncateString(mailContent.TextContent, 100)
		encryptedContentPreview, err := ageService.EncryptString(userPublicKey, contentPreview)
		if err != nil {
			log.Error().Err(err).Msg("Failed to encrypt content preview")
			haveErrors = true
			continue
		}
		encryptedNotificationContent := payloads.NewMailReceivedPayload(encryptedMailContent.Headers["From"].(string), encryptedMailContent.Headers["Subject"].(string), encryptedContentPreview)
		encryptedNotifications[userID.Hex()] = *encryptedNotificationContent
	}

	if haveErrors {
		log.Error().Msg("Errors occurred while processing email")
		return
	}

	// upload the attachments to s3 in bulk
	log.Debug().Msg("Uploading attachments to S3 in bulk")
	uploadedKeys, err := s3Service.BulkUploadFiles(context.TODO(), encryptedAttachments)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upload attachments to S3")
		return
	}

	// create calendars for users with calendar payloads
	for userID, calendarPayload := range calendarPayloads {
		log.Debug().Str("userID", userID).Msg("Creating calendar for user")
		log.Debug().Str("userID", userID).Interface("calendarPayload", calendarPayload).Msg("Calendar payload")
		if calendarPayload == nil {
			continue
		}
		// Create a new calendar request
		req := calendarclient.CreateCreateCalendarRequest(&authv1.User{Id: userID}, calendarPayload)

		// Call the CreateCalendar method
		response, err := calendarService.CreateCalendar(context.TODO(), req)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create calendar")
			return
		}

		// convert the returned calendar event ID to ObjectID
		calendarEventID, err := primitive.ObjectIDFromHex(*response.Msg.Id)
		if err != nil {
			log.Error().Err(err).Msg("Failed to convert calendar event ID to ObjectID")
			return
		}

		// link the created calendar event to the mail entity
		mail := encryptedMails[userID]
		mail.CalendarEvent = &calendarEventID
		encryptedMails[userID] = mail
	}

	// save the mail documents with s3 references to mongodb
	log.Debug().Int("count", len(encryptedMails)).Msg("Saving mail documents to MongoDB")
	// gather all the mail entities into a slice
	mailsToCreate := make([]models.Mail, 0, len(encryptedMails))
	for _, mail := range encryptedMails {
		mailsToCreate = append(mailsToCreate, mail)
	}
	_, err = mailRepository.CreateMany(context.TODO(), mailsToCreate)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save mail documents to MongoDB")
		s3Service.BulkDeleteFiles(context.TODO(), uploadedKeys)
		return
	}

	// send notifications to the user
	log.Debug().Msg("Sending notifications to users")
	for userID, notification := range encryptedNotifications {

		// Get user devices using gRPC client
		log.Debug().Str("userID", userID).Msg("Getting user devices for notification")
		req := &connect.Request[userv1.GetUserDevicesRequest]{
			Msg: &userv1.GetUserDevicesRequest{
				User: &authv1.User{
					Id: userID,
				},
			},
		}

		userService, err := userclient.NewUserClient()
		if err != nil {
			log.Error().Err(err).Msg("Failed to create user client")
			return
		}

		log.Debug().Str("userID", userID).Msg("Fetching user devices via gRPC")
		resp, err := userService.GetUserDevices(context.TODO(), req)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to get user devices for user: %s", userID)
			continue
		}

		deviceTokens := []string{}
		for _, device := range resp.Msg.Devices {
			if device.FcmToken != "" {
				deviceTokens = append(deviceTokens, device.FcmToken)
			}
		}

		if len(deviceTokens) == 0 {
			log.Debug().Str("userID", userID).Msg("No device tokens found for user")
			continue
		}
		data := notification.GetData()

		log.Debug().Msgf("Payload: %v", payload)
		log.Debug().Msgf("Data: %v", data)

		// send the notification to the user
		log.Debug().Str("userID", userID).Msg("Sending FCM notification to user devices")
		fcmutils.SendMulticast(context.TODO(), fcmClient, data, deviceTokens)
	}

	m.Ack(false)
}

func processMessageBody(entity *message.Entity, mailContent *models.RawMail) {
	// Extract all headers first
	headers := make(map[string]interface{})
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

func processMultipartMessage(entity *message.Entity, mailContent *models.RawMail) {
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

func processMessagePart(part *message.Entity, mailContent *models.RawMail) {
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
		attachment := models.RawAttachment{
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

// Truncate returns the first n runes of s.
func truncateString(s string, n int) string {
	if len(s) <= n {
		return s
	}
	for i := range s {
		if n == 0 {
			return s[:i]
		}
		n--
	}
	return s
}
