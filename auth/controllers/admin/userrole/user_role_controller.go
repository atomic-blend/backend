package userrole

import (
	userrolerepo "github.com/atomic-blend/backend/shared/repositories/user_role"

	"github.com/gin-gonic/gin"
)

// Controller handles user role related operations
type Controller struct {
	userRoleRepo userrolerepo.Interface
}

// NewUserRoleController creates a new user role controller instance
func NewUserRoleController(userRoleRepo userrolerepo.Interface) *Controller {
	return &Controller{
		userRoleRepo: userRoleRepo,
	}
}

// SetupRoutes sets up the user role routes
func (c *Controller) SetupRoutes(router *gin.RouterGroup) {
	userRoleRoutes := router.Group("/user-roles")
	{
		userRoleRoutes.GET("", c.GetAllRoles)
		userRoleRoutes.GET("/:id", c.GetRoleByID)
		userRoleRoutes.POST("", c.CreateRole)
		userRoleRoutes.PUT("/:id", c.UpdateRole)
		userRoleRoutes.DELETE("/:id", c.DeleteRole)
	}
}
