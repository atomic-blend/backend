package send_mail

import (
	"net/http"

	"github.com/atomic-blend/backend/mail/auth"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteSendMail soft deletes a send mail by marking it as trashed
// @Summary Delete send mail
// @Description Soft delete a send mail by marking it as trashed
// @Tags SendMail
// @Produce json
// @Param id path string true "Send Mail ID"
// @Success 204 "No Content"
// @Failure 401 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /mail/send/{id} [delete]
func (c *Controller) DeleteSendMail(ctx *gin.Context) {
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

	// Soft delete (mark as trashed)
	err = c.sendMailRepo.Delete(ctx, sendMailID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
