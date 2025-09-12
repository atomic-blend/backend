package mail

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PutActionsPayload represents the request payload for updating mail actions
type PutActionsPayload struct {
	Read       []string `json:"read,omitempty"`
	Unread     []string `json:"unread,omitempty"`
	Archived   []string `json:"archived,omitempty"`
	Unarchived []string `json:"unarchived,omitempty"`
}

// PutMailActions updates the actions of a mail
func (c *Controller) PutMailActions(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Debug: Log request headers and content type
	log.Debug().Str("content-type", ctx.GetHeader("Content-Type")).Msg("Request Content-Type")
	log.Debug().Interface("headers", ctx.Request.Header).Msg("Request Headers")

	// Debug: Read raw body
	body, err := ctx.GetRawData()
	if err != nil {
		log.Error().Err(err).Msg("Failed to read raw body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	log.Debug().Str("raw-body", string(body)).Msg("Raw Request Body")

	// Reset the body for JSON binding
	ctx.Request.Body = io.NopCloser(strings.NewReader(string(body)))

	var payload PutActionsPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		log.Error().Err(err).Msg("JSON binding failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	log.Debug().Interface("payload", payload).Msg("Parsed Payload")

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
		ids = payload.Archived
	} else {
		ids = payload.Unarchived
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
