package calendar

import (
	"net/http"

	"github.com/atomic-blend/backend/calendar/models"
	"github.com/atomic-blend/backend/shared/middlewares/auth"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateCalendarRequest is the payload for updating a calendar
type UpdateCalendarRequest struct {
	Calendar *models.Calendar `json:"calendar" binding:"required"`
}

// UpdateCalendar updates an existing calendar
func (c *Controller) UpdateCalendar(ctx *gin.Context) {
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	idStr := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid calendar ID"})
		return
	}

	existing, err := c.calendarRepo.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existing == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Calendar not found"})
		return
	}
	if existing.UserID != nil && *existing.UserID != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req UpdateCalendarRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Calendar == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "calendar required"})
		return
	}

	// Ensure owner preserved
	req.Calendar.UserID = &authUser.UserID

	updated, err := c.calendarRepo.Update(ctx, id, req.Calendar)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updated)
}
