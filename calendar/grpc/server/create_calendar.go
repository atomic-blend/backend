package global

import (
	"context"

	"connectrpc.com/connect"
	calendarv1 "github.com/atomic-blend/backend/grpc/gen/calendar/v1"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateCalendar handles the creation of a new calendar
func (s *GrpcServer) CreateCalendar(ctx context.Context, req *connect.Request[calendarv1.CreateCalendarRequest]) (*connect.Response[calendarv1.CreateCalendarResponse], error) {
	user := req.Msg.GetUser()
	if user == nil {
		log.Error().Msg("User is required")
		return connect.NewResponse(&calendarv1.CreateCalendarResponse{
			Id:    nil,
			Error: prt("User is required"),
		}), nil
	}

	userIDHex := user.GetId()
	if userIDHex == "" {
		log.Error().Msg("User ID is required")
		return connect.NewResponse(&calendarv1.CreateCalendarResponse{
			Id:    nil,
			Error: prt("User ID is required"),
		}), nil
	}

	// Convert user ID string to ObjectID
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		log.Error().Err(err).Msg("Invalid user ID format")
		return connect.NewResponse(&calendarv1.CreateCalendarResponse{
			Id:    nil,
			Error: prt("Invalid user ID format"),
		}), nil
	}

	log.Info().Str("userID", userID.Hex()).Msg("Creating calendar event for user")

	// TODO: create the calendar event inside the repo

	return connect.NewResponse(&calendarv1.CreateCalendarResponse{
		Id: nil,
	}), nil
}

func prt(s string) *string {
	return &s
}
