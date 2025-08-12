package amqpservice

import (
	"github.com/atomic-blend/backend/mail/services/amqp/interfaces"
	"github.com/atomic-blend/backend/mail/services/amqp/utils"
	"github.com/streadway/amqp"
)

// AMQPServiceWrapper wraps the existing AMQP functionality
type AMQPServiceWrapper struct{}

// NewAMQPService creates a new AMQP service wrapper
func NewAMQPService() interfaces.AMQPServiceInterface {
	return &AMQPServiceWrapper{}
}

// PublishMessage publishes a message to the AMQP broker
func (a *AMQPServiceWrapper) PublishMessage(exchangeName string, topic string, message map[string]interface{}) {
	utils.PublishMessage(exchangeName, topic, message)
}

// InitProducerAmqp initializes the AMQP producer
func (a *AMQPServiceWrapper) InitProducerAmqp() {
	utils.InitProducerAmqp()
}

// InitConsumerAmqp initializes the AMQP consumer
func (a *AMQPServiceWrapper) InitConsumerAmqp() {
	utils.InitConsumerAmqp()
}

// Messages returns the AMQP messages
func (a *AMQPServiceWrapper) Messages() <-chan amqp.Delivery {
	return utils.Messages
}