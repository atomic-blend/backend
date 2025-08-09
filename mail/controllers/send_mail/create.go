package send_mail

import (
	"context"
	"net/http"
	"os"

	"connectrpc.com/connect"
	userv1 "github.com/atomic-blend/backend/grpc/gen/user/v1"
	"github.com/atomic-blend/backend/mail/auth"
	"github.com/atomic-blend/backend/mail/grpc/clients"
	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/mail/utils/amqp"
	"github.com/atomic-blend/backend/mail/utils/s3"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

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

	// get the user public key from the auth service via grpc
	log.Info().Msg("Instantiating user client")
	userClient, err := clients.NewUserClient()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user client")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user client"})
		return
	}

	//TODO: [DONE] support getting user public key from userID
	log.Info().Msg("Getting user public key")
	userPublicKey, err := userClient.GetUserPublicKey(context.Background(), &connect.Request[userv1.GetUserPublicKeyRequest]{
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

	s3Service, err := s3.NewS3Service(os.Getenv("AWS_BUCKET"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create S3 service")
		return
	}

	// Create send mail entity
	sendMail := &models.SendMail{
		Mail:       encryptedMail.ToMailEntity(),
		SendStatus: models.SendStatusPending,
		Trashed:    false,
	}

	// Generate encrypted attachments upload requests to send them to S3
	encryptedAttachments := make([]*awss3.PutObjectInput, 0)
	for _, attachment := range encryptedMail.Attachments {
		uniqueFileID := uuid.New().String()
		payload, err := s3Service.GenerateUploadPayload(context.Background(), attachment.Data, "send_mail/attachments/"+userPublicKey.Msg.UserId, uniqueFileID, map[string]string{})
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
	uploadedKeys, err := s3Service.BulkUploadFiles(context.Background(), encryptedAttachments)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upload attachments to S3")
		return
	}

	// TODO: [DONE] setup transactional upload for mail with rollback on failure (s3 and mongo)
	createdSendMail, err := c.sendMailRepo.Create(ctx, sendMail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		s3Service.BulkDeleteFiles(context.Background(), uploadedKeys)
		return
	}

	sendMail.ID = createdSendMail.ID
	sendMail.CreatedAt = createdSendMail.CreatedAt
	sendMail.UpdatedAt = createdSendMail.UpdatedAt

	// TODO: fix headers null in db mail
	//TODO: fix raw mail in payload have all the fields not defined

	log.Debug().Interface("send_mail", rawMail).Msg("Publishing send mail message to queue")

	//TODO: [DONE] publish to message queue
	amqp.PublishMessage("mail", "sent", map[string]interface{}{
		"send_mail_id": sendMail.ID.Hex(),
		"content":      rawMail, // Use the raw mail content for processing
	})

	ctx.JSON(http.StatusCreated, createdSendMail)
}
