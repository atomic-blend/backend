package webhooks

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v83/webhook"
)

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

	switch event.Type {
	case "invoice.paid":
		// Handle successful payment
		//TODO: update user subscription status to active + remove any failed at
		log.Info().Msg("Invoice payment succeeded")
	case "invoice.payment_failed":
		// Handle failed payment
		//TODO: update user subscription status to failed + store failed at
		log.Info().Msg("Invoice payment failed")
	case "customer.subscription.deleted":
		// Handle subscription cancellation (by payment or by user)
		//TODO: update user subscription status to cancelled + store cancelled at
		log.Info().Msg("Customer subscription deleted")
	default:
		log.Warn().Str("event_type", string(event.Type)).Msg("Unhandled Stripe webhook event type")
	}

	log.Debug().Msg("Stripe webhook event processed successfully")
	ctx.JSON(200, gin.H{"status": "success"})
}
