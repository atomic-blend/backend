package webhooks

import (
	"auth/models"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EventsToHandle is an array of event types to handle
var EventsToHandle = []string{
	"INITIAL_PURCHASE",
	"RENEWAL",
}

// HandleRevenueCatWebhook handles incoming webhooks from RevenueCat
func (c *Controller) HandleRevenueCatWebhook(ctx *gin.Context) {
	// Parse the request body
	var webhookData models.RevenueCatPayload

	if err := ctx.ShouldBindJSON(&webhookData); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if slices.Contains(EventsToHandle, webhookData.Event.Type) {
		log.Info().Msgf("Handling RevenueCat event: %s", webhookData.Event.Type)
		purchaseEntry := models.NewRevenueCatPurchase(webhookData.Event)
		userID, err := primitive.ObjectIDFromHex(webhookData.Event.AppUserID)
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
		log.Info().Msgf("Ignoring RevenueCat event: %s", webhookData.Event.Type)
		ctx.JSON(200, gin.H{"status": "ignored"})
		return
	}

	ctx.JSON(200, gin.H{"status": "success"})
}
