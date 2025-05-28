package time_entry

import (
	"atomic_blend_api/auth"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetByID retrieves a specific time entry by ID
func (tc *TimeEntryController) GetByID(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if authUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	idParam := c.Param("id")
	ctx := context.Background()

	timeEntry, err := tc.timeEntryRepository.GetByID(ctx, idParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Time entry not found"})
		return
	}

	// Verify user owns this time entry
	if timeEntry.User == nil || *timeEntry.User != authUser.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, timeEntry)
}