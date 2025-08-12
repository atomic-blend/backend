package admin

import (
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/atomic-blend/backend/auth/controllers/admin/userrole"
	userrolerepo "github.com/atomic-blend/backend/shared/repositories/user_role"

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
		userRoleRepo := userrolerepo.NewUserRoleRepository(database)
		userRoleController := userrole.NewUserRoleController(userRoleRepo)
		userRoleController.SetupRoutes(adminRoutes)
	}
}
