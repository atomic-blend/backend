package timeentrycontroller

import (
	"context"
	"net/http"
	"productivity/models"
	"time"

	"github.com/gin-gonic/gin"
)

// Update updates an existing time entry
func (tc *Controller) Update(c *gin.Context) {
	// TODO: replace that with grpc call
	authUser := auth.GetAuthUser(c)
	if authUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	idParam := c.Param("id")
	ctx := context.Background()

	// Verify time entry exists and user owns it
	existingEntry, err := tc.timeEntryRepository.GetByID(ctx, idParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Time entry not found"})
		return
	}

	if existingEntry.User == nil || *existingEntry.User != authUser.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var updateData models.TimeEntry
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Preserve original fields and update timestamp
	updateData.ID = existingEntry.ID
	updateData.User = existingEntry.User
	updateData.CreatedAt = existingEntry.CreatedAt
	updateData.UpdatedAt = time.Now().Format(time.RFC3339)

	updatedTimeEntry, err := tc.timeEntryRepository.Update(ctx, idParam, &updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update time entry"})
		return
	}

	c.JSON(http.StatusOK, updatedTimeEntry)
}
