package user

import (
	"context"

	"connectrpc.com/connect"
	"github.com/atomic-blend/backend/grpc/gen/productivity"
)

func (s *grpcServer) DeleteUserData(ctx context.Context, req *connect.Request[productivity.DeleteUserDataRequest]) (*connect.Response[productivity.DeleteUserDataResponse], error) {
	return connect.NewResponse(&productivity.DeleteUserDataResponse{
		Success: false,
	}), nil
}
