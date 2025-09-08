package mail

import (
	"net/http"
	"time"

	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PutActionsPayload struct {
	Read   []string `json:"read,omitempty"`
	Unread []string `json:"unread,omitempty"`
}

func (c *Controller) PutMailActions(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var payload PutActionsPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Process "read" actions
	for _, idStr := range payload.Read {
		mailID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			continue // Skip invalid IDs
		}

		mail, err := c.mailRepo.GetByID(ctx, mailID)
		if err != nil || mail == nil || mail.UserID != authUser.UserID {
			continue // Skip if mail not found or doesn't belong to user
		}

		read := true
		mail.Read = &read
		now := primitive.NewDateTimeFromTime(time.Now())
		mail.UpdatedAt = &now

		if err := c.mailRepo.Update(ctx, mail); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update mail status"})
			return
		}
	}

	// Process "unread" actions
	for _, idStr := range payload.Unread {
		mailID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			continue // Skip invalid IDs
		}

		mail, err := c.mailRepo.GetByID(ctx, mailID)
		if err != nil || mail == nil || mail.UserID != authUser.UserID {
			continue // Skip if mail not found or doesn't belong to user
		}

		read := false
		mail.Read = &read
		now := primitive.NewDateTimeFromTime(time.Now())
		mail.UpdatedAt = &now

		if err := c.mailRepo.Update(ctx, mail); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update mail status"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Mail actions updated successfully"})
}
