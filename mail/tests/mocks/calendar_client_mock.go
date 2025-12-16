package mocks

import (
	"context"

	"connectrpc.com/connect"
	calendarv1 "github.com/atomic-blend/backend/grpc/gen/calendar/v1"
	"github.com/stretchr/testify/mock"
)

// MockCalendarClient provides a mock implementation of CalendarClient
type MockCalendarClient struct {
	mock.Mock
}

// DeleteUserData deletes user data
func (m *MockCalendarClient) DeleteUserData(ctx context.Context, req *connect.Request[calendarv1.DeleteUserDataRequest]) (*connect.Response[calendarv1.DeleteUserDataResponse], error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*connect.Response[calendarv1.DeleteUserDataResponse]), args.Error(1)
}

// CreateCalendar creates a calendar
func (m *MockCalendarClient) CreateCalendar(ctx context.Context, req *connect.Request[calendarv1.CreateCalendarRequest]) (*connect.Response[calendarv1.CreateCalendarResponse], error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*connect.Response[calendarv1.CreateCalendarResponse]), args.Error(1)
}
