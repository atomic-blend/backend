package payment

import (
	"fmt"
	"net/http"
	"os"

	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Subscribe handles user subscription requests
// @Summary Subscribe to a plan
// @Description Subscribe the authenticated user to a plan
// @Tags Payment
// @Produce json
// @Success 200 {object} map[string]interface{} "Subscription created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /payment/subscribe [post]
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
	trialDays := os.Getenv("STRIPE_CLOUD_TRIAL_DAYS")
	var trialDaysInt64 int64 = 0
	if trialDays != "" {
		fmt.Sscanf(trialDays, "%d", &trialDaysInt64)
	}
	subscription := c.stripeService.CreateSubscription(ctx, stripeCustomer.ID, priceID, trialDaysInt64)
	if subscription == nil {
		log.Error().Msg("Failed to create subscription")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot_create_subscription"})
		return
	}

	//TODO: store inside the user the subcription ID and status to trialing

	fmt.Printf("%+v\n", subscription)

	log.Debug().Msgf("Subscription created with ID: %s", subscription.ID)
	log.Debug().Msgf("Pending Setup Intent: %v", subscription.PendingSetupIntent)

	ctx.JSON(http.StatusOK, gin.H{"pending_setup_intent": gin.H{
		"secret":    subscription.PendingSetupIntent.ClientSecret,
		"intent_id": subscription.PendingSetupIntent.ID,
	}})
}
