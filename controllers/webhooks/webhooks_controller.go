package webhooks

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/repositories"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type WebhooksController struct {
	userRepo repositories.UserRepositoryInterface
}

func NewWebhooksController(userRepo repositories.UserRepositoryInterface) *WebhooksController {
	return &WebhooksController{
		userRepo: userRepo,
	}
}

// SetupRoutes configures all user-related routes
func SetupRoutes(router *gin.Engine, db *mongo.Database) {
	userRepo := repositories.NewUserRepository(db)
	webhooksController := NewWebhooksController(userRepo)

	// Public user routes (if any)
	revenueCatGroup := router.Group("/webhooks/revenuecat")

	// Protected user routes (require authentication with static token from env)
	bearerToken := "Bearer " + os.Getenv("REVENUE_CAT_WEBHOOK_TOKEN")
	protectedUserRoutes := auth.RequireStaticStringMiddleware(revenueCatGroup, bearerToken)
	{
		protectedUserRoutes.POST("", webhooksController.HandleRevenueCatWebhook)
	}

	// add other groups with static security for webhooks as needed
}
