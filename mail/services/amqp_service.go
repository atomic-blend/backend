package services

import (
	"github.com/atomic-blend/backend/mail/services/interfaces"
	"github.com/atomic-blend/backend/mail/utils/amqp"
)

// AMQPServiceWrapper wraps the existing AMQP functionality
type AMQPServiceWrapper struct{}

// NewAMQPService creates a new AMQP service wrapper
func NewAMQPService() interfaces.AMQPServiceInterface {
	return &AMQPServiceWrapper{}
}

// PublishMessage publishes a message to the AMQP broker
func (a *AMQPServiceWrapper) PublishMessage(exchangeName string, topic string, message map[string]interface{}) {
	amqp.PublishMessage(exchangeName, topic, message)
}
