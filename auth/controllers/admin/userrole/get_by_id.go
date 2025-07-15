package userrole

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetRoleByID retrieves a user role by ID
// @Summary Get user role by ID
// @Description Get a user role by its ID
// @Tags User Roles
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} models.UserRoleEntity
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/user-roles/{id} [get]
func (c *Controller) GetRoleByID(ctx *gin.Context) {
	id := ctx.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	role, err := c.userRoleRepo.GetByID(ctx, objID)
	if err != nil {
		if err.Error() == "user role not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User role not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, role)
}
