package tasks

import (
	"net/http"
	"atomic-blend/backend/productivity/auth"
	"atomic-blend/backend/productivity/models"

	"github.com/gin-gonic/gin"
)

// CreateTask creates a new task
// @Summary Create task
// @Description Create a new task
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task body models.TaskEntity true "Task"
// @Success 201 {object} models.TaskEntity
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tasks [post]
func (c *TaskController) CreateTask(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var task models.TaskEntity
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set task owner to authenticated user
	task.User = authUser.UserID

	// Set default values if needed
	if task.Completed == nil {
		completed := false
		task.Completed = &completed
	}

	// Validate tags if any are provided
	if task.Tags != nil && len(*task.Tags) > 0 {
		var validatedTags []*models.Tag

		// Check if all tags exist and belong to the user
		for _, tag := range *task.Tags {
			if tag.ID == nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tag ID is required"})
				return
			}

			// Fetch the tag from the database to verify it exists
			dbTag, err := c.tagRepo.GetByID(ctx, *tag.ID)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error validating tags: " + err.Error()})
				return
			}
			if dbTag == nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tag not found: " + tag.ID.Hex()})
				return
			}

			// Make sure the tag belongs to the user
			if dbTag.UserID == nil || *dbTag.UserID != authUser.UserID {
				ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to use this tag: " + tag.ID.Hex()})
				return
			}

			// Add the validated tag from the database
			validatedTags = append(validatedTags, dbTag)
		}

		// Replace the tags with validated tags from the database
		task.Tags = &validatedTags
	}

	createdTask, err := c.taskRepo.Create(ctx, &task)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdTask)
}
