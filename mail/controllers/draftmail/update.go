package draftmail

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	userv1 "github.com/atomic-blend/backend/grpc/gen/user/v1"
	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateDraftMailRequest represents the request payload for updating a draft mail
type UpdateDraftMailRequest struct {
	Mail *models.Mail `json:"mail" binding:"required"`
}

// UpdateDraftMail updates a draft mail entry
// @Summary Update a draft mail
// @Description Update a draft mail entry with the provided mail data
// @Tags DraftMail
// @Accept json
// @Produce json
// @Param id path string true "Draft Mail ID"
// @Param body body UpdateDraftMailRequest true "Draft mail data"
// @Success 200 {object} models.SendMail
// @Failure 401 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /mail/draft/{id} [put]
func (c *Controller) UpdateDraftMail(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Parse draft mail ID from URL parameter
	draftMailIDStr := ctx.Param("id")
	draftMailID, err := primitive.ObjectIDFromHex(draftMailIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid draft mail ID"})
		return
	}

	// First get the draft mail to check ownership
	existingDraftMail, err := c.draftMailRepo.GetByID(ctx, draftMailID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Debug().Interface("draftMailID", draftMailID).Msg("Draft mail ID")
	log.Debug().Interface("draftMailIDStr", draftMailIDStr).Msg("Draft mail ID string")
	log.Debug().Interface("err", err).Msg("Error")
	log.Debug().Interface("existingDraftMail", existingDraftMail).Msg("Existing draft mail")

	if existingDraftMail == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Draft mail not found"})
		return
	}

	log.Debug().Interface("existingDraftMail", existingDraftMail.Mail.UserID).Msg("Existing draft mail user ID")
	log.Debug().Interface("authUser", authUser.UserID).Msg("Auth user user ID")
	log.Debug().Interface("existingDraftMail.Mail.UserID == authUser.UserID", existingDraftMail.Mail.UserID == authUser.UserID).Msg("Existing draft mail user ID equals auth user user ID")

	// Check if the draft mail belongs to the authenticated user
	if existingDraftMail.Mail != nil && existingDraftMail.Mail.UserID != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
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

	log.Debug().Interface("raw_mail", rawMail).Msg("Received raw mail for draft update")

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

	// Transform the raw mail to an encrypted mail
	encryptedMail, err := rawMail.Encrypt(userPublicKey.Msg.PublicKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to encrypt mail")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt mail"})
		return
	}

	log.Debug().Interface("encrypted_mail", encryptedMail).Msg("Encrypted mail ready for draft update")

	mailEntity := encryptedMail.ToMailEntity()
	mailEntity.UserID = authUser.UserID
	// Prepare update data
	updateData := bson.M{
		"mail": mailEntity,
	}

	// Handle attachments if any
	if len(encryptedMail.Attachments) > 0 {
		// Generate encrypted attachments upload requests to send them to S3
		encryptedAttachments := make([]*awss3.PutObjectInput, 0)
		attachments := make([]models.MailAttachment, 0)

		for _, attachment := range encryptedMail.Attachments {
			uniqueFileID := uuid.New().String()
			payload, err := c.s3Service.GenerateUploadPayload(context.Background(), attachment.Data, "draft_mail/attachments/"+userPublicKey.Msg.UserId, uniqueFileID, map[string]string{})
			if err != nil {
				log.Error().Err(err).Msg("Failed to upload attachment to S3")
				continue
			}
			attachments = append(attachments, models.MailAttachment{
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

		// Update the mail entity with new attachments
		mailEntity := encryptedMail.ToMailEntity()
		mailEntity.Attachments = attachments
		updateData["mail"] = mailEntity

		// Update draft mail with rollback on failure (s3 and mongo)
		updatedDraftMail, err := c.draftMailRepo.Update(ctx, draftMailID, updateData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.s3Service.BulkDeleteFiles(context.Background(), uploadedKeys)
			return
		}

		ctx.JSON(http.StatusOK, updatedDraftMail)
		return
	}

	// Update draft mail without new attachments
	updatedDraftMail, err := c.draftMailRepo.Update(ctx, draftMailID, updateData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if updatedDraftMail == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Draft mail not found"})
		return
	}

	log.Debug().Interface("draft_mail", rawMail).Msg("Draft mail updated successfully")

	ctx.JSON(http.StatusOK, updatedDraftMail)
}
