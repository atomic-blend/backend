package mail

import (
	"net/http"

	"github.com/atomic-blend/backend/mail/auth"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetMailByID retrieves a mail by its ID for the authenticated user
// @Summary Get mail by ID
// @Description Get a specific mail by its ID for the authenticated user
// @Tags Mail
// @Produce json
// @Param id path string true "Mail ID"
// @Success 200 {object} models.Mail
// @Failure 401 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /mail/{id} [get]
func (c *MailController) GetMailByID(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Parse mail ID from URL parameter
	mailIDStr := ctx.Param("id")
	mailID, err := primitive.ObjectIDFromHex(mailIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mail ID"})
		return
	}

	// Get mail by ID
	mail, err := c.mailRepo.GetByID(ctx, mailID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if mail == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Mail not found"})
		return
	}

	// Check if the mail belongs to the authenticated user
	if mail.UserID != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	ctx.JSON(http.StatusOK, mail)
}
