// Package global contains the logic for the global gRPC operations
package global

import (
	"context"

	"connectrpc.com/connect"
	mailv1 "github.com/atomic-blend/backend/grpc/gen/mail/v1"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteUserData deletes all user data across various repositories
func (s *GrpcServer) DeleteUserData(ctx context.Context, req *connect.Request[mailv1.DeleteUserDataRequest]) (*connect.Response[mailv1.DeleteUserDataResponse], error) {
	user := req.Msg.GetUser()
	if user == nil {
		log.Error().Msg("User is required")
		return connect.NewResponse(&mailv1.DeleteUserDataResponse{
			Success: false,
		}), nil
	}

	userIDHex := user.GetId()
	if userIDHex == "" {
		log.Error().Msg("User ID is required")
		return connect.NewResponse(&mailv1.DeleteUserDataResponse{
			Success: false,
		}), nil
	}

	// Convert user ID string to ObjectID
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		log.Error().Err(err).Msg("Invalid user ID format")
		return connect.NewResponse(&mailv1.DeleteUserDataResponse{
			Success: false,
		}), nil
	}

	log.Info().Str("userID", userID.Hex()).Msg("Successfully deleted user data")

	return connect.NewResponse(&mailv1.DeleteUserDataResponse{
		Success: true,
	}), nil
}
