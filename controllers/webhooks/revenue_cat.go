package webhooks

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// array of event types to handle
var EVENTS_TO_HANDLE = []string{
	"INITIAL_PURCHASE",
	"RENEWAL",
}

func (c *WebhooksController) HandleRevenueCatWebhook(ctx *gin.Context) {
	// Parse the request body
	var webhookData struct {
		Event struct {
			Type string `json:"type"`
		} `json:"event"`
	}

	if err := ctx.ShouldBindJSON(&webhookData); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if (slices.Contains(EVENTS_TO_HANDLE, webhookData.Event.Type)) {
		log.Info().Msgf("Handling RevenueCat event: %s", webhookData.Event.Type)
	} else {
		log.Info().Msgf("Ignoring RevenueCat event: %s", webhookData.Event.Type)
		ctx.JSON(200, gin.H{"status": "ignored"})
		return
	}
	
	ctx.JSON(200, gin.H{"status": "success"})
}
