package mail

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

func RouteMessage(message *amqp.Delivery) {
	switch message.RoutingKey {
	case "received":
		// Parse the AMQP payload into our structured format
		var payload ReceivedMailPayload
		err := json.Unmarshal(message.Body, &payload)
		if err != nil {
			log.Error().Err(err).Msg("Error unmarshalling AMQP payload")
			return
		}

		// Validate required fields
		if payload.Content == "" {
			log.Error().Msg("Content is empty in AMQP payload")
			return
		}

		// Call receiveMail with the complete payload
		receiveMail(message, payload)
	case "sent":
		//routeSentMessage()
	}
}
