package amqpworker

import (
	"github.com/atomic-blend/backend/mail-server/amqp-worker/mail"
	"github.com/streadway/amqp"
)

// RouteMessage routes a message to the appropriate worker
func RouteMessage(message *amqp.Delivery) {
	switch message.Exchange {
	case "mail":
		mail.RouteMessage(message)
	}
}