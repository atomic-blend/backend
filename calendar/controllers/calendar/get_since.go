package calendar

import (
	"net/http"
	"time"

	"github.com/atomic-blend/backend/calendar/models"
	"github.com/atomic-blend/backend/shared/middlewares/auth"

	"github.com/gin-gonic/gin"
)

// GetCalendarsSince retrieves calendars updated since a specific date
func (c *Controller) GetCalendarsSince(ctx *gin.Context) {
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	sinceStr := ctx.Query("since")
	if sinceStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameter: since"})
		return
	}

	sinceTime, err := time.Parse(time.RFC3339, sinceStr)
	if err != nil {
		sinceTime, err = time.Parse("2006-01-02T15:04:05-07:00", sinceStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Expected ISO8601"})
			return
		}
	}

	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	calendars, totalCount, err := c.calendarRepo.GetSince(ctx, authUser.UserID, sinceTime, int64(page), int64(size))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (totalCount + int64(size) - 1) / int64(size)

	if calendars == nil {
		calendars = make([]*models.Calendar, 0)
	}

	response := PaginatedCalendarResponse{
		Calendars:  calendars,
		TotalCount: totalCount,
		Page:       int64(page),
		Size:       int64(size),
		TotalPages: totalPages,
	}

	ctx.JSON(http.StatusOK, response)
}
