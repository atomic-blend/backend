package habits

import (
	"net/http"
	"github.com/atomic-blend/backend/productivity/auth"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteHabit deletes a habit
// @Summary Delete habit
// @Description Delete a habit
// @Tags Habits
// @Param id path string true "Habit ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /habits/{id} [delete]
func (c *HabitController) DeleteHabit(ctx *gin.Context) {
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

	// First get the habit to check ownership
	habit, err := c.habitRepo.GetByID(ctx, objID)
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
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this habit"})
		return
	}

	err = c.habitRepo.Delete(ctx, objID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Habit deleted successfully"})
}
