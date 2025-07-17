package interfaces

import (
	"context"

	"connectrpc.com/connect"
	"github.com/atomic-blend/backend/grpc/gen/auth"
)

// UserClientInterface defines the methods for user-related gRPC operations
type UserClientInterface interface {
	GetUserDevices(context.Context, *connect.Request[auth.GetUserDevicesRequest]) (*connect.Response[auth.GetUserDevicesResponse], error)
}