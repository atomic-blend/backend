package stripe

import (
	"github.com/atomic-blend/backend/shared/repositories/user"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v83"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StripeServiceInferface interface {
	GetOrCreateCustomer(ctx *gin.Context, userID primitive.ObjectID)
}

type StripeService struct {
	userService  *user.Repository
	stripeClient *stripe.Client
}

func NewStripeService(userRepo *user.Repository, stripeKey *string) StripeServiceInferface {
	sc := stripe.NewClient(*stripeKey)
	return &StripeService{
		userService:  userRepo,
		stripeClient: sc,
	}
}

func (s *StripeService) GetOrCreateCustomer(ctx *gin.Context, userID primitive.ObjectID) {
	user, err := s.userService.FindByID(ctx, userID)
	if err != nil {
		log.Error().Str("user_id", userID.Hex()).Msg("cannot find user")
		return
	}

	if user.StripeCustomerId == nil {
		// TODO: create stripe customer
	} else {
		//TODO: get customer and return it
	}
}
