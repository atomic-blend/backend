package amqpservice

import (
	"github.com/atomic-blend/backend/mail/services/amqp/interfaces"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

// MockAMQPService provides a mock implementation of AMQP service
type MockAMQPService struct {
	mock.Mock
}

// Ensure MockAMQPService implements the interface
var _ interfaces.AMQPServiceInterface = (*MockAMQPService)(nil)

// PublishMessage publishes a message to the AMQP broker
func (m *MockAMQPService) PublishMessage(exchangeName string, topic string, message map[string]interface{}) {
	m.Called(exchangeName, topic, message)
}

// InitProducerAmqp initializes the AMQP producer
func (m *MockAMQPService) InitProducerAmqp() {
	m.Called()
}

// InitConsumerAmqp initializes the AMQP consumer
func (m *MockAMQPService) InitConsumerAmqp() {
	m.Called()
}

// Messages returns the AMQP messages
func (m *MockAMQPService) Messages() <-chan amqp.Delivery {
	return m.Called().Get(0).(<-chan amqp.Delivery)
}