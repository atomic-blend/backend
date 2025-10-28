// Package payment contains controllers and routes for payment-related actions
package payment

import (
	"os"

	"github.com/atomic-blend/backend/auth/utils/stripe"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/atomic-blend/backend/shared/repositories/user"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller is a controller for payment-related actions
type Controller struct {
	userRepo      user.Interface
	stripeService stripe.Interface
}

// NewController creates a new instance of the payment controller
func NewController(stripeService stripe.Interface, userRepo user.Interface) *Controller {
	return &Controller{
		stripeService: stripeService,
		userRepo:      userRepo,
	}
}

// SetupRoutes configures all payment-related routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	userRepo := user.NewUserRepository(database)
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		log.Error().Msg("stripe secret key not found, skipping setup routes for payment controller")
		return
	}
	stripeService := stripe.NewStripeService(userRepo, &stripeKey)
	paymentController := NewController(stripeService, userRepo)

	paymentGroup := router.Group("/payment")

	protectedPaymentRoutes := auth.RequireAuth(paymentGroup)
	{
		protectedPaymentRoutes.POST("checkout", paymentController.Checkout)
	}
}
