package userrole

import (
	"auth/repositories"

	"github.com/gin-gonic/gin"
)

// Controller handles user role related operations
type Controller struct {
	userRoleRepo repositories.UserRoleRepositoryInterface
}

// NewUserRoleController creates a new user role controller instance
func NewUserRoleController(userRoleRepo repositories.UserRoleRepositoryInterface) *Controller {
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
