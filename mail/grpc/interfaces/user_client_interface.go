package interfaces

import (
	"context"

	"connectrpc.com/connect"
	userv1 "github.com/atomic-blend/backend/grpc/gen/user/v1"
)

// UserClientInterface defines the methods for user-related gRPC operations
type UserClientInterface interface {
	GetUserDevices(context.Context, *connect.Request[userv1.GetUserDevicesRequest]) (*connect.Response[userv1.GetUserDevicesResponse], error)
	GetUserPublicKey(context.Context, *connect.Request[userv1.GetUserPublicKeyRequest]) (*connect.Response[userv1.GetUserPublicKeyResponse], error)
}