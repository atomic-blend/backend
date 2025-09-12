package mail

import (
	"net/http"
	"time"

	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PutActionsPayload represents the request payload for updating mail actions
type PutActionsPayload struct {
	Read      []string `json:"read,omitempty"`
	Unread    []string `json:"unread,omitempty"`
	Archive   []string `json:"archive,omitempty"`
	Unarchive []string `json:"unarchive,omitempty"`
}

// PutMailActions updates the actions of a mail
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
	_processRead(ctx, payload, c, authUser, true)

	// Process "unread" actions
	_processRead(ctx, payload, c, authUser, false)

	// Process "archive" actions
	_processArchive(ctx, payload, c, authUser, true)

	// Process "unarchive" actions
	_processArchive(ctx, payload, c, authUser, false)

	ctx.JSON(http.StatusOK, gin.H{"message": "Mail actions updated successfully"})
}

func _processRead(ctx *gin.Context, payload PutActionsPayload, c *Controller, authUser *auth.UserAuthInfo, read bool) {
	var ids []string
	if read {
		ids = payload.Read
	} else {
		ids = payload.Unread
	}

	for _, idStr := range ids {
		mailID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			continue // Skip invalid IDs
		}

		mail, err := c.mailRepo.GetByID(ctx, mailID)
		if err != nil || mail == nil || mail.UserID != authUser.UserID {
			continue // Skip if mail not found or doesn't belong to user
		}

		mail.Read = &read
		now := primitive.NewDateTimeFromTime(time.Now())
		mail.UpdatedAt = &now

		if err := c.mailRepo.Update(ctx, mail); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update mail status"})
			return
		}
	}
}

func _processArchive(ctx *gin.Context, payload PutActionsPayload, c *Controller, authUser *auth.UserAuthInfo, archived bool) {
	var ids []string
	if archived {
		ids = payload.Archive
	} else {
		ids = payload.Unarchive
	}

	for _, idStr := range ids {
		mailID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			continue // Skip invalid IDs
		}

		mail, err := c.mailRepo.GetByID(ctx, mailID)
		if err != nil || mail == nil || mail.UserID != authUser.UserID {
			continue // Skip if mail not found or doesn't belong to user
		}

		mail.Archived = &archived
		if archived {
			trashed := false
			mail.Trashed = &trashed
		}
		now := primitive.NewDateTimeFromTime(time.Now())
		mail.UpdatedAt = &now

		if err := c.mailRepo.Update(ctx, mail); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update mail status"})
			return
		}
	}
}
