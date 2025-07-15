package subscription

import (
	"auth/repositories"
	"auth/utils/db"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IsUserSubscribed checks if a user has an active subscription
func IsUserSubscribed(ctx *gin.Context, userID primitive.ObjectID) bool {
	userRepo := repositories.NewUserRepository(db.Database)

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
		// compare expiration at ms with current time
		if purchase.PurchaseData.ExpirationAtMs > 0 && purchase.PurchaseData.ExpirationAtMs > time.Now().UnixMilli() {
			return true
		}
	}
	return false
}
