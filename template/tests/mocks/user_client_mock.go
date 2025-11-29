package mocks

import (
	"context"

	"connectrpc.com/connect"
	userv1 "github.com/atomic-blend/backend/grpc/gen/user/v1"
	"github.com/stretchr/testify/mock"
)

// MockUserClient provides a mock implementation of UserClient
type MockUserClient struct {
	mock.Mock
}

// GetUserPublicKey gets the public key for a user
func (m *MockUserClient) GetUserPublicKey(ctx context.Context, req *connect.Request[userv1.GetUserPublicKeyRequest]) (*connect.Response[userv1.GetUserPublicKeyResponse], error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*connect.Response[userv1.GetUserPublicKeyResponse]), args.Error(1)
}

// GetUserDevices gets the devices for a user
func (m *MockUserClient) GetUserDevices(ctx context.Context, req *connect.Request[userv1.GetUserDevicesRequest]) (*connect.Response[userv1.GetUserDevicesResponse], error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*connect.Response[userv1.GetUserDevicesResponse]), args.Error(1)
}
