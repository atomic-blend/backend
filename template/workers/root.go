package workers

import (
	"github.com/atomic-blend/backend/template/workers/template"
	"github.com/streadway/amqp"
)

// RouteMessage routes a message to the appropriate worker
func RouteMessage(message *amqp.Delivery) {
	switch message.Exchange {
	case "template":
		template.RouteMessage(message)
	}
}
