package auth

import (
	"atomic_blend_api/controllers/auth"
	"atomic_blend_api/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupRoutes configures all auth-related routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	userRepo := repositories.NewUserRepository(database)
	authController := auth.NewController(userRepo)

	// JWKS endpoint for token verification
	router.GET("/.well-known/jwks.json", authController.GetJWKS)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)

		
	}
}
