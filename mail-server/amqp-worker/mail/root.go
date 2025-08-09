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
		// Parse the AMQP payload into our structured format
		var payload models.RawMail
		err := json.Unmarshal(message.Body, &payload)
		if err != nil {
			log.Error().Err(err).Msg("Error unmarshalling AMQP payload")
			return
		}

		// Call sendMail with the complete payload
		processSendMailMessage(message, payload)
	}
}
