package tags

import (
	"net/http"
	"atomic-blend/backend/productivity/auth"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetTagByID retrieves a tag by its ID
// @Summary Get tag by ID
// @Description Get a tag by its ID
// @Tags Tags
// @Produce json
// @Param id path string true "Tag ID"
// @Success 200 {object} models.Tag
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tags/{id} [get]
func (c *TagController) GetTagByID(ctx *gin.Context) {
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

	// Get the tag
	tag, err := c.tagRepo.GetByID(ctx, objID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if tag exists
	if tag == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	// Check if the authenticated user owns this tag (if tag has user_id field)
	if tag.UserID != nil && *tag.UserID != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this tag"})
		return
	}

	ctx.JSON(http.StatusOK, tag)
}
