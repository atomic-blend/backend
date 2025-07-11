package timeentrycontroller

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAll retrieves all time entries for the authenticated user
func (tc *Controller) GetAll(c *gin.Context) {
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

	if timeEntries == nil {
		timeEntries = []*models.TimeEntry{} // Ensure we return an empty array instead of null
	}

	c.JSON(http.StatusOK, timeEntries)
}
