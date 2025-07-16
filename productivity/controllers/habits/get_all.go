package habits

import (
	"net/http"
	"github.com/atomic-blend/backend/productivity/auth"
	"github.com/atomic-blend/backend/productivity/models"

	"github.com/gin-gonic/gin"
)

// GetAllHabits retrieves all habits for the authenticated user
// @Summary Get all habits
// @Description Get all habits for the authenticated user
// @Tags Habits
// @Produce json
// @Success 200 {array} models.Habit
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /habits [get]
func (c *HabitController) GetAllHabits(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Only get habits for the authenticated user
	habits, err := c.habitRepo.GetAll(ctx, &authUser.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ensure habits is never null (return empty array instead)
	if habits == nil {
		habits = []*models.Habit{}
	}

	// Load entries for each habit
	for _, habit := range habits {
		entries, err := c.habitRepo.GetEntriesByHabitID(ctx, habit.ID)
		if err == nil { // Only assign if no error
			habit.Entries = entries
		}
	}

	ctx.JSON(http.StatusOK, habits)
}
