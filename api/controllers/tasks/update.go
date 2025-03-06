package tasks

import (
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
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tasks/{id} [put]
func (c *TaskController) UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")

	// Check if task exists
	existingTask, err := c.taskRepo.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingTask == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Bind updated task data
	var updatedTask models.TaskEntity
	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the task
	result, err := c.taskRepo.Update(ctx, id, &updatedTask)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
