package smtpserver

import (
	"errors"
	"io"
)

// MockReader is a mock io.Reader that can simulate errors
type MockReader struct {
	shouldError bool
	content     string
	readCount   int
}

func (m *MockReader) Read(p []byte) (n int, err error) {
	if m.shouldError {
		return 0, errors.New("read error")
	}

	if m.readCount >= len(m.content) {
		return 0, io.EOF
	}

	n = copy(p, m.content[m.readCount:])
	m.readCount += n
	return n, nil
}

// MockAMQP is a mock for the AMQP functionality
type MockAMQP struct {
	publishedMessages []map[string]interface{}
}

func (m *MockAMQP) PublishMessage(exchangeName string, topic string, message map[string]interface{}) {
	m.publishedMessages = append(m.publishedMessages, message)
}

// GetPublishedMessages returns all published messages for testing
func (m *MockAMQP) GetPublishedMessages() []map[string]interface{} {
	return m.publishedMessages
}

// ClearPublishedMessages clears the published messages
func (m *MockAMQP) ClearPublishedMessages() {
	m.publishedMessages = nil
}
