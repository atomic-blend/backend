package server

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/atomic-blend/backend/shared/models"
	userv1 "github.com/atomic-blend/backend/grpc/gen/user/v1"
)

// GetUserPublicKey is the gRPC method to retrieve user public key
func (userGrpcServer *UserGrpcServer) GetUserPublicKey(ctx context.Context, req *connect.Request[userv1.GetUserPublicKeyRequest]) (*connect.Response[userv1.GetUserPublicKeyResponse], error) {
	user := &models.UserEntity{}
	
	if req.Msg.Id != "" {
		// Call the repository method to get user public key by ID
		user, _ = userGrpcServer.userRepo.GetByID(ctx, req.Msg.Id)
	} else if req.Msg.Email != "" {
		// Call the repository method to get user public key by email
		user, _ = userGrpcServer.userRepo.GetByEmail(ctx, req.Msg.Email)
	} else {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("either id or email must be provided"))
	}
	
	if user == nil || user.KeySet == nil || user.KeySet.PublicKey == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found or public key not set"))
	}

	resp := &userv1.GetUserPublicKeyResponse{
		UserId:    user.ID.Hex(),
		PublicKey: *user.KeySet.PublicKey,
	}

	return connect.NewResponse(resp), nil
}
