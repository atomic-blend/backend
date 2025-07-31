package workers

import "github.com/atomic-blend/backend/mail/workers/mail"

func RouteMessage(exchange string, routingKey string, body []byte) {
	switch exchange {
	case "mail":
		mail.RouteMessage(routingKey, body)
	}
}