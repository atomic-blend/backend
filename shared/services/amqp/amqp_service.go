// Package amqpservice contains the AMQP service
package amqpservice

import (
	amqpinterfaces "github.com/atomic-blend/backend/shared/services/amqp/interfaces"
	amqputils "github.com/atomic-blend/backend/shared/services/amqp/utils"
	"github.com/streadway/amqp"
)

// AMQPServiceWrapper wraps the existing AMQP functionality
type AMQPServiceWrapper struct{}

// NewAMQPService creates a new AMQP service wrapper
func NewAMQPService() amqpinterfaces.AMQPServiceInterface {
	return &AMQPServiceWrapper{}
}

// PublishMessage publishes a message to the AMQP broker
func (a *AMQPServiceWrapper) PublishMessage(exchangeName string, topic string, message map[string]interface{}, headers *amqp.Table) {
	amqputils.PublishMessage(exchangeName, topic, message, headers)
}

// InitProducerAMQP initializes the AMQP producer
func (a *AMQPServiceWrapper) InitProducerAMQP() {
	amqputils.InitProducerAMQP()
}

// InitConsumerAMQP initializes the AMQP consumer
func (a *AMQPServiceWrapper) InitConsumerAMQP() {
	amqputils.InitConsumerAMQP()
}

// Messages returns the AMQP messages
func (a *AMQPServiceWrapper) Messages() <-chan amqp.Delivery {
	return amqputils.Messages
}
