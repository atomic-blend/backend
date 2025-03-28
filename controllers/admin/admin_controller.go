package admin

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/controllers/admin/userrole"
	"atomic_blend_api/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller is a controller for admin-related actions
type Controller struct {
}

// NewAdminController creates a new admin controller
func NewAdminController() *Controller {
	return &Controller{}
}

// SetupRoutes configures all admin-related routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	adminRoutes := router.Group("/admin")
	auth.RequireRoleMiddleware(adminRoutes, "admin")
	{
		userRoleRepo := repositories.NewUserRoleRepository(database)
		userRoleController := userrole.NewUserRoleController(userRoleRepo)
		userRoleController.SetupRoutes(adminRoutes)
	}
}
