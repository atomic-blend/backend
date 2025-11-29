package global

import (
	"context"

	"connectrpc.com/connect"
	calendarv1 "github.com/atomic-blend/backend/grpc/gen/calendar/v1"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteUserData handles the deletion of all user data
func (s *GrpcServer) DeleteUserData(ctx context.Context, req *connect.Request[calendarv1.DeleteUserDataRequest]) (*connect.Response[calendarv1.DeleteUserDataResponse], error) {
	user := req.Msg.GetUser()
	if user == nil {
		log.Error().Msg("User is required")
		return connect.NewResponse(&calendarv1.DeleteUserDataResponse{
			Success: false,
		}), nil
	}

	userIDHex := user.GetId()
	if userIDHex == "" {
		log.Error().Msg("User ID is required")
		return connect.NewResponse(&calendarv1.DeleteUserDataResponse{
			Success: false,
		}), nil
	}

	// Convert user ID string to ObjectID
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		log.Error().Err(err).Msg("Invalid user ID format")
		return connect.NewResponse(&calendarv1.DeleteUserDataResponse{
			Success: false,
		}), nil
	}

	//TODO: Implement deletion logic across various repositories

	log.Info().Str("userID", userID.Hex()).Msg("Successfully deleted user data")

	return connect.NewResponse(&calendarv1.DeleteUserDataResponse{
		Success: true,
	}), nil
}
