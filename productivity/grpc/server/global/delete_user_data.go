package user

import (
	"context"

	"connectrpc.com/connect"
	"github.com/atomic-blend/backend/grpc/gen/productivity"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteUserData deletes all user data across various repositories
func (s *GrpcServer) DeleteUserData(ctx context.Context, req *connect.Request[productivity.DeleteUserDataRequest]) (*connect.Response[productivity.DeleteUserDataResponse], error) {
	user := req.Msg.GetUser()
	if user == nil {
		log.Error().Msg("User is required")
		return connect.NewResponse(&productivity.DeleteUserDataResponse{
			Success: false,
		}), nil
	}

	userIDHex := user.GetId()
	if userIDHex == "" {
		log.Error().Msg("User ID is required")
		return connect.NewResponse(&productivity.DeleteUserDataResponse{
			Success: false,
		}), nil
	}

	// Convert user ID string to ObjectID
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		log.Error().Err(err).Msg("Invalid user ID format")
		return connect.NewResponse(&productivity.DeleteUserDataResponse{
			Success: false,
		}), nil
	}

	// Delete all tasks for the user
	if err := s.taskRepo.DeleteByUserID(ctx, userID); err != nil {
		log.Error().Err(err).Msg("Failed to delete user tasks")
		return connect.NewResponse(&productivity.DeleteUserDataResponse{
			Success: false,
		}), nil
	}

	// Delete all habit entries for the user
	if err := s.habitRepo.DeleteEntriesByUserID(ctx, userID); err != nil {
		log.Error().Err(err).Msg("Failed to delete user habit entries")
		return connect.NewResponse(&productivity.DeleteUserDataResponse{
			Success: false,
		}), nil
	}

	// Delete all habits for the user
	if err := s.habitRepo.DeleteByUserID(ctx, userID); err != nil {
		log.Error().Err(err).Msg("Failed to delete user habits")
		return connect.NewResponse(&productivity.DeleteUserDataResponse{
			Success: false,
		}), nil
	}

	// Delete all notes for the user
	if err := s.noteRepo.DeleteByUserID(ctx, userID); err != nil {
		log.Error().Err(err).Msg("Failed to delete user notes")
		return connect.NewResponse(&productivity.DeleteUserDataResponse{
			Success: false,
		}), nil
	}

	// Delete all tags for the user
	if err := s.tagRepo.DeleteByUserID(ctx, userID); err != nil {
		log.Error().Err(err).Msg("Failed to delete user tags")
		return connect.NewResponse(&productivity.DeleteUserDataResponse{
			Success: false,
		}), nil
	}

	// Delete all folders for the user (this also cleans folder references from tasks)
	if err := s.folderRepo.DeleteByUserID(ctx, userID); err != nil {
		log.Error().Err(err).Msg("Failed to delete user folders")
		return connect.NewResponse(&productivity.DeleteUserDataResponse{
			Success: false,
		}), nil
	}

	// Delete all time entries for the user
	if err := s.timeEntryRepo.DeleteByUserID(ctx, userID); err != nil {
		log.Error().Err(err).Msg("Failed to delete user time entries")
		return connect.NewResponse(&productivity.DeleteUserDataResponse{
			Success: false,
		}), nil
	}

	log.Info().Str("userID", userIDHex).Msg("Successfully deleted user data")

	return connect.NewResponse(&productivity.DeleteUserDataResponse{
		Success: true,
	}), nil
}
