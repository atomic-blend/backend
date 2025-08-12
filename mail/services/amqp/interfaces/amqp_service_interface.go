package interfaces

import "github.com/streadway/amqp"

// AMQPServiceInterface defines the interface for AMQP operations
type AMQPServiceInterface interface {
	Messages() <-chan amqp.Delivery
	InitProducerAmqp()
	InitConsumerAmqp()
	PublishMessage(exchangeName string, topic string, message map[string]interface{})
}
