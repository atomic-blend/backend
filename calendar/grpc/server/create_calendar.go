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
			Error: ptr("User is required"),
		}), nil
	}

	userIDHex := user.GetId()
	if userIDHex == "" {
		log.Error().Msg("User ID is required")
		return connect.NewResponse(&calendarv1.CreateCalendarResponse{
			Id:    nil,
			Error: ptr("User ID is required"),
		}), nil
	}

	// Convert user ID string to ObjectID
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		log.Error().Err(err).Msg("Invalid user ID format")
		return connect.NewResponse(&calendarv1.CreateCalendarResponse{
			Id:    nil,
			Error: ptr("Invalid user ID format"),
		}), nil
	}

	log.Info().Str("userID", userID.Hex()).Msg("Creating calendar event for user")

	// Get the calendar from the request
	calendarProto := req.Msg.GetCalendar()
	if calendarProto == nil {
		log.Error().Msg("Calendar is required")
		return connect.NewResponse(&calendarv1.CreateCalendarResponse{
			Id:    nil,
			Error: ptr("Calendar is required"),
		}), nil
	}

	// Convert protobuf calendar to models calendar
	calendarModel, err := convertProtoToModel(calendarProto)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert calendar")
		return connect.NewResponse(&calendarv1.CreateCalendarResponse{
			Id:    nil,
			Error: ptr("Failed to convert calendar"),
		}), nil
	}

	var defaultCalendar *models.Calendar
	calendarRepo := repositories.NewCalendarRepository(db.Database)

	// if a user don't have a default calendar, create one
	defaultCalendar, err = calendarRepo.GetByNameWithContext(ctx, userID, nil)
	if err != nil || defaultCalendar == nil {
		if err != nil {
			log.Info().Err(err).Msg("Error finding default calendar, creating one")
		} else {
			log.Info().Msg("No default calendar found, creating one")
		}
		newDefaultCalendar := &models.Calendar{
			UserID: &userID,
		}
		defaultCalendar, err = calendarRepo.CreateWithContext(ctx, newDefaultCalendar)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create default calendar")
			return connect.NewResponse(&calendarv1.CreateCalendarResponse{
				Id:    nil,
				Error: ptr("Failed to create default calendar"),
			}), nil
		}
	}

	log.Info().Str("calendarID", defaultCalendar.ID.Hex()).Msg("Default calendar found")

	if len(calendarModel.Events) == 0 {
		log.Error().Msg("No events to create")
		return connect.NewResponse(&calendarv1.CreateCalendarResponse{
			Id:    nil,
			Error: ptr("no_events_to_create"),
		}), nil
	}

	// if event with same UID exists, update it instead and return existing one
	eventRepo := repositories.NewEventRepository(db.Database)
	existingEvent, _ := eventRepo.GetByUIDWithContext(ctx, calendarModel.Events[0].UID)
	var newEvent *models.Event

	if existingEvent != nil {
		log.Info().Str("eventUID", existingEvent.UID).Msg("Event with same UID exists, updating it")
		_, err = eventRepo.UpdateWithContext(ctx, existingEvent.ID, &calendarModel.Events[0])
		if err != nil {
			log.Error().Err(err).Msg("Failed to update existing event")
			return connect.NewResponse(&calendarv1.CreateCalendarResponse{
				Id:    nil,
				Error: ptr("Failed to update existing event"),
			}), nil
		}
		updated := true
		return connect.NewResponse(&calendarv1.CreateCalendarResponse{
			Id:      ptr(existingEvent.ID.Hex()),
			Updated: &updated,
		}), nil
	}

	log.Info().Msg("Creating new event")
	event := &calendarModel.Events[0]
	event.CalendarID = defaultCalendar.ID
	event.ID = primitive.NewObjectID()
	event.UserID = userID
	newEvent, err = eventRepo.CreateWithContext(ctx, event)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create event")
		return connect.NewResponse(&calendarv1.CreateCalendarResponse{
			Id:    nil,
			Error: ptr("Failed to create event"),
		}), nil
	}

	var returnEventID *string
	if existingEvent != nil {
		returnEventID = ptr(existingEvent.ID.Hex())
	} else {
		returnEventID = ptr(newEvent.ID.Hex())
	}
	return connect.NewResponse(&calendarv1.CreateCalendarResponse{
		Id: returnEventID,
	}), nil
}

func ptr(s string) *string {
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
