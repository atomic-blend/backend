package send_mail

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	userv1 "github.com/atomic-blend/backend/grpc/gen/user/v1"
	"github.com/atomic-blend/backend/mail/auth"
	"github.com/atomic-blend/backend/mail/grpc/clients"
	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/mail/utils/amqp"
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
	var req models.RawMail
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	//TODO: support getting user public key from userID
	log.Info().Msg("Getting user public key")
	userPublicKey, err := userClient.GetUserPublicKey(context.Background(), &connect.Request[userv1.GetUserPublicKeyRequest]{
		Msg: &userv1.GetUserPublicKeyRequest{
			Email: string(authUser.UserID.Hex()),
		},
	})
	if err != nil {
		log.Info().Msg("User not found, skipping")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	log.Info().Str("public_key", userPublicKey.Msg.PublicKey).Msg("User public key retrieved successfully")

	// TODO: transform the raw mail to an encrypted mail
	encryptedMail := &models.Mail{}

	// Create send mail entity
	sendMail := &models.SendMail{
		Mail:       encryptedMail,
		SendStatus: models.SendStatusPending,
		Trashed:    false,
	}

	// TODO: setup transactional upload for mail with rollback on failure (s3 and mongo)
	createdSendMail, err := c.sendMailRepo.Create(ctx, sendMail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sendMail.ID = createdSendMail.ID
	sendMail.CreatedAt = createdSendMail.CreatedAt
	sendMail.UpdatedAt = createdSendMail.UpdatedAt

	//TODO: publish to message queue
	amqp.PublishMessage("mail", "mail:sent", map[string]interface{}{
		"send_mail_id": sendMail.ID.Hex(),
		"content":      sendMail.Mail,
	})

	ctx.JSON(http.StatusCreated, createdSendMail)
}
