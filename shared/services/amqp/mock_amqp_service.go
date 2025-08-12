package amqpservice

import (
	amqpinterfaces "github.com/atomic-blend/backend/shared/services/amqp/interfaces"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

// MockAMQPService provides a mock implementation of AMQP service
type MockAMQPService struct {
	mock.Mock
}

// Ensure MockAMQPService implements the interface
var _ amqpinterfaces.AMQPServiceInterface = (*MockAMQPService)(nil)

// PublishMessage publishes a message to the AMQP broker
func (m *MockAMQPService) PublishMessage(exchangeName string, topic string, message map[string]interface{}, headers *amqp.Table) {
	m.Called(exchangeName, topic, message, headers)
}

// InitProducerAMQP initializes the AMQP producer
func (m *MockAMQPService) InitProducerAMQP() {
	m.Called()
}

// InitConsumerAMQP initializes the AMQP consumer
func (m *MockAMQPService) InitConsumerAMQP() {
	m.Called()
}

// Messages returns the AMQP messages
func (m *MockAMQPService) Messages() <-chan amqp.Delivery {
	return m.Called().Get(0).(<-chan amqp.Delivery)
}