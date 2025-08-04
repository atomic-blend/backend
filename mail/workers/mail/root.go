package mail

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

func RouteMessage(routingKey string, body []byte) {
	switch routingKey {
	case "received":
		// Parse the AMQP payload into our structured format
		var payload ReceivedMailPayload
		err := json.Unmarshal(body, &payload)
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
		receiveMail(payload)
	case "sent":
		//routeSentMessage()
	}
}
