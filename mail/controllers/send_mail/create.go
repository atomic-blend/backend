package send_mail

import (
	"net/http"

	"github.com/atomic-blend/backend/mail/auth"
	"github.com/atomic-blend/backend/mail/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	var req CreateSendMailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the mail belongs to the authenticated user
	if req.Mail.UserID == primitive.NilObjectID {
		req.Mail.UserID = authUser.UserID
	} else if req.Mail.UserID != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Create send mail entity
	sendMail := &models.SendMail{
		Mail:         req.Mail,
		SendStatus:   models.SendStatusPending,
		RetryCounter: nil, // Will be managed by worker
		Trashed:      false,
	}

	// Save to database
	createdSendMail, err := c.sendMailRepo.Create(ctx, sendMail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdSendMail)
}
