package workers

import (
	"github.com/atomic-blend/backend/mail/workers/mail"
	"github.com/streadway/amqp"
)

func RouteMessage(message *amqp.Delivery) {
	switch message.Exchange {
	case "mail":
		mail.RouteMessage(message)
	}
}