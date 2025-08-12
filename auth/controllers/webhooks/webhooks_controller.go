package webhooks

import (
	"os"

	staticstringmiddleware "github.com/atomic-blend/backend/shared/middlewares/static_string"
	userrepo "github.com/atomic-blend/backend/shared/repositories/user"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller handles webhooks from external services like RevenueCat
type Controller struct {
	userRepo userrepo.Interface
}

// NewWebhooksController creates a new instance of WebhooksController
func NewWebhooksController(userRepo userrepo.Interface) *Controller {
	return &Controller{
		userRepo: userRepo,
	}
}

// SetupRoutes configures all user-related routes
func SetupRoutes(router *gin.Engine, db *mongo.Database) {
	userRepo := userrepo.NewUserRepository(db)
	webhooksController := NewWebhooksController(userRepo)

	// Public user routes (if any)
	revenueCatGroup := router.Group("/webhooks/revenuecat")

	// Protected user routes (require authentication with static token from env)
	bearerToken := "Bearer " + os.Getenv("REVENUE_CAT_WEBHOOK_TOKEN")
	protectedUserRoutes := staticstringmiddleware.RequireStaticStringMiddleware(revenueCatGroup, bearerToken)
	{
		protectedUserRoutes.POST("", webhooksController.HandleRevenueCatWebhook)
	}

	// add other groups with static security for webhooks as needed
}
