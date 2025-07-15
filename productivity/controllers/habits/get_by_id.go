package habits

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetHabitByID retrieves a habit by its ID
// @Summary Get habit by ID
// @Description Get a habit by its ID
// @Tags Habits
// @Produce json
// @Param id path string true "Habit ID"
// @Success 200 {object} models.Habit
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /habits/{id} [get]
func (c *HabitController) GetHabitByID(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Get habit ID from URL
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Habit ID is required"})
		return
	}

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid habit ID format"})
		return
	}

	// Get the habit
	habit, err := c.habitRepo.GetByID(ctx, objID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if habit exists
	if habit == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	// Check if the authenticated user owns this habit
	if habit.UserID != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this habit"})
		return
	}

	// Get entries for this habit
	entries, err := c.habitRepo.GetEntriesByHabitID(ctx, habit.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load habit entries: " + err.Error()})
		return
	}

	// Assign entries to the habit
	habit.Entries = entries

	ctx.JSON(http.StatusOK, habit)
}
