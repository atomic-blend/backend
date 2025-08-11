package send_mail

import (
	"net/http"

	"github.com/atomic-blend/backend/mail/auth"
	"github.com/atomic-blend/backend/mail/models"

	"github.com/gin-gonic/gin"
	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateSendMailStatusRequest represents the request payload for updating send mail status
type UpdateSendMailStatusRequest struct {
	Status models.SendStatus `json:"status" binding:"required"`
}

// UpdateSendMailStatus updates the status of a send mail
// @Summary Update send mail status
// @Description Update the status of a send mail (pending, sent, failed)
// @Tags SendMail
// @Accept json
// @Produce json
// @Param id path string true "Send Mail ID"
// @Param body body UpdateSendMailStatusRequest true "Status update data"
// @Success 200 {object} models.SendMail
// @Failure 401 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /mail/send/{id}/status [patch]
func (c *Controller) UpdateSendMailStatus(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Parse send mail ID from URL parameter
	sendMailIDStr := ctx.Param("id")
	sendMailID, err := primitive.ObjectIDFromHex(sendMailIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid send mail ID"})
		return
	}

	// Bind JSON payload
	var req UpdateSendMailStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	switch req.Status {
	case models.SendStatusPending, models.SendStatusSent, models.SendStatusFailed:
		// Valid status
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be one of: pending, sent, failed"})
		return
	}

	// First get the send mail to check ownership
	existingSendMail, err := c.sendMailRepo.GetByID(ctx, sendMailID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingSendMail == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Send mail not found"})
		return
	}

	// Check if the send mail belongs to the authenticated user
	if existingSendMail.Mail != nil && existingSendMail.Mail.UserID != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Update status
	update := bson.M{
		"send_status": req.Status,
	}
	updatedSendMail, err := c.sendMailRepo.Update(ctx, sendMailID, update)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if updatedSendMail == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Send mail not found"})
		return
	}

	ctx.JSON(http.StatusOK, updatedSendMail)
}
