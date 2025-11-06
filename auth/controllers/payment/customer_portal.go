package payment

import (
	"net/http"

	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/gin-gonic/gin"
)

// CustomerPortal handles the creation of a Stripe customer portal session
func (c *Controller) CustomerPortal(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// get the stripe customer from the database
	stripeCustomer := c.stripeService.GetOrCreateCustomer(ctx, authUser.UserID)
	if stripeCustomer == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot_get_stripe_customer"})
		return
	}

	// create the customer portal session
	portalSession, err := c.stripeService.CreateCustomerPortalSession(ctx, stripeCustomer.ID, nil)
	if err != nil || portalSession == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot_create_customer_portal_session"})
		return
	}

	// return the portal session URL
	ctx.JSON(http.StatusOK, gin.H{"url": portalSession.URL})
}
