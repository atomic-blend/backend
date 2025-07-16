package habits

import (
	"net/http"
	"github.com/atomic-blend/backend/productivity/auth"
	"github.com/atomic-blend/backend/productivity/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateHabit creates a new habit
// @Summary Create habit
// @Description Create a new habit
// @Tags Habits
// @Accept json
// @Produce json
// @Param habit body models.Habit true "Habit"
// @Success 201 {object} models.Habit
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /habits [post]
func (c *HabitController) CreateHabit(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	userHabits, err := c.habitRepo.GetAll(ctx, &authUser.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if user has reached habit limit
	if len(userHabits) >= 3 {
		if authUser.Claims.UserID != nil && !*authUser.Claims.IsSubscribed {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You must be subscribed to create more than 3 habits"})
			return
		}
	}

	var habit models.Habit
	if err := ctx.ShouldBindJSON(&habit); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set habit owner to authenticated user
	habit.UserID = authUser.UserID

	// Generate ID if not provided
	if habit.ID.IsZero() {
		habit.ID = primitive.NewObjectID()
	}

	// Set timestamps
	now := time.Now().Format(time.RFC3339)
	habit.CreatedAt = &now
	habit.UpdatedAt = &now

	// Create habit in database
	createdHabit, err := c.habitRepo.Create(ctx, &habit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdHabit)
}
