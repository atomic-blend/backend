package tags

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateTag creates a new tag
// @Summary Create tag
// @Description Create a new tag
// @Tags Tags
// @Accept json
// @Produce json
// @Param tag body models.Tag true "Tag"
// @Success 201 {object} models.Tag
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tags [post]
func (c *TagController) CreateTag(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var tag models.Tag
	if err := ctx.ShouldBindJSON(&tag); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set tag owner to authenticated user
	tag.UserID = &authUser.UserID

	// Set timestamps
	now := primitive.NewDateTimeFromTime(time.Now())
	tag.CreatedAt = &now
	tag.UpdatedAt = &now

	createdTag, err := c.tagRepo.Create(ctx, &tag)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdTag)
}
