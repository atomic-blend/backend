package tasks

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateTask updates an existing task
// @Summary Update task
// @Description Update an existing task
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body models.TaskEntity true "Task"
// @Success 200 {object} models.TaskEntity
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tasks/{id} [put]
func (c *TaskController) UpdateTask(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Get task ID from URL
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	// First get the task to check ownership
	existingTask, err := c.taskRepo.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if existingTask == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Check if the authenticated user owns this task
	if existingTask.User != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this task"})
		return
	}

	// Bind the updated task data
	var task models.TaskEntity
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Make sure to preserve the owner and ID
	task.User = existingTask.User
	task.ID = existingTask.ID

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

	updatedTask, err := c.taskRepo.Update(ctx, id, &task)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedTask)
}
