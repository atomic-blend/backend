package timeentrycontroller

import (
	"context"
	"net/http"
	"productivity/auth"

	"github.com/gin-gonic/gin"
)

// Delete deletes a time entry
func (tc *Controller) Delete(c *gin.Context) {
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

	err = tc.timeEntryRepository.Delete(ctx, idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete time entry"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Time entry deleted successfully"})
}
