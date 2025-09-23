package draftmail

import (
	"net/http"

	"github.com/atomic-blend/backend/shared/middlewares/auth"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetDraftMailByID retrieves a draft mail by its ID for the authenticated user
// @Summary Get draft mail by ID
// @Description Get a specific draft mail by its ID for the authenticated user
// @Tags DraftMail
// @Produce json
// @Param id path string true "Draft Mail ID"
// @Success 200 {object} models.SendMail
// @Failure 401 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /mail/draft/{id} [get]
func (c *Controller) GetDraftMailByID(ctx *gin.Context) {
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

	// Get draft mail by ID
	draftMail, err := c.draftMailRepo.GetByID(ctx, draftMailID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if draftMail == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Draft mail not found"})
		return
	}

	// Check if the draft mail belongs to the authenticated user
	if draftMail.Mail != nil && draftMail.Mail.UserID != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	ctx.JSON(http.StatusOK, draftMail)
}
