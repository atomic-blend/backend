package habits

import (
	"net/http"
	"github.com/atomic-blend/backend/productivity/auth"
	"github.com/atomic-blend/backend/productivity/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateHabit updates an existing habit
// @Summary Update habit
// @Description Update an existing habit
// @Tags Habits
// @Accept json
// @Produce json
// @Param id path string true "Habit ID"
// @Param habit body models.Habit true "Updated Habit"
// @Success 200 {object} models.Habit
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /habits/{id} [put]
func (c *HabitController) UpdateHabit(ctx *gin.Context) {
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
	existingHabit, err := c.habitRepo.GetByID(ctx, objID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingHabit == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	// Check if the authenticated user owns this habit
	if existingHabit.UserID != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this habit"})
		return
	}

	// Parse the updated habit from request body
	var updatedHabit models.Habit
	if err := ctx.ShouldBindJSON(&updatedHabit); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Preserve important fields from existing habit
	updatedHabit.ID = existingHabit.ID
	updatedHabit.UserID = existingHabit.UserID
	updatedHabit.CreatedAt = existingHabit.CreatedAt

	// Update timestamp
	now := time.Now().Format(time.RFC3339)
	updatedHabit.UpdatedAt = &now

	// Save the updated habit
	result, err := c.habitRepo.Update(ctx, &updatedHabit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
