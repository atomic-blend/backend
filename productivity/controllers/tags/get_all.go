package tags

import (
	"net/http"
	"github.com/atomic-blend/backend/productivity/auth"
	"github.com/atomic-blend/backend/productivity/models"

	"github.com/gin-gonic/gin"
)

// GetAllTags retrieves all tags for the authenticated user
// @Summary Get all tags
// @Description Get all tags for the authenticated user
// @Tags Tags
// @Produce json
// @Success 200 {array} models.Tag
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tags [get]
func (c *TagController) GetAllTags(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Only get tags for the authenticated user
	tags, err := c.tagRepo.GetAll(ctx, &authUser.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ensure tags is never null (return empty array instead)
	if tags == nil {
		tags = []*models.Tag{}
	}

	ctx.JSON(http.StatusOK, tags)
}
