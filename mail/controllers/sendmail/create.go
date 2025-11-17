// Package sendmail contains the controller for the send mail endpoint
package sendmail

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"connectrpc.com/connect"
	userv1 "github.com/atomic-blend/backend/grpc/gen/user/v1"
	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

// CreateSendMailRequest represents the request payload for creating a send mail
type CreateSendMailRequest struct {
	Mail *models.Mail `json:"mail" binding:"required"`
}

// CreateSendMail creates a new send mail entry
// @Summary Create a new send mail
// @Description Create a new send mail entry with the provided mail data
// @Tags SendMail
// @Accept json
// @Produce json
// @Param body body CreateSendMailRequest true "Send mail data"
// @Success 201 {object} models.SendMail
// @Failure 401 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /mail/send [post]
func (c *Controller) CreateSendMail(ctx *gin.Context) {
	// get public address from environment for generating message IDs
	publicAddress := os.Getenv("PUBLIC_ADDRESS")
	if publicAddress == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Public address not configured"})
		return
	}

	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Bind JSON payload
	var rawMail models.RawMail
	if err := ctx.ShouldBindJSON(&rawMail); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Basic validation - ensure at least some content is provided
	if rawMail.TextContent == "" && rawMail.HTMLContent == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Mail content is required (either text or HTML)"})
		return
	}

	// Normalize headers to preserve list structure
	if rawMail.Headers != nil {
		normalizedHeaders := make(map[string]interface{})
		for key, value := range rawMail.Headers {
			switch v := value.(type) {
			case []interface{}:
				// Convert []interface{} to []string
				stringSlice := make([]string, len(v))
				for i, item := range v {
					if str, ok := item.(string); ok {
						stringSlice[i] = str
					} else {
						stringSlice[i] = fmt.Sprintf("%v", item)
					}
				}
				normalizedHeaders[key] = stringSlice
			default:
				// Keep other types as they are
				normalizedHeaders[key] = value
			}
		}
		rawMail.Headers = normalizedHeaders
	}

	// // generate a message ID header if not provided
	if rawMail.Headers == nil {
		rawMail.Headers = make(map[string]interface{})
	}

	messageID := fmt.Sprintf("<%s@%s>", uuid.New().String(), publicAddress)
	rawMail.Headers["Message-ID"] = messageID

	// Set In-Reply-To header if replying to another mail
	if rawMail.InReplyTo != nil {
		func() {
			// fetch the original mail to get its message ID
			originalMailID := *rawMail.InReplyTo
			objectID, err := primitive.ObjectIDFromHex(originalMailID)
			if err != nil {
				log.Error().Err(err).Str("original_mail_id", originalMailID).Msg("Invalid original mail ID for In-Reply-To header")
				return
			}
			originalMail, err := c.mailRepo.GetByID(ctx, objectID)
			if err != nil {
				log.Error().Err(err).Str("original_mail_id", originalMailID).Msg("Failed to fetch original mail for In-Reply-To header")
				return
			}

			originalMailMessageID := ""
			if originalMail != nil && originalMail.Headers != nil {
				if msgID, ok := originalMail.Headers["Message-ID"].(string); ok {
					originalMailMessageID = msgID
				}
			}

			if originalMailMessageID == "" {
				log.Error().Str("original_mail_id", originalMailID).Msg("Original mail has no Message-ID header for In-Reply-To")
				return
			}

			log.Debug().Str("in_reply_to", originalMailMessageID).Msg("Setting In-Reply-To header")
			rawMail.Headers["In-Reply-To"] = originalMailMessageID

			// Also set References header to include new message ID
			if refs, ok := originalMail.Headers["References"].(string); ok && refs != "" {
				rawMail.Headers["References"] = fmt.Sprintf("%s %s", refs, originalMailMessageID)
			} else {
				rawMail.Headers["References"] = originalMailMessageID
			}
		}()
	}

	//TODO: check email validity here
	log.Debug().Interface("raw_mail", rawMail).Msg("Received raw mail for sending")

	// get the user public key from the auth service via grpc
	log.Info().Msg("Getting user public key")
	userPublicKey, err := c.userClient.GetUserPublicKey(context.Background(), &connect.Request[userv1.GetUserPublicKeyRequest]{
		Msg: &userv1.GetUserPublicKeyRequest{
			Id: string(authUser.UserID.Hex()),
		},
	})
	if err != nil {
		log.Info().Msg("User not found, skipping")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	log.Info().Str("public_key", userPublicKey.Msg.PublicKey).Msg("User public key retrieved successfully")

	// TODO: [DONE] transform the raw mail to an encrypted mail
	encryptedMail, err := rawMail.Encrypt(userPublicKey.Msg.PublicKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to encrypt mail")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt mail"})
		return
	}

	// set the user ID on the encrypted mail
	encryptedMailEntity := encryptedMail.ToMailEntity()
	encryptedMailEntity.UserID = authUser.UserID

	log.Debug().Interface("encrypted_mail", encryptedMail).Msg("Encrypted mail ready for sending")

	// Create send mail entity
	sendMail := &models.SendMail{
		Mail:       encryptedMailEntity,
		SendStatus: models.SendStatusPending,
		Trashed:    false,
	}

	// Generate encrypted attachments upload requests to send them to S3
	encryptedAttachments := make([]*awss3.PutObjectInput, 0)
	for _, attachment := range encryptedMail.Attachments {
		uniqueFileID := uuid.New().String()
		payload, err := c.s3Service.GenerateUploadPayload(context.Background(), attachment.Data, "send_mail/attachments/"+userPublicKey.Msg.UserId, uniqueFileID, map[string]string{})
		if err != nil {
			log.Error().Err(err).Msg("Failed to upload attachment to S3")
		}
		sendMail.Mail.Attachments = append(sendMail.Mail.Attachments, models.MailAttachment{
			StoragePath: *payload.Key,
			Filename:    attachment.Filename,
			ContentType: attachment.ContentType,
			StorageType: "s3",
			Size:        int64(len(attachment.Data)),
		})
		encryptedAttachments = append(encryptedAttachments, payload)
	}

	// upload the attachments to s3 in bulk
	uploadedKeys, err := c.s3Service.BulkUploadFiles(context.Background(), encryptedAttachments)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upload attachments to S3")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload attachments"})
		return
	}

	// TODO: [DONE] setup transactional upload for mail with rollback on failure (s3 and mongo)
	createdSendMail, err := c.sendMailRepo.Create(ctx, sendMail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.s3Service.BulkDeleteFiles(context.Background(), uploadedKeys)
		return
	}

	sendMail.ID = createdSendMail.ID
	sendMail.CreatedAt = createdSendMail.CreatedAt
	sendMail.UpdatedAt = createdSendMail.UpdatedAt

	log.Debug().Interface("send_mail", rawMail).Msg("Publishing send mail message to queue")

	c.amqpService.PublishMessage("mail", "sent", map[string]interface{}{
		"send_mail_id": sendMail.ID.Hex(),
		"content":      rawMail, // Use the raw mail content for processing
	}, nil)

	ctx.JSON(http.StatusCreated, createdSendMail)
}
