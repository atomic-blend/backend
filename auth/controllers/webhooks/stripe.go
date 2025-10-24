package webhooks

import (
	"io"
	"os"
	"time"

	stripeutils "github.com/atomic-blend/backend/auth/utils/stripe"
	"github.com/atomic-blend/backend/shared/repositories/user"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/webhook"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HandleStripeWebhook processes incoming Stripe webhook events
func (c *Controller) HandleStripeWebhook(ctx *gin.Context) {
	stripeSecret := os.Getenv("STRIPE_WEBHOOK_TOKEN")
	if stripeSecret == "" {
		log.Error().Msg("Stripe webhook token not set in environment variables")
		ctx.JSON(500, gin.H{"error": "Server configuration error"})
		return
	}

	log.Debug().Msg("Verifying payload signature")
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading request body")
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	event, err := webhook.ConstructEvent(bodyBytes, ctx.Request.Header.Get("Stripe-Signature"), stripeSecret)
	if err != nil {
		log.Error().Err(err).Msg("Error verifying Stripe webhook signature")
		ctx.JSON(400, gin.H{"error": "Invalid signature"})
		return
	}

	log.Info().Str("event_type", string(event.Type)).Msg("Processing Stripe webhook event")

	log.Debug().Interface("event", event.Data.Object).Msg("Handling specific event type")

	switch event.Type {
	case "invoice.paid":
		handleInvoicePaid(ctx, c.stripeService, c.userRepo, &event)
	case "invoice.payment_failed":
		handleInvoicePaymentFailed(ctx, c.stripeService, c.userRepo, &event)
	case "customer.subscription.deleted":
		handleSubscriptionDeleted(ctx, c.stripeService, c.userRepo, &event)
	default:
		log.Warn().Str("event_type", string(event.Type)).Msg("Unhandled Stripe webhook event type")
		ctx.JSON(200, gin.H{"result": "unhandled_event_type"})
		return
	}

	log.Debug().Msg("Stripe webhook event processed successfully")
}

func getCustomer(ctx *gin.Context, stripeService stripeutils.Interface, event *stripe.Event) *stripe.Customer {
	log.Debug().Msg("Fetching customer for Stripe event")
	customerID := event.Data.Object["customer"].(string)
	log.Debug().Str("customer_id", customerID).Msg("Retrieving customer from Stripe")
	customer, err := stripeService.GetCustomer(ctx, customerID, &stripe.CustomerRetrieveParams{})
	if err != nil {
		log.Error().Err(err).Str("customer_id", customerID).Msg("Error fetching customer for invoice.paid event")
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return nil
	}
	return customer
}

func handleInvoicePaid(ctx *gin.Context, stripeService stripeutils.Interface, userRepo user.Interface, event *stripe.Event) {
	customer := getCustomer(ctx, stripeService, event)
	if customer == nil {
		return
	}

	log.Debug().Str("customer_id", customer.ID).Msg("Fetched customer for invoice.paid event")

	log.Info().Msg("finding user for invoice.paid event")
	appUserID := customer.Metadata["app_user_id"]
	if appUserID == "" {
		log.Error().Str("customer_id", customer.ID).Msg("No app_user_id metadata found on customer")
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	appUserObjectID, err := primitive.ObjectIDFromHex(appUserID)
	if err != nil {
		log.Error().Err(err).Str("app_user_id", appUserID).Msg("Invalid app_user_id metadata on customer")
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	user, err := userRepo.FindByID(ctx, appUserObjectID)
	if err != nil {
		log.Error().Err(err).Str("app_user_id", appUserID).Msg("Error fetching user for invoice.paid event")
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	log.Info().Msg("updating user subscription status for invoice.paid event")
	user.SubscriptionStatus = stripe.String("active")
	subscriptionID := event.Data.Object["parent"].(map[string]interface{})["subscription_details"].(map[string]interface{})["subscription"].(string)
	user.StripeSubscriptionID = &subscriptionID
	user.FailedAt = nil
	user.CancelledAt = nil

	_, err = userRepo.Update(ctx, user)
	if err != nil {
		log.Error().Err(err).Str("app_user_id", appUserID).Msg("Error updating user subscription status for invoice.paid event")
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	log.Info().Msg("Invoice payment succeeded")
	ctx.JSON(200, gin.H{"status": "success"})
}

func handleInvoicePaymentFailed(ctx *gin.Context, stripeService stripeutils.Interface, userRepo user.Interface, event *stripe.Event) {
	customer := getCustomer(ctx, stripeService, event)
	if customer == nil {
		return
	}

	log.Info().Msg("finding user for invoice.payment_failed event")
	appUserID := customer.Metadata["app_user_id"]
	if appUserID == "" {
		log.Error().Str("customer_id", customer.ID).Msg("No app_user_id metadata found on customer")
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	appUserObjectID, err := primitive.ObjectIDFromHex(appUserID)
	if err != nil {
		log.Error().Err(err).Str("app_user_id", appUserID).Msg("Invalid app_user_id metadata on customer")
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	user, err := userRepo.FindByID(ctx, appUserObjectID)
	if err != nil {
		log.Error().Err(err).Str("app_user_id", appUserID).Msg("Error fetching user for invoice.paid event")
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	log.Info().Msg("updating user subscription status for invoice.payment_failed event")
	user.SubscriptionStatus = stripe.String("past_due")
	now := primitive.NewDateTimeFromTime(time.Now())
	user.FailedAt = &now
	user.CancelledAt = nil

	_, err = userRepo.Update(ctx, user)
	if err != nil {
		log.Error().Err(err).Str("app_user_id", appUserID).Msg("Error updating user subscription status for invoice.payment_failed event")
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	log.Info().Msg("Invoice payment failed")
	ctx.JSON(200, gin.H{"status": "success"})
}

func handleSubscriptionDeleted(ctx *gin.Context, stripeService stripeutils.Interface, userRepo user.Interface, event *stripe.Event) {
	customer := getCustomer(ctx, stripeService, event)
	if customer == nil {
		return
	}

	log.Info().Msg("finding user for customer.subscription.deleted event")
	appUserID := customer.Metadata["app_user_id"]
	if appUserID == "" {
		log.Error().Str("customer_id", customer.ID).Msg("No app_user_id metadata found on customer")
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	appUserObjectID, err := primitive.ObjectIDFromHex(appUserID)
	if err != nil {
		log.Error().Err(err).Str("app_user_id", appUserID).Msg("Invalid app_user_id metadata on customer")
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	user, err := userRepo.FindByID(ctx, appUserObjectID)
	if err != nil {
		log.Error().Err(err).Str("app_user_id", appUserID).Msg("Error fetching user for invoice.paid event")
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	log.Info().Msg("updating user subscription status for customer.subscription.deleted event")
	user.SubscriptionStatus = stripe.String("cancelled")
	now := primitive.NewDateTimeFromTime(time.Now())
	user.CancelledAt = &now

	_, err = userRepo.Update(ctx, user)
	if err != nil {
		log.Error().Err(err).Str("app_user_id", appUserID).Msg("Error updating user subscription status for invoice.payment_failed event")
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	log.Info().Msg("Customer subscription deleted")
	ctx.JSON(200, gin.H{"status": "success"})
}
