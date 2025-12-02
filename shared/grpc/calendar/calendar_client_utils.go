// Package calendarclient contains the logic for creating calendar requests
package calendarclient

import (
	"connectrpc.com/connect"
	authv1 "github.com/atomic-blend/backend/grpc/gen/auth/v1"
	calendarv1 "github.com/atomic-blend/backend/grpc/gen/calendar/v1"
)

// CreateDeleteUserDataRequest creates a new DeleteUserDataRequest with the given user
func CreateDeleteUserDataRequest(user *authv1.User) *connect.Request[calendarv1.DeleteUserDataRequest] {
	req := &calendarv1.DeleteUserDataRequest{
		User: user,
	}
	return connect.NewRequest(req)
}

// CreateCreateCalendarRequest creates a new CreateCalendarRequest with the given user and calendar
func CreateCreateCalendarRequest(user *authv1.User, calendar *calendarv1.Calendar) *connect.Request[calendarv1.CreateCalendarRequest] {
	req := &calendarv1.CreateCalendarRequest{
		User:     user,
		Calendar: calendar,
	}
	return connect.NewRequest(req)
}
