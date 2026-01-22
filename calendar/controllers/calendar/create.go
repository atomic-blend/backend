package calendar

import (
	"net/http"

	"github.com/atomic-blend/backend/calendar/models"
	"github.com/atomic-blend/backend/shared/middlewares/auth"

	"github.com/gin-gonic/gin"
)

// CreateCalendarRequest is the request payload for creating a calendar
type CreateCalendarRequest struct {
	Calendar *models.Calendar `json:"calendar" binding:"required"`
}

// CreateCalendar creates a new calendar
func (c *Controller) CreateCalendar(ctx *gin.Context) {
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var req CreateCalendarRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cal := req.Calendar
	if cal == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "calendar required"})
		return
	}

	// Set the authenticated user as owner
	cal.UserID = &authUser.UserID

	created, err := c.calendarRepo.Create(ctx, cal)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, created)
}
