package webhooks

import (
	"atomic-blend/backend/auth/auth"
	"atomic-blend/backend/auth/repositories"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller handles webhooks from external services like RevenueCat
type Controller struct {
	userRepo repositories.UserRepositoryInterface
}

// NewWebhooksController creates a new instance of WebhooksController
func NewWebhooksController(userRepo repositories.UserRepositoryInterface) *Controller {
	return &Controller{
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
