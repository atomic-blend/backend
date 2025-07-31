package main

import (
	"github.com/atomic-blend/backend/mail/utils/amqp"
	"github.com/atomic-blend/backend/mail/workers"
	"github.com/rs/zerolog/log"
)

func processMessages() {
	for m := range amqp.Messages {
		exchange := m.Exchange
		routingKey := m.RoutingKey
		body := m.Body

		log.Debug().Str("exchange", exchange).Str("routingKey", routingKey).Msg("routing message to exchange")
		workers.RouteMessage(exchange, routingKey, body)
	}
}
