package interfaces

import (
	"context"

	"connectrpc.com/connect"
	"github.com/atomic-blend/backend/grpc/gen/productivity"
)

// ProductivityClientInterface defines the interface for the productivity gRPC client
type ProductivityClientInterface interface {
	DeleteUserData(ctx context.Context, req *connect.Request[productivity.DeleteUserDataRequest]) (*connect.Response[productivity.DeleteUserDataResponse], error)
}
