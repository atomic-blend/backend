// Package mail contains the logic for processing mail messages
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
		log.Info().Msgf("📤 Processing sent message: %s", string(message.Body))

		// The message body contains a wrapper with send_mail_id and content fields
		var messageWrapper map[string]interface{}
		err := json.Unmarshal(message.Body, &messageWrapper)
		if err != nil {
			log.Error().Err(err).Msg("Error unmarshalling AMQP payload wrapper")
			return
		}

		// Extract the content field which contains the actual RawMail data
		contentData, ok := messageWrapper["content"]
		if !ok {
			log.Error().Msg("Message wrapper missing 'content' field")
			return
		}

		// Convert the content to JSON bytes for proper unmarshalling
		contentBytes, err := json.Marshal(contentData)
		if err != nil {
			log.Error().Err(err).Msg("Error marshalling content data")
			return
		}

		var payload models.RawMail
		err = json.Unmarshal(contentBytes, &payload)
		if err != nil {
			log.Error().Err(err).Msg("Error unmarshalling RawMail content")
			return
		}

		log.Info().Msgf("📤 Processing sent message: %+v", payload)

		// Call sendMail with the complete payload
		processSendMailMessage(message, payload)
	default:
		log.Warn().Str("routing_key", message.RoutingKey).Msg("⚠️ Unknown routing key, message not processed")
	}
}
