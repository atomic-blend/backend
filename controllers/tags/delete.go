package tags

import (
	"atomic_blend_api/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteTag deletes a tag by ID
// @Summary Delete tag
// @Description Delete a tag by ID
// @Tags Tags
// @Param id path string true "Tag ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tags/{id} [delete]
func (c *TagController) DeleteTag(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Get tag ID from URL
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tag ID is required"})
		return
	}

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID format"})
		return
	}

	// First get the tag to check ownership
	existingTag, err := c.tagRepo.GetByID(ctx, objID)
	if err != nil || existingTag == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	// Check if the authenticated user owns this tag
	if existingTag.UserID != nil && *existingTag.UserID != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this tag"})
		return
	}

	// Get all tasks for the user
	tasks, err := c.taskRepo.GetAll(ctx, &authUser.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving tasks: " + err.Error()})
		return
	}

	// Remove the tag from all tasks of this user
	for _, task := range tasks {
		if task.Tags != nil && len(*task.Tags) > 0 {
			// Check if the task contains the tag to be deleted
			updatedTags := removeTagFromSlice(*task.Tags, objID)
			if len(updatedTags) != len(*task.Tags) {
				// Tag was found and removed
				task.Tags = &updatedTags
				// Update the task in the database
				_, err := c.taskRepo.Update(ctx, task.ID, task)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating task: " + err.Error()})
					return
				}
			}
		}
	}

	// Finally delete the tag itself
	err = c.tagRepo.Delete(ctx, objID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}
