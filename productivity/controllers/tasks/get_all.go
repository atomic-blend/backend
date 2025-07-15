package tasks

import (
	"net/http"
	"productivity/auth"
	"productivity/models"

	"github.com/gin-gonic/gin"
)

// GetAllTasks retrieves all tasks for the authenticated user
// @Summary Get all tasks
// @Description Get all tasks for the authenticated user
// @Tags Tasks
// @Produce json
// @Success 200 {array} models.TaskEntity
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tasks [get]
func (c *TaskController) GetAllTasks(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Only get tasks for the authenticated user
	tasks, err := c.taskRepo.GetAll(ctx, &authUser.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ensure tasks is never null (return empty array instead)
	if tasks == nil {
		tasks = []*models.TaskEntity{} 
	}

	ctx.JSON(http.StatusOK, tasks)
}
