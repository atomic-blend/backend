package tasks

import (
	"net/http"
	"strconv"

	"github.com/atomic-blend/backend/productivity/models"
	"github.com/atomic-blend/backend/shared/middlewares/auth"

	"github.com/gin-gonic/gin"
)

// PaginatedTaskResponse represents the paginated response for tasks
type PaginatedTaskResponse struct {
	Tasks      []*models.TaskEntity `json:"tasks"`
	TotalCount int64                `json:"total_count"`
	Page       int64                `json:"page"`
	Size       int64                `json:"size"`
	TotalPages int64                `json:"total_pages"`
}

// GetAllTasks retrieves tasks for the authenticated user with optional pagination
// @Summary Get all tasks
// @Description Get tasks for the authenticated user with optional pagination. If both page and limit are provided, returns paginated results with total count. If either is missing, returns all tasks.
// @Tags Tasks
// @Produce json
// @Param page query int false "Page number (1-based)"
// @Param limit query int false "Number of tasks per page"
// @Success 200 {object} PaginatedTaskResponse
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

	// Parse pagination parameters
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")

	var page, limit *int64
	var err error

	if pageStr != "" {
		pageVal, err := strconv.ParseInt(pageStr, 10, 64)
		if err != nil || pageVal < 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
			return
		}
		page = &pageVal
	}

	if limitStr != "" {
		limitVal, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil || limitVal < 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
		limit = &limitVal
	}

	// Get tasks with pagination
	tasks, totalCount, err := c.taskRepo.GetAll(ctx, &authUser.UserID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ensure tasks is never null (return empty array instead)
	if tasks == nil {
		tasks = []*models.TaskEntity{}
	}

	// Create response
	response := PaginatedTaskResponse{
		Tasks:      tasks,
		TotalCount: totalCount,
	}

	// Add pagination metadata if both page and limit were provided
	if page != nil && limit != nil {
		totalPages := (totalCount + *limit - 1) / *limit // Ceiling division
		response.Page = *page
		response.Size = *limit
		response.TotalPages = totalPages
	}

	ctx.JSON(http.StatusOK, response)
}
