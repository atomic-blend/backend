package sendmail

import (
	"net/http"

	"github.com/atomic-blend/backend/shared/middlewares/auth"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetSendMailByID retrieves a send mail by its ID for the authenticated user
// @Summary Get send mail by ID
// @Description Get a specific send mail by its ID for the authenticated user
// @Tags SendMail
// @Produce json
// @Param id path string true "Send Mail ID"
// @Success 200 {object} models.SendMail
// @Failure 401 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /mail/send/{id} [get]
func (c *Controller) GetSendMailByID(ctx *gin.Context) {
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

	// Get send mail by ID
	sendMail, err := c.sendMailRepo.GetByID(ctx, sendMailID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if sendMail == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Send mail not found"})
		return
	}

	// Check if the send mail belongs to the authenticated user
	if sendMail.Mail != nil && sendMail.Mail.UserID != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	ctx.JSON(http.StatusOK, sendMail)
}
