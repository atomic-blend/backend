package habits

import (
	"net/http"
	"github.com/atomic-blend/backend/shared/middlewares/auth"

	"github.com/atomic-blend/backend/productivity/models"

	"github.com/gin-gonic/gin"
)

// AddHabitEntry adds a new entry to a habit
// @Summary Add habit entry
// @Description Add a new entry to a habit
// @Tags Habits
// @Accept json
// @Produce json
// @Param entry body models.HabitEntry true "Habit Entry"
// @Success 201 {object} models.HabitEntry
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /habits/entry/add [post]
func (c *HabitController) AddHabitEntry(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var entry models.HabitEntry
	if err := ctx.ShouldBindJSON(&entry); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set entry owner to authenticated user
	entry.UserID = authUser.UserID

	// Get the habit to check ownership
	habit, err := c.habitRepo.GetByID(ctx, entry.HabitID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if habit == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	// Check if the authenticated user owns this habit
	if habit.UserID != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to add entries to this habit"})
		return
	}

	// Create entry in database
	createdEntry, err := c.habitRepo.AddEntry(ctx, &entry)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdEntry)
}
