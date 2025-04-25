package tags

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateTag updates an existing tag
// @Summary Update tag
// @Description Update an existing tag
// @Tags Tags
// @Accept json
// @Produce json
// @Param id path string true "Tag ID"
// @Param tag body models.Tag true "Tag"
// @Success 200 {object} models.Tag
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tags/{id} [put]
func (c *TagController) UpdateTag(ctx *gin.Context) {
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
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingTag == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	// Check if the authenticated user owns this tag
	if existingTag.UserID != nil && *existingTag.UserID != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this tag"})
		return
	}

	// Parse the updated tag from request body
	var updatedTag models.Tag
	if err := ctx.ShouldBindJSON(&updatedTag); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if updatedTag.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	// check if user is allowed to update the tag
	if updatedTag.UserID != nil && *updatedTag.UserID != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this tag"})
		return
	}

	// Preserve important fields from existing tag
	updatedTag.ID = existingTag.ID
	updatedTag.UserID = existingTag.UserID
	updatedTag.CreatedAt = existingTag.CreatedAt

	// Update timestamp
	now := primitive.NewDateTimeFromTime(time.Now())
	updatedTag.UpdatedAt = &now

	// Save the updated tag
	result, err := c.tagRepo.Update(ctx, &updatedTag)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
