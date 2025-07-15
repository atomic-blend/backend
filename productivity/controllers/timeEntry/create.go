package timeentrycontroller

import (
	"context"
	"net/http"
	"productivity/auth"
	"productivity/models"
	"time"

	"github.com/gin-gonic/gin"
)

// Create creates a new time entry
func (tc *Controller) Create(c *gin.Context) {
	// TODO: replace that with grpc call
	authUser := auth.GetAuthUser(c)
	if authUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var timeEntry models.TimeEntry
	if err := c.ShouldBindJSON(&timeEntry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set user ID and timestamps
	timeEntry.User = &authUser.UserID
	now := time.Now().Format(time.RFC3339)
	timeEntry.CreatedAt = now
	timeEntry.UpdatedAt = now

	ctx := context.Background()
	createdTimeEntry, err := tc.timeEntryRepository.Create(ctx, &timeEntry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create time entry"})
		return
	}

	c.JSON(http.StatusCreated, createdTimeEntry)
}
