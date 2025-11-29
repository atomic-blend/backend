package main

import (
	amqpinterfaces "github.com/atomic-blend/backend/shared/services/amqp/interfaces"
	"github.com/atomic-blend/backend/mail/workers"
)

func processMessages(amqpService amqpinterfaces.AMQPServiceInterface) {
	for m := range amqpService.Messages() {
		workers.RouteMessage(&m)
	}
}
