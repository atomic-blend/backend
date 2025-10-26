package payment

import (
	"fmt"
	"net/http"
	"os"

	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// CheckoutRequest represents the expected request body for the Checkout endpoint
type CheckoutRequest struct {
	// Add any fields if needed in the future
	SuccessURL *string `json:"success_url"`
	CancelURL  *string `json:"cancel_url"`
}

// Checkout handles user checkout requests
// @Summary Checkout for a subscription
// @Description Checkout the authenticated user for a subscription
// @Tags Payment
// @Produce json
// @Success 200 {object} map[string]interface{} "Subscription created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /payment/subscribe [post]
func (c *Controller) Checkout(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Parse request body
	var checkoutReq CheckoutRequest
	if err := ctx.ShouldBindJSON(&checkoutReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request_body"})
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

	isExisting := stripeCustomer.Subscriptions != nil && len(stripeCustomer.Subscriptions.Data) > 0

	// if user already have the subscription, return error
	if isExisting {
		log.Debug().Msg("User already has a subscription, cannot proceed to checkout")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "subscription_already_exists"})
		return
	}

	trialDays := int64(0)
	trialDaysEnv := os.Getenv("STRIPE_CLOUD_TRIAL_DAYS")
	if trialDaysEnv != "" {
		fmt.Sscanf(trialDaysEnv, "%d", &trialDays)
	}

	checkoutSession, err := c.stripeService.CreateCheckoutSession(ctx, stripeCustomer.ID, trialDays, checkoutReq.SuccessURL, checkoutReq.CancelURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot_create_checkout_session"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"session": checkoutSession.URL})
}
