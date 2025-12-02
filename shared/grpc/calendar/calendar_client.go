package calendarclient

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	calendarv1 "github.com/atomic-blend/backend/grpc/gen/calendar/v1"
	"github.com/atomic-blend/backend/grpc/gen/calendar/v1/calendarv1connect"
	grpcclientutils "github.com/atomic-blend/backend/shared/utils/grpc_client_utils"
)

// CalendarClient is the client for calendar-related gRPC operations
type CalendarClient struct {
	client calendarv1connect.CalendarServiceClient
}

var _ Interface = (*CalendarClient)(nil)

// NewCalendarClient creates a new calendar client
func NewCalendarClient() (*CalendarClient, error) {
	httpClient := &http.Client{}
	baseURL, err := grpcclientutils.GetServiceBaseURL("calendar")
	if err != nil {
		return nil, err
	}

	client := calendarv1connect.NewCalendarServiceClient(httpClient, baseURL)
	return &CalendarClient{client: client}, nil
}

// DeleteUserData calls the DeleteUserData method on the calendar service
func (c *CalendarClient) DeleteUserData(ctx context.Context, req *connect.Request[calendarv1.DeleteUserDataRequest]) (*connect.Response[calendarv1.DeleteUserDataResponse], error) {
	return c.client.DeleteUserData(ctx, req)
}

// CreateCalendar calls the CreateCalendar method on the calendar service
func (c *CalendarClient) CreateCalendar(ctx context.Context, req *connect.Request[calendarv1.CreateCalendarRequest]) (*connect.Response[calendarv1.CreateCalendarResponse], error) {
	return c.client.CreateCalendar(ctx, req)
}
