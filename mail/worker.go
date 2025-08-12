package main

import (
	"github.com/atomic-blend/backend/mail/services/amqp/interfaces"
	"github.com/atomic-blend/backend/mail/workers"
)

func processMessages(amqpService interfaces.AMQPServiceInterface) {
	for m := range amqpService.Messages() {
		workers.RouteMessage(&m)
	}
}
