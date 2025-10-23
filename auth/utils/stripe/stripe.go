package stripe

import (
	"context"

	"github.com/atomic-blend/backend/shared/repositories/user"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v83"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StripeServiceInferface interface {
	GetOrCreateCustomer(ctx *gin.Context, userID primitive.ObjectID) *stripe.Customer
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

func (s *StripeService) GetOrCreateCustomer(ctx *gin.Context, userID primitive.ObjectID) *stripe.Customer {
	user, err := s.userService.FindByID(ctx, userID)
	if err != nil {
		log.Error().Str("user_id", userID.Hex()).Msg("cannot find user")
		return nil
	}

	if user.StripeCustomerId == nil {
		// create stripe customer
		params := &stripe.CustomerCreateParams{
			Name:  stripe.String(*user.FirstName + " " + *user.LastName),
			Email: stripe.String(*user.Email),
		}
		result, err := s.stripeClient.V1Customers.Create(context.TODO(), params)
		if err != nil {
			log.Error().Err(err).Msg("error during creation of the stripe customer")
			return nil
		}
		return result
	} else {
		// get customer and return it
		result, err := s.stripeClient.V1Customers.Retrieve(context.TODO(), *user.StripeCustomerId, nil)
		if err != nil {
			log.Error().Err(err).Msg("cannot get stripe customer")
			return nil
		}
		return result
	}
}
