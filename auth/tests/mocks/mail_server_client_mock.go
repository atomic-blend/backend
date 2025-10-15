// Package mocks provides mock implementations for testing
package mocks

import (
	"context"

	"connectrpc.com/connect"
	mailserverv1 "github.com/atomic-blend/backend/grpc/gen/mail-server/v1"
	"github.com/stretchr/testify/mock"
)

// MockMailServerClient is a mock for the MailServerClientInterface
type MockMailServerClient struct {
	mock.Mock
}

// SendMailInternal mocks the SendMailInternal method
func (m *MockMailServerClient) SendMailInternal(ctx context.Context, req *connect.Request[mailserverv1.SendMailInternalRequest]) (*connect.Response[mailserverv1.SendMailInternalResponse], error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*connect.Response[mailserverv1.SendMailInternalResponse]), args.Error(1)
}
