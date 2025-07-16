package habits

import (
	"net/http"
	"github.com/atomic-blend/backend/productivity/auth"
	"github.com/atomic-blend/backend/productivity/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EditHabitEntry updates an existing habit entry
// @Summary Edit habit entry
// @Description Update an existing habit entry
// @Tags Habits
// @Accept json
// @Produce json
// @Param id path string true "Entry ID"
// @Param entry body models.HabitEntry true "Updated Habit Entry"
// @Success 200 {object} models.HabitEntry
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /habits/entry/edit/{id} [put]
func (c *HabitController) EditHabitEntry(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Get entry ID from URL
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Entry ID is required"})
		return
	}

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entry ID format"})
		return
	}

	var updatedEntry models.HabitEntry
	if err := ctx.ShouldBindJSON(&updatedEntry); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure we're updating the correct entry
	updatedEntry.ID = objID

	// Get the habit to check ownership
	habit, err := c.habitRepo.GetByID(ctx, updatedEntry.HabitID)
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
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to edit entries for this habit"})
		return
	}

	// Set user ID to authenticated user
	updatedEntry.UserID = authUser.UserID

	// Update entry in database
	result, err := c.habitRepo.UpdateEntry(ctx, &updatedEntry)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
