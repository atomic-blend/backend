package draftmail

import (
	"net/http"

	"github.com/atomic-blend/backend/shared/middlewares/auth"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteDraftMail soft deletes a draft mail by marking it as trashed
// @Summary Delete draft mail
// @Description Soft delete a draft mail by marking it as trashed
// @Tags DraftMail
// @Produce json
// @Param id path string true "Draft Mail ID"
// @Success 204 "No Content"
// @Failure 401 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /mail/draft/{id} [delete]
func (c *Controller) DeleteDraftMail(ctx *gin.Context) {
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

	if existingDraftMail == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Draft mail not found"})
		return
	}

	// Check if the draft mail belongs to the authenticated user
	if existingDraftMail.Mail != nil && existingDraftMail.Mail.UserID != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Soft delete (mark as trashed)
	err = c.draftMailRepo.Delete(ctx, draftMailID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
