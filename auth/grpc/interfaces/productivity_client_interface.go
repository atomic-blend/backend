package interfaces

import (
	"context"

	"connectrpc.com/connect"
	"github.com/atomic-blend/backend/grpc/gen/productivity/v1"
)

// ProductivityClientInterface defines the interface for the productivity gRPC client
type ProductivityClientInterface interface {
	DeleteUserData(ctx context.Context, req *connect.Request[productivityv1.DeleteUserDataRequest]) (*connect.Response[productivityv1.DeleteUserDataResponse], error)
}
