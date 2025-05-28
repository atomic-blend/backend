package time_entry

import (
	"atomic_blend_api/auth"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAll retrieves all time entries for the authenticated user
func (tc *TimeEntryController) GetAll(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if authUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	ctx := context.Background()
	timeEntries, err := tc.timeEntryRepository.GetAll(ctx, &authUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve time entries"})
		return
	}

	c.JSON(http.StatusOK, timeEntries)
}