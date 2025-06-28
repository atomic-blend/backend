package tasks

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BulkUpdateTasks updates multiple tasks in a single request
// @Summary Bulk update tasks
// @Description Update multiple tasks, handling conflicts when database version is more recent
// @Tags Tasks
// @Accept json
// @Produce json
// @Param tasks body models.BulkTaskRequest true "Tasks to update"
// @Success 200 {object} models.BulkTaskResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tasks/bulk [put]
func (c *TaskController) BulkUpdateTasks(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Bind the request payload
	var request models.BulkTaskRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate that tasks are provided
	if len(request.Tasks) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "At least one task is required"})
		return
	}

	// Validate ownership and tags for all tasks
	for i, task := range request.Tasks {
		// Validate task ID is provided
		if task.ID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Task ID is required for bulk update",
				"index": i,
			})
			return
		}

		// Set the user ID to ensure ownership
		task.User = authUser.UserID

		// Validate tags if any are provided
		if task.Tags != nil && len(*task.Tags) > 0 {
			var validatedTags []*models.Tag

			// Check if all tags exist and belong to the user
			for _, tag := range *task.Tags {
				if tag.ID == nil {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"error": "Tag ID is required",
						"index": i,
					})
					return
				}			// Fetch the tag from the database to verify it exists
			dbTag, err := c.tagRepo.GetByID(ctx, *tag.ID)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": "Error validating tags: " + err.Error(),
						"index": i,
					})
					return
				}
				if dbTag == nil {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"error": "Tag not found: " + tag.ID.Hex(),
						"index": i,
					})
					return
				}

				// Make sure the tag belongs to the user
				if dbTag.UserID == nil || *dbTag.UserID != authUser.UserID {
					ctx.JSON(http.StatusForbidden, gin.H{
						"error": "You don't have permission to use this tag: " + tag.ID.Hex(),
						"index": i,
					})
					return
				}

				// Add the validated tag from the database
				validatedTags = append(validatedTags, dbTag)
			}

			// Replace the tags with validated tags from the database
			task.Tags = &validatedTags
		}
	}

	// Perform bulk update
	updated, conflicts, err := c.taskRepo.BulkUpdate(ctx, request.Tasks)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the response
	response := models.BulkTaskResponse{
		Updated:   updated,
		Conflicts: conflicts,
	}

	ctx.JSON(http.StatusOK, response)
}
