package admin

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/controllers/admin/user_role"
	"atomic_blend_api/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminController struct {
}

func NewAdminController() *AdminController {
	return &AdminController{}
}

func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	adminRoutes := router.Group("/admin")
	auth.RequireRole(adminRoutes, "admin")
	{
		userRoleRepo := repositories.NewUserRoleRepository(database)
		userRoleController := user_role.NewUserRoleController(userRoleRepo)
		userRoleController.SetupRoutes(adminRoutes)
	}
}
