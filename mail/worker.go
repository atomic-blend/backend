package main

import (
	"github.com/atomic-blend/backend/mail/utils/amqp"
	"github.com/rs/zerolog/log"
)

func processMessages() {
	for m := range amqp.Messages {
		exchange := m.Exchange
		routingKey := m.RoutingKey
		body := m.Body

		// do this in a separate mail module
		switch exchange {
		case "mail":
			switch routingKey {
			case "received":
				//TODO:
				log.Debug().Str("body", string(body)).Msg("Received message")
			case "sent":
				//TODO:
				log.Debug().Str("body", string(body)).Msg("Sent message")
			}
		}
	}
}
