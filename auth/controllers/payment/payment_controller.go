package payment

import (
	"os"

	"github.com/atomic-blend/backend/auth/utils/stripe"
	"github.com/atomic-blend/backend/shared/repositories/user"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type Controller struct {
	stripeService stripe.StripeServiceInferface
}

func NewController(stripeService stripe.StripeServiceInferface) *Controller {
	return &Controller{
		stripeService: stripeService,
	}
}

func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	userRepo := user.NewUserRepository(database)
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		log.Error().Msg("stripe secret key not found, skipping setup routes for payment controller")
		return
	}
	stripeService := stripe.NewStripeService(userRepo, &stripeKey)
	paymentController := NewController(stripeService)

	paymentGroup := router.Group("/payment")
	{
		paymentGroup.POST("subscribe", paymentController.Subscribe)
	}
}
