package tasks

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
	"net/http"

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
		// Check if all tags exist and belong to the user
		for _, tagID := range *task.Tags {
			tag, err := c.tagRepo.GetByID(ctx, tagID)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error validating tags: " + err.Error()})
				return
			}
			if tag == nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tag not found: " + tagID.Hex()})
				return
			}
			// Make sure the tag belongs to the user
			if tag.UserID == nil || *tag.UserID != authUser.UserID {
				ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to use this tag: " + tagID.Hex()})
				return
			}
		}
	}

	createdTask, err := c.taskRepo.Create(ctx, &task)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdTask)
}
