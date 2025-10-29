package subscription

import (
	"time"

	"github.com/atomic-blend/backend/shared/repositories/user"
	"github.com/atomic-blend/backend/shared/utils/db"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IsUserSubscribed checks if a user has an active subscription
func IsUserSubscribed(ctx *gin.Context, userID primitive.ObjectID) bool {
	userRepo := user.NewUserRepository(db.Database)

	user, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		return false
	}

	if user == nil {
		return false
	}

	// check in the user's purchase if the user has an active subscription
	purchases := user.Purchases
	for _, purchase := range purchases {
		// Handle RevenueCat purchases
		if rcData, ok := purchase.GetRevenueCatData(); ok {
			if rcData.ExpirationAtMs > 0 && rcData.ExpirationAtMs > time.Now().UnixMilli() {
				return true
			}
		}

		// Handle Stripe purchases
		if stripeData, ok := purchase.GetStripeData(); ok {
			if stripeData.Status == "active" {
				return true
			}
		}
	}
	return false
}
