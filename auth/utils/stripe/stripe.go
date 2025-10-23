package stripe

import (
	"context"

	"github.com/atomic-blend/backend/shared/repositories/user"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v83"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Interface interface {
	GetOrCreateCustomer(ctx *gin.Context, userID primitive.ObjectID) *stripe.Customer
}

type Service struct {
	userService  *user.Repository
	stripeClient *stripe.Client
}

func NewStripeService(userRepo *user.Repository, stripeKey *string) Interface {
	sc := stripe.NewClient(*stripeKey)
	return &Service{
		userService:  userRepo,
		stripeClient: sc,
	}
}

func (s *Service) GetOrCreateCustomer(ctx *gin.Context, userID primitive.ObjectID) *stripe.Customer {
	userEntity, err := s.userService.FindByID(ctx, userID)
	if err != nil {
		log.Error().Str("user_id", userID.Hex()).Msg("cannot find userEntity")
		return nil
	}

	if userEntity.StripeCustomerId == nil {
		// create stripe customer
		params := &stripe.CustomerCreateParams{
			Name:  stripe.String(*userEntity.FirstName + " " + *userEntity.LastName),
			Email: stripe.String(*userEntity.Email),
		}
		result, err := s.stripeClient.V1Customers.Create(context.TODO(), params)
		if err != nil {
			log.Error().Err(err).Msg("error during creation of the stripe customer")
			return nil
		}

		userEntity.StripeCustomerId = &result.ID

		// save customer id to user
		_, err = s.userService.Update(ctx, userEntity)
		if err != nil {
			log.Error().Err(err).Msg("cannot save stripe customer id to user")
			return nil
		}
		
		return result
	} else {
		// get customer and return it
		result, err := s.stripeClient.V1Customers.Retrieve(context.TODO(), *userEntity.StripeCustomerId, nil)
		if err != nil {
			log.Error().Err(err).Msg("cannot get stripe customer")
			return nil
		}
		return result
	}
}
