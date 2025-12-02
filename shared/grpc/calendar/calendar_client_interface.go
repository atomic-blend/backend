// Package calendarclient contains the interfaces for the calendar-related gRPC operations
package calendarclient

import (
	"context"

	"connectrpc.com/connect"
	calendarv1 "github.com/atomic-blend/backend/grpc/gen/calendar/v1"
)

// Interface defines the methods for calendar-related gRPC operations
type Interface interface {
	DeleteUserData(context.Context, *connect.Request[calendarv1.DeleteUserDataRequest]) (*connect.Response[calendarv1.DeleteUserDataResponse], error)
	CreateCalendar(context.Context, *connect.Request[calendarv1.CreateCalendarRequest]) (*connect.Response[calendarv1.CreateCalendarResponse], error)
}