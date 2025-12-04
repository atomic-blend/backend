package global

import (
	"context"

	"connectrpc.com/connect"
	"github.com/atomic-blend/backend/calendar/models"
	"github.com/atomic-blend/backend/calendar/repositories"
	calendarv1 "github.com/atomic-blend/backend/grpc/gen/calendar/v1"
	"github.com/atomic-blend/backend/shared/utils/db"
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

	// Get the calendar from the request
	calendarProto := req.Msg.GetCalendar()
	if calendarProto == nil {
		log.Error().Msg("Calendar is required")
		return connect.NewResponse(&calendarv1.CreateCalendarResponse{
			Id:    nil,
			Error: prt("Calendar is required"),
		}), nil
	}

	// Convert protobuf calendar to models calendar
	calendarModel, err := convertProtoToModel(calendarProto)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert calendar")
		return connect.NewResponse(&calendarv1.CreateCalendarResponse{
			Id:    nil,
			Error: prt("Failed to convert calendar"),
		}), nil
	}

	// Set the user ID
	calendarModel.UserID = &userID

	// Create the calendar in the repository
	calendarRepo := repositories.NewCalendarRepository(db.Database)
	createdCalendar, err := calendarRepo.CreateWithContext(ctx, calendarModel)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create calendar")
		return connect.NewResponse(&calendarv1.CreateCalendarResponse{
			Id:    nil,
			Error: prt("Failed to create calendar"),
		}), nil
	}

	calendarID := createdCalendar.ID.Hex()
	return connect.NewResponse(&calendarv1.CreateCalendarResponse{
		Id: &calendarID,
	}), nil
}

func prt(s string) *string {
	return &s
}

// convertProtoToModel converts a calendarv1.Calendar to models.Calendar
func convertProtoToModel(proto *calendarv1.Calendar) (*models.Calendar, error) {
	model := &models.Calendar{
		ProdID:  proto.ProdId,
		Version: proto.Version,
		Method:  proto.Method,
	}

	// Convert events
	if len(proto.Events) > 0 {
		model.Events = make([]models.Event, len(proto.Events))
		for i, eventProto := range proto.Events {
			event, err := convertEventProtoToModel(eventProto)
			if err != nil {
				return nil, err
			}
			model.Events[i] = *event
		}
	}

	// Convert timezones
	if len(proto.Timezones) > 0 {
		model.Timezones = make([]models.VTimezone, len(proto.Timezones))
		for i, tzProto := range proto.Timezones {
			model.Timezones[i] = models.VTimezone{
				TZID:   tzProto.Tzid,
				Offset: int(tzProto.Offset),
			}
		}
	}

	// Convert timestamps
	if proto.CreatedAt != nil {
		createdAt := proto.CreatedAt.AsTime()
		createdAtDateTime := primitive.NewDateTimeFromTime(createdAt)
		model.CreatedAt = &createdAtDateTime
	}
	if proto.UpdatedAt != nil {
		updatedAt := proto.UpdatedAt.AsTime()
		updatedAtDateTime := primitive.NewDateTimeFromTime(updatedAt)
		model.UpdatedAt = &updatedAtDateTime
	}

	return model, nil
}

// convertEventProtoToModel converts a calendarv1.Event to models.Event
func convertEventProtoToModel(proto *calendarv1.Event) (*models.Event, error) {
	event := &models.Event{
		UID:         proto.Uid,
		Summary:     proto.Summary,
		Description: proto.Description,
		Location:    &proto.Location,
		Status:      proto.Status,
		Sequence:    int(proto.Sequence),
	}

	// Convert timestamps
	if proto.DtStamp != nil {
		event.DTStamp = proto.DtStamp.AsTime()
	}
	if proto.Start != nil {
		startTime := proto.Start.Time.AsTime()
		event.Start = &models.TimeSpec{
			Time: startTime,
			TZID: proto.Start.Tzid,
		}
	}
	if proto.End != nil {
		endTime := proto.End.Time.AsTime()
		event.End = &models.TimeSpec{
			Time: endTime,
			TZID: proto.End.Tzid,
		}
	}
	if proto.Duration != nil {
		duration := proto.Duration.AsDuration()
		event.Duration = &duration
	}
	if proto.LastModified != nil {
		lastModified := proto.LastModified.AsTime()
		event.LastModified = &lastModified
	}
	if proto.CreatedAt != nil {
		createdAt := proto.CreatedAt.AsTime()
		createdAtDateTime := primitive.NewDateTimeFromTime(createdAt)
		event.CreatedAt = &createdAtDateTime
	}
	if proto.UpdatedAt != nil {
		updatedAt := proto.UpdatedAt.AsTime()
		updatedAtDateTime := primitive.NewDateTimeFromTime(updatedAt)
		event.UpdatedAt = &updatedAtDateTime
	}

	// Convert organizer
	if proto.Organizer != nil {
		event.Organizer = &models.CalAddress{
			Email:    proto.Organizer.Email,
			CN:       proto.Organizer.Cn,
			Role:     proto.Organizer.Role,
			PartStat: proto.Organizer.PartStat,
			RSVP:     proto.Organizer.Rsvp,
		}
	}

	// Convert attendees
	if len(proto.Attendees) > 0 {
		event.Attendees = make([]models.CalAddress, len(proto.Attendees))
		for i, attendeeProto := range proto.Attendees {
			event.Attendees[i] = models.CalAddress{
				Email:    attendeeProto.Email,
				CN:       attendeeProto.Cn,
				Role:     attendeeProto.Role,
				PartStat: attendeeProto.PartStat,
				RSVP:     attendeeProto.Rsvp,
			}
		}
	}

	// Convert recurrence
	if proto.Recurrence != nil {
		event.Recurrence = &models.RRule{
			Freq:     proto.Recurrence.Freq,
			Interval: int(proto.Recurrence.Interval),
			Count:    int(proto.Recurrence.Count),
		}
		if proto.Recurrence.Until != nil {
			until := proto.Recurrence.Until.AsTime()
			event.Recurrence.Until = &until
		}
		if len(proto.Recurrence.ByDay) > 0 {
			event.Recurrence.ByDay = make([]string, len(proto.Recurrence.ByDay))
			copy(event.Recurrence.ByDay, proto.Recurrence.ByDay)
		}
		if len(proto.Recurrence.ByMonth) > 0 {
			event.Recurrence.ByMonth = make([]int, len(proto.Recurrence.ByMonth))
			for i, month := range proto.Recurrence.ByMonth {
				event.Recurrence.ByMonth[i] = int(month)
			}
		}
		if len(proto.Recurrence.BySetPos) > 0 {
			event.Recurrence.BySetPos = make([]int, len(proto.Recurrence.BySetPos))
			for i, pos := range proto.Recurrence.BySetPos {
				event.Recurrence.BySetPos[i] = int(pos)
			}
		}
	}

	// Convert exception dates
	if len(proto.ExDates) > 0 {
		event.ExDates = make([]models.TimeSpec, len(proto.ExDates))
		for i, exDateProto := range proto.ExDates {
			exDateTime := exDateProto.Time.AsTime()
			event.ExDates[i] = models.TimeSpec{
				Time: exDateTime,
				TZID: exDateProto.Tzid,
			}
		}
	}

	// Convert recurrence dates
	if len(proto.RDates) > 0 {
		event.RDates = make([]models.TimeSpec, len(proto.RDates))
		for i, rDateProto := range proto.RDates {
			rDateTime := rDateProto.Time.AsTime()
			event.RDates[i] = models.TimeSpec{
				Time: rDateTime,
				TZID: rDateProto.Tzid,
			}
		}
	}

	return event, nil
}
