package mail

import (
	"encoding/json"

	"github.com/atomic-blend/backend/mail/models"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

// RouteMessage routes a message to the appropriate worker
func RouteMessage(message *amqp.Delivery) {
	switch message.RoutingKey {
	case "sent":
		log.Info().Msg("📤 Processing sent message")
		// Parse the AMQP payload into our structured format
		var payload models.RawMail
		err := json.Unmarshal(message.Body, &payload)
		if err != nil {
			log.Error().Err(err).Msg("Error unmarshalling AMQP payload")
			return
		}

		// Call sendMail with the complete payload
		processSendMailMessage(message, payload)
	case "send_retry":
		log.Info().Msg("🔄 Processing send retry message")
		// Parse the AMQP payload into our structured format
		var payload models.RawMail
		err := json.Unmarshal(message.Body, &payload)
		if err != nil {
			log.Error().Err(err).Msg("Error unmarshalling AMQP payload for retry")
			return
		}

		// Call sendMail with the complete payload
		processSendMailMessage(message, payload)
	default:
		log.Warn().Str("routing_key", message.RoutingKey).Msg("⚠️ Unknown routing key, message not processed")
	}
}
