package userrole

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllRoles retrieves all user roles
// @Summary Get all user roles
// @Description Get a list of all user roles
// @Tags User Roles
// @Produce json
// @Success 200 {array} models.UserRoleEntity
// @Failure 500 {object} map[string]interface{}
// @Router /admin/user-roles [get]
func (c *Controller) GetAllRoles(ctx *gin.Context) {
	roles, err := c.userRoleRepo.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, roles)
}
