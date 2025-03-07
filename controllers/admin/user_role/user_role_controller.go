package user_role

import (
	"atomic_blend_api/repositories"

	"github.com/gin-gonic/gin"
)

// UserRoleController handles user role related operations
type UserRoleController struct {
	userRoleRepo repositories.UserRoleRepositoryInterface
}

// NewUserRoleController creates a new user role controller instance
func NewUserRoleController(userRoleRepo repositories.UserRoleRepositoryInterface) *UserRoleController {
	return &UserRoleController{
		userRoleRepo: userRoleRepo,
	}
}

// SetupRoutes sets up the user role routes
func (c *UserRoleController) SetupRoutes(router *gin.RouterGroup) {
	userRoleRoutes := router.Group("/user-roles")
	{
		userRoleRoutes.GET("", c.GetAllRoles)
		userRoleRoutes.GET("/:id", c.GetRoleByID)
		userRoleRoutes.POST("", c.CreateRole)
		userRoleRoutes.PUT("/:id", c.UpdateRole)
		userRoleRoutes.DELETE("/:id", c.DeleteRole)
	}
}
