// Package amqpinterfaces contains the interfaces for the AMQP service
package amqpinterfaces

// AMQPServiceInterface defines the interface for AMQP operations
type AMQPServiceInterface interface {
	PublishMessage(exchangeName string, topic string, message map[string]interface{})
}
