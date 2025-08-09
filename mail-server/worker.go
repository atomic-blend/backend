package main

import (
	"github.com/atomic-blend/backend/mail/utils/amqp"
	"github.com/atomic-blend/backend/mail/workers"
)

func processMessages() {
	for m := range amqp.Messages {
		workers.RouteMessage(&m)
	}
}
