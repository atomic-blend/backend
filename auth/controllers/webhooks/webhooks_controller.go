package webhooks

import (
	"os"

	"github.com/atomic-blend/backend/auth/utils/stripe"
	staticstringmiddleware "github.com/atomic-blend/backend/shared/middlewares/static_string"
	userrepo "github.com/atomic-blend/backend/shared/repositories/user"
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller handles webhooks from external services like RevenueCat
type Controller struct {
	userRepo    userrepo.Interface
	stripeService stripe.Interface
}

// NewWebhooksController creates a new instance of WebhooksController
func NewWebhooksController(userRepo userrepo.Interface, stripeService stripe.Interface) *Controller {
	return &Controller{
		stripeService: stripeService,
		userRepo: userRepo,
	}
}

// SetupRoutes configures all user-related routes
func SetupRoutes(router *gin.Engine, db *mongo.Database) {
	userRepo := userrepo.NewUserRepository(db)
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		log.Error().Msg("stripe secret key not found, skipping setup routes for payment controller")
		return
	}
	stripeService := stripe.NewStripeService(userRepo, &stripeKey)
	webhooksController := NewWebhooksController(userRepo, stripeService)

	// Public user routes (if any)
	revenueCatGroup := router.Group("/webhooks/revenuecat")

	// Protected user routes (require authentication with static token from env)
	bearerToken := "Bearer " + os.Getenv("REVENUE_CAT_WEBHOOK_TOKEN")
	if bearerToken != "Bearer " {
		rcRoutes := staticstringmiddleware.RequireStaticStringMiddleware(revenueCatGroup, bearerToken)
		{
			rcRoutes.POST("", webhooksController.HandleRevenueCatWebhook)
		}
	}

	// add other groups with static security for webhooks as needed
	stripeGroup := router.Group("/webhooks/stripe")
	{
		stripeGroup.POST("", webhooksController.HandleStripeWebhook)
	}
}
