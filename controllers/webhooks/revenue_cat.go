package webhooks

import (
	"atomic_blend_api/models"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// array of event types to handle
var EVENTS_TO_HANDLE = []string{
	"INITIAL_PURCHASE",
	"RENEWAL",
}

func (c *WebhooksController) HandleRevenueCatWebhook(ctx *gin.Context) {
	// Parse the request body
	var webhookData models.RevenueCatPurchaseData 

	if err := ctx.ShouldBindJSON(&webhookData); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if (slices.Contains(EVENTS_TO_HANDLE, webhookData.Type)) {
		log.Info().Msgf("Handling RevenueCat event: %s", webhookData.Type)
		purchaseEntry := models.NewRevenueCatPurchase(webhookData)
		userID, err := primitive.ObjectIDFromHex(webhookData.AppUserID)
		if err != nil {
			log.Error().Err(err).Msg("Invalid user ID")
			ctx.JSON(400, gin.H{"error": "Invalid user ID"})
			return
		}
		if err := c.userRepo.AddPurchase(ctx, userID, &purchaseEntry); err != nil {
			log.Error().Err(err).Msg("Failed to create purchase entry")
			ctx.JSON(500, gin.H{"error": "Failed to create purchase entry"})
			return
		}
	} else {
		log.Info().Msgf("Ignoring RevenueCat event: %s", webhookData.Type)
		ctx.JSON(200, gin.H{"status": "ignored"})
		return
	}
	
	ctx.JSON(200, gin.H{"status": "success"})
}
