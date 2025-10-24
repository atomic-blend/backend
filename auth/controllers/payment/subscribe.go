package payment

import (
	"net/http"
	"os"

	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (c *Controller) Subscribe(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Get or create the stripe Customer
	log.Debug().Msg("Fetching or creating Stripe customer")
	stripeCustomer := c.stripeService.GetOrCreateCustomer(ctx, authUser.UserID)
	if stripeCustomer == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot_get_stripe_customer"})
		return
	}

	log.Debug().Msgf("Stripe Customer ID: %s", stripeCustomer.ID)

	priceID := os.Getenv("STRIPE_CLOUD_SUBSCRIPTION_PRICE_ID")

	// if user already have the subscription, return error
	log.Debug().Msgf("subscriptions count: %d", len(stripeCustomer.Subscriptions.Data))
	if stripeCustomer.Subscriptions != nil && len(stripeCustomer.Subscriptions.Data) > 0 {
		log.Debug().Msg("User already has a subscription, fetching existing subscription")
		subscription := c.stripeService.GetSubscription(ctx, stripeCustomer.ID, priceID)
		if subscription != nil && subscription.PendingSetupIntent != nil {
			ctx.JSON(http.StatusOK, gin.H{"subscription": gin.H{
				"secret": subscription.PendingSetupIntent.ClientSecret,
				"intent": subscription.PendingSetupIntent.ID,
			}})
		} else {
			log.Error().Msg("Subscription already exists without pending setup intent")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "subscription_already_exists"})
		}
		return
	}

	// create the subscription
	if priceID == "" {
		log.Error().Msg("STRIPE_CLOUD_SUBSCRIPTION_PRICE_ID is not set")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "subscription_price_not_configured"})
		return
	}

	log.Debug().Msgf("Creating subscription for Price ID: %s", priceID)
	subscription := c.stripeService.CreateSubscription(ctx, stripeCustomer.ID, priceID)
	if subscription == nil {
		log.Error().Msg("Failed to create subscription")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot_create_subscription"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"subscription": gin.H{
		"secret": subscription.PendingSetupIntent.ClientSecret,
		"intent": subscription.PendingSetupIntent.ID,
	}})
}
