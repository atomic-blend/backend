// Package amqpservice contains the AMQP service
package amqpservice

import (
	amqpinterfaces "github.com/atomic-blend/backend/shared/services/amqp/interfaces"
	amqputils "github.com/atomic-blend/backend/shared/services/amqp/utils"
	"github.com/streadway/amqp"
)

// AMQPServiceWrapper wraps the existing AMQP functionality
type AMQPServiceWrapper struct {
	workerName string
}

// NewAMQPService creates a new AMQP service wrapper
func NewAMQPService(workerName string) amqpinterfaces.AMQPServiceInterface {
	return &AMQPServiceWrapper{
		workerName: workerName,
	}
}

// PublishMessage publishes a message to the AMQP broker
func (a *AMQPServiceWrapper) PublishMessage(exchangeName string, topic string, message map[string]interface{}, headers *amqp.Table) {
	amqputils.PublishMessage(exchangeName, topic, message, headers)
}

// InitProducerAMQP initializes the AMQP producer
func (a *AMQPServiceWrapper) InitProducerAMQP() {
	amqputils.InitProducerAMQP(a.workerName)
}

// InitConsumerAMQP initializes the AMQP consumer
func (a *AMQPServiceWrapper) InitConsumerAMQP() {
	amqputils.InitConsumerAMQP(a.workerName)
}

// Messages returns the AMQP messages
func (a *AMQPServiceWrapper) Messages() <-chan amqp.Delivery {
	return amqputils.Messages
}
