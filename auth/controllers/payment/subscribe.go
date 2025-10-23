package payment

import (
	"net/http"

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
	}

	log.Debug().Msgf("Stripe Customer ID: %s", stripeCustomer.ID)
}
