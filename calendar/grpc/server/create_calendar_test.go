package global

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/atomic-blend/backend/grpc/gen/auth/v1"
	"github.com/atomic-blend/backend/grpc/gen/calendar/v1"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"
	"github.com/atomic-blend/backend/shared/utils/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func setupGrpcTest(t *testing.T) (*GrpcServer, func()) {
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	client, err := inmemorymongo.ConnectToInMemoryDB(mongoServer.URI())
	require.NoError(t, err)

	testDB := client.Database("test_db")
	// Set the global Database to test db
	db.Database = testDB

	server := NewGrpcServer()

	cleanup := func() {
		client.Disconnect(context.Background())
		mongoServer.Stop()
		// Reset global db
		db.Database = nil
	}

	return server, cleanup
}

func TestCreateCalendar(t *testing.T) {
	server, cleanup := setupGrpcTest(t)
	defer cleanup()

	userID := primitive.NewObjectID()

	// Create a test calendar proto
	calendarProto := &calendarv1.Calendar{
		ProdId:  "-//AB//EN",
		Version: "2.0",
		Events:  []*calendarv1.Event{}, // empty for simplicity
	}

	userProto := &authv1.User{
		Id: userID.Hex(),
	}

	req := &connect.Request[calendarv1.CreateCalendarRequest]{
		Msg: &calendarv1.CreateCalendarRequest{
			User:     userProto,
			Calendar: calendarProto,
		},
	}

	resp, err := server.CreateCalendar(context.Background(), req)

	require.NoError(t, err)
	require.Nil(t, resp.Msg.Id)
	assert.NotEmpty(t, resp.Msg.Error)
	assert.Equal(t, "no_events_to_create", *resp.Msg.Error)

	// You can add more assertions here, like checking the database directly
	// But for now, just check response
}

func TestCreateCalendar_InvalidUserID(t *testing.T) {
	server, cleanup := setupGrpcTest(t)
	defer cleanup()

	// Create a test calendar proto
	calendarProto := &calendarv1.Calendar{
		ProdId:  "-//AB//EN",
		Version: "2.0",
		Events:  []*calendarv1.Event{},
	}

	userProto := &authv1.User{
		Id: "invalid-id",
	}

	req := &connect.Request[calendarv1.CreateCalendarRequest]{
		Msg: &calendarv1.CreateCalendarRequest{
			User:     userProto,
			Calendar: calendarProto,
		},
	}

	resp, err := server.CreateCalendar(context.Background(), req)

	require.NoError(t, err)
	assert.Nil(t, resp.Msg.Id)
	assert.NotEmpty(t, resp.Msg.Error)
	assert.Equal(t, "Invalid user ID format", *resp.Msg.Error)
}

func TestCreateCalendar_WithEvent(t *testing.T) {
	server, cleanup := setupGrpcTest(t)
	defer cleanup()

	userID := primitive.NewObjectID()

	// Create a test event proto
	eventProto := &calendarv1.Event{
		Uid:         "test-uid",
		Summary:     "Team Meeting",
		Description: "Weekly team sync",
		Location:    "Conference Room A",
		Status:      "CONFIRMED",
		Sequence:    0,
		Start: &calendarv1.TimeSpec{
			Time: &timestamppb.Timestamp{Seconds: 1733107200}, // 2024-12-02 00:00:00 UTC
			Tzid: "UTC",
		},
		End: &calendarv1.TimeSpec{
			Time: &timestamppb.Timestamp{Seconds: 1733110800}, // 2024-12-02 01:00:00 UTC
			Tzid: "UTC",
		},
		Organizer: &calendarv1.CalAddress{
			Email:    "organizer@example.com",
			Cn:       "John Doe",
			Role:     "CHAIR",
			PartStat: "ACCEPTED",
			Rsvp:     true,
		},
		Attendees: []*calendarv1.CalAddress{
			{
				Email:    "attendee1@example.com",
				Cn:       "Jane Smith",
				Role:     "REQ-PARTICIPANT",
				PartStat: "NEEDS-ACTION",
				Rsvp:     true,
			},
		},
	}

	// Create a test calendar proto with the event
	calendarProto := &calendarv1.Calendar{
		ProdId:  "-//AB//EN",
		Version: "2.0",
		Events:  []*calendarv1.Event{eventProto},
	}

	userProto := &authv1.User{
		Id: userID.Hex(),
	}

	req := &connect.Request[calendarv1.CreateCalendarRequest]{
		Msg: &calendarv1.CreateCalendarRequest{
			User:     userProto,
			Calendar: calendarProto,
		},
	}

	resp, err := server.CreateCalendar(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp.Msg.Id)
	assert.Empty(t, resp.Msg.Error)
}
