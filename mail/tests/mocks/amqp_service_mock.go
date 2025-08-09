package mocks

import (
	"github.com/atomic-blend/backend/mail/services/interfaces"
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
