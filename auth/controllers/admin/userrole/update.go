package userrole

import (
	"atomic-blend/backend/auth/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateRole updates an existing user role
// @Summary Update user role
// @Description Update an existing user role
// @Tags User Roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param role body models.UserRoleEntity true "User Role"
// @Success 200 {object} models.UserRoleEntity
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/user-roles/{id} [put]
func (c *Controller) UpdateRole(ctx *gin.Context) {
	id := ctx.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Check if role exists
	_, err = c.userRoleRepo.GetByID(ctx, objID)
	if err != nil {
		if err.Error() == "user role not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User role not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updatedRole models.UserRoleEntity
	if err := ctx.ShouldBindJSON(&updatedRole); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the ID from the URL parameter
	updatedRole.ID = &objID

	// Update the role in the database
	result, err := c.userRoleRepo.Update(ctx, &updatedRole)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
