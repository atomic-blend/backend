package mocks

import (
	"context"

	"connectrpc.com/connect"
	"github.com/atomic-blend/backend/grpc/gen/productivity"
	"github.com/stretchr/testify/mock"
)

// MockProductivityClient is a mock for the ProductivityClientInterface
type MockProductivityClient struct {
	mock.Mock
}

// DeleteUserData mocks the DeleteUserData method
func (m *MockProductivityClient) DeleteUserData(ctx context.Context, req *connect.Request[productivity.DeleteUserDataRequest]) (*connect.Response[productivity.DeleteUserDataResponse], error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*connect.Response[productivity.DeleteUserDataResponse]), args.Error(1)
}
