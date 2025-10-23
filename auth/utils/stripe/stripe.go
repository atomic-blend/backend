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
	CreateSubscription(ctx *gin.Context, customerID string, priceID string) *stripe.Subscription
}

type Service struct {
	userService  user.Interface
	stripeClient ClientInterface
}

func NewStripeService(userRepo *user.Repository, stripeKey *string) Interface {
	sc := stripe.NewClient(*stripeKey)
	wrapper := &ClientWrapper{client: sc}
	return &Service{
		userService:  userRepo,
		stripeClient: wrapper,
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
		result, err := s.stripeClient.CreateCustomer(context.TODO(), params)
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
		params := &stripe.CustomerRetrieveParams{
			Expand: []*string{stripe.String("subscriptions")},
		}
		result, err := s.stripeClient.GetCustomer(context.TODO(), *userEntity.StripeCustomerId, params)
		if err != nil {
			log.Error().Err(err).Msg("cannot get stripe customer")
			return nil
		}
		return result
	}
}

func (s *Service) CreateSubscription(ctx *gin.Context, customerID string, priceID string) *stripe.Subscription {
	params := &stripe.SubscriptionCreateParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionCreateItemParams{
			&stripe.SubscriptionCreateItemParams{
				Price: stripe.String(priceID),
			},
		},
		PaymentBehavior: stripe.String(string("default_incomplete")),
		PaymentSettings: &stripe.SubscriptionCreatePaymentSettingsParams{
			SaveDefaultPaymentMethod: stripe.String(string("on_subscription")),
		},
		Expand: []*string{stripe.String("pending_setup_intent")},
	}
	result, err := s.stripeClient.CreateSubscription(context.TODO(), params)
	if err != nil {
		log.Error().Err(err).Msg("error during creation of the stripe subscription")
		return nil
	}
	return result
}
