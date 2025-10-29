package stripe

import (
	"context"
	"os"
	"time"

	"github.com/atomic-blend/backend/shared/repositories/user"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v83"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Interface defines the methods for interacting with Stripe services
type Interface interface {
	GetOrCreateCustomer(ctx *gin.Context, userID primitive.ObjectID) *stripe.Customer
	GetCustomer(ctx *gin.Context, customerID string, params *stripe.CustomerRetrieveParams) (*stripe.Customer, error)
	CreateSubscription(ctx *gin.Context, customerID string, priceID string, trialDays int64) *stripe.Subscription
	GetSubscription(ctx *gin.Context, customerID string, priceID string) *stripe.Subscription
	CreateInvoice(ctx *gin.Context, customerID string, subscriptionID string) *stripe.Invoice
	CreateInvoiceItem(ctx *gin.Context, customerID string, amount float64, description string) *stripe.InvoiceItem
	FinalizeInvoice(ctx *gin.Context, invoiceID string) *stripe.Invoice
	GetEphemeralKeys(ctx *gin.Context, customerID string) *stripe.EphemeralKey
	CreateCheckoutSession(ctx *gin.Context, customerID string, trialDays int64, successURL *string, cancelURL *string) (*stripe.CheckoutSession, error)
}

// Service implements the Stripe service interface
type Service struct {
	userService  user.Interface
	stripeClient ClientInterface
}

// NewStripeService creates a new instance of the Stripe service
func NewStripeService(userRepo *user.Repository, stripeKey *string) Interface {
	sc := stripe.NewClient(*stripeKey)
	wrapper := &ClientWrapper{client: sc}
	return &Service{
		userService:  userRepo,
		stripeClient: wrapper,
	}
}

// GetOrCreateCustomer retrieves an existing Stripe customer or creates a new one if it doesn't exist
func (s *Service) GetOrCreateCustomer(ctx *gin.Context, userID primitive.ObjectID) *stripe.Customer {
	userEntity, err := s.userService.FindByID(ctx, userID)
	if err != nil {
		log.Error().Str("user_id", userID.Hex()).Msg("cannot find userEntity")
		return nil
	}

	if userEntity.StripeCustomerID == nil {
		// create stripe customer
		params := &stripe.CustomerCreateParams{
			Name:  stripe.String(*userEntity.FirstName + " " + *userEntity.LastName),
			Email: stripe.String(*userEntity.Email),
			Metadata: map[string]string{
				"app_user_id": userEntity.ID.Hex(),
			},
		}
		result, err := s.stripeClient.CreateCustomer(context.TODO(), params)
		if err != nil {
			log.Error().Err(err).Msg("error during creation of the stripe customer")
			return nil
		}

		userEntity.StripeCustomerID = &result.ID

		// save customer id to user
		_, err = s.userService.Update(ctx, userEntity)
		if err != nil {
			log.Error().Err(err).Msg("cannot save stripe customer id to user")
			return nil
		}

		return result
	}
	// get customer and return it
	params := &stripe.CustomerRetrieveParams{
		Expand: []*string{stripe.String("subscriptions")},
	}
	result, err := s.stripeClient.GetCustomer(context.TODO(), *userEntity.StripeCustomerID, params)
	if err != nil {
		log.Error().Err(err).Msg("cannot get stripe customer")
		return nil
	}
	log.Debug().Str("customer_id", result.ID).Interface("metadata", result.Metadata).Msg("Retrieved existing Stripe customer")
	return result
}

// CreateSubscription creates a new Stripe subscription for the given customer and price ID
func (s *Service) CreateSubscription(ctx *gin.Context, customerID string, priceID string, trialDays int64) *stripe.Subscription {
	trialEnd := time.Now().AddDate(0, 0, int(trialDays)).Unix()
	params := &stripe.SubscriptionCreateParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionCreateItemParams{
			{
				Price: stripe.String(priceID),
			},
		},
		TrialEnd:           stripe.Int64(trialEnd),
		BillingCycleAnchor: stripe.Int64(trialEnd + 1000),
		CollectionMethod:   stripe.String("charge_automatically"),
		PaymentBehavior:    stripe.String("default_incomplete"),
		TrialSettings: &stripe.SubscriptionCreateTrialSettingsParams{
			EndBehavior: &stripe.SubscriptionCreateTrialSettingsEndBehaviorParams{
				MissingPaymentMethod: stripe.String("cancel"),
			},
		},
		Expand: []*string{stripe.String("latest_invoice.payment_intent"), stripe.String("pending_setup_intent")},
	}
	result, err := s.stripeClient.CreateSubscription(context.TODO(), params)
	if err != nil {
		log.Error().Err(err).Msg("error during creation of the stripe subscription")
		return nil
	}
	return result
}

// GetSubscription retrieves an existing Stripe subscription for the given customer and price ID
func (s *Service) GetSubscription(ctx *gin.Context, customerID string, priceID string) *stripe.Subscription {
	// get customer and return it
	params := &stripe.CustomerRetrieveParams{
		Expand: []*string{stripe.String("subscriptions.data.pending_setup_intent")},
	}
	result, err := s.stripeClient.GetCustomer(context.TODO(), customerID, params)
	if err != nil {
		log.Error().Err(err).Msg("cannot get stripe customer")
		return nil
	}

	for _, sub := range result.Subscriptions.Data {
		for _, item := range sub.Items.Data {
			if item.Price.ID == priceID {
				return sub
			}
		}
	}
	return nil
}

// CreateInvoice creates a new Stripe invoice for the given customer and subscription ID
func (s *Service) CreateInvoice(ctx *gin.Context, customerID string, subscriptionID string) *stripe.Invoice {
	params := &stripe.InvoiceCreateParams{
		Customer:         stripe.String(customerID),
		Subscription:     stripe.String(subscriptionID),
		CollectionMethod: stripe.String("charge_automatically"),
		AutoAdvance:      stripe.Bool(true),
	}

	result, err := s.stripeClient.CreateInvoice(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("error during creation of the stripe invoice")
		return nil
	}
	return result
}

// CreateInvoiceItem creates a new Stripe invoice item for the given customer, amount, and description
func (s *Service) CreateInvoiceItem(ctx *gin.Context, customerID string, amount float64, description string) *stripe.InvoiceItem {
	params := &stripe.InvoiceItemCreateParams{
		Customer:    stripe.String(customerID),
		Amount:      stripe.Int64(int64(amount * 100)), // amount in cents
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String(description),
	}

	result, err := s.stripeClient.CreateInvoiceItem(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("error during creation of the stripe invoice item")
		return nil
	}
	return result
}

// FinalizeInvoice finalizes a Stripe invoice by its ID
func (s *Service) FinalizeInvoice(ctx *gin.Context, invoiceID string) *stripe.Invoice {
	params := &stripe.InvoiceFinalizeInvoiceParams{}
	result, err := s.stripeClient.FinalizeInvoice(ctx, invoiceID, params)
	if err != nil {
		log.Error().Err(err).Msg("error during finalization of the stripe invoice")
		return nil
	}
	return result
}

// GetEphemeralKeys retrieves ephemeral keys for the given customer ID
func (s *Service) GetEphemeralKeys(ctx *gin.Context, customerID string) *stripe.EphemeralKey {
	params := &stripe.EphemeralKeyCreateParams{
		Customer: stripe.String(customerID),
	}
	result, err := s.stripeClient.GetEphemeralKeys(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("error during retrieval of the stripe ephemeral keys")
		return nil
	}
	return result
}

// GetCustomer retrieves a Stripe customer by ID.
func (s *Service) GetCustomer(ctx *gin.Context, customerID string, params *stripe.CustomerRetrieveParams) (*stripe.Customer, error) {
	return s.stripeClient.GetCustomer(context.TODO(), customerID, params)
}

// CreateCheckoutSession creates a new Stripe checkout session for the given customer ID and trial days
func (s *Service) CreateCheckoutSession(ctx *gin.Context, customerID string, trialDays int64, successURL *string, cancelURL *string) (*stripe.CheckoutSession, error) {
	publicAddress := os.Getenv("PUBLIC_ADDRESS")
	if publicAddress == "" {
		publicAddress = "http://localhost:53631"
	}
	https := os.Getenv("HTTPS") == "true"

	var baseURL string
	if https {
		baseURL = "https://" + publicAddress + "/#"
	} else {
		baseURL = "http://" + publicAddress + "/#"
	}

	cloudPriceID := os.Getenv("STRIPE_CLOUD_SUBSCRIPTION_PRICE_ID")
	if cloudPriceID == "" {
		log.Error().Msg("STRIPE_CLOUD_SUBSCRIPTION_PRICE_ID is not set")
		return nil, nil
	}

	storagePriceID := os.Getenv("STRIPE_STORAGE_PRICE_ID")
	if storagePriceID == "" {
		log.Error().Msg("STRIPE_STORAGE_PRICE_ID is not set")
		return nil, nil
	}

	betaCouponID := os.Getenv("STRIPE_BETA_COUPON_ID")
	if betaCouponID != "" {
		log.Debug().Msgf("Applying beta coupon: %s", betaCouponID)
	}

	discounts := []*stripe.CheckoutSessionCreateDiscountParams{}
	if betaCouponID != "" {
		discounts = append(discounts, &stripe.CheckoutSessionCreateDiscountParams{
			Coupon: stripe.String(betaCouponID),
		})
	}

	if successURL == nil {
		successURL = stripe.String(baseURL + "/paywall?success=true")
	}
	if cancelURL == nil {
		cancelURL = stripe.String(baseURL + "/paywall?canceled=true")
	}

	params := &stripe.CheckoutSessionCreateParams{
		Customer: stripe.String(customerID),
		LineItems: []*stripe.CheckoutSessionCreateLineItemParams{
			{
				Price:    stripe.String(cloudPriceID),
				Quantity: stripe.Int64(1),
			},
			{
				Price: stripe.String(storagePriceID),
			},
		},
		Discounts: discounts,
		SubscriptionData: &stripe.CheckoutSessionCreateSubscriptionDataParams{
			TrialPeriodDays: stripe.Int64(trialDays),
			// BillingCycleAnchor: stripe.Int64(trialEnd + 100),
		},
		Mode:                     stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL:               successURL,
		CancelURL:                cancelURL,
		BillingAddressCollection: stripe.String(string(stripe.CheckoutSessionBillingAddressCollectionAuto)),
	}

	return s.stripeClient.CreateCheckoutSession(context.TODO(), params)
}
