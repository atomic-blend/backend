// Package amqpinterfaces contains the interfaces for the AMQP service
package amqpinterfaces

import "github.com/streadway/amqp"

// AMQPServiceInterface defines the interface for AMQP operations
type AMQPServiceInterface interface {
	Messages() <-chan amqp.Delivery
	InitProducerAMQP()
	InitConsumerAMQP()
	PublishMessage(exchangeName string, topic string, message map[string]interface{})
}
