package userrole

import (
	"auth/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateRole creates a new user role
// @Summary Create user role
// @Description Create a new user role
// @Tags User Roles
// @Accept json
// @Produce json
// @Param role body models.UserRoleEntity true "User Role"
// @Success 201 {object} models.UserRoleEntity
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/user-roles [post]
func (c *Controller) CreateRole(ctx *gin.Context) {
	var role models.UserRoleEntity

	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if role with the same name already exists
	existingRole, err := c.userRoleRepo.GetByName(ctx, role.Name)
	if err == nil && existingRole != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "A role with this name already exists"})
		return
	}

	createdRole, err := c.userRoleRepo.Create(ctx, &role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdRole)
}
