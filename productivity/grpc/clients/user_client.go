package clients

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/atomic-blend/backend/grpc/gen/auth"
	"github.com/atomic-blend/backend/grpc/gen/auth/authconnect"
	grpcclientutils "github.com/atomic-blend/backend/productivity/utils/grpc_client_utils"
)

// UserClient is the client for user-related gRPC operations
type UserClient struct {
	client authconnect.UserServiceClient
}

// NewUserClient creates a new user client
func NewUserClient() (*UserClient, error) {
	httpClient := &http.Client{}
	baseURL, err := grpcclientutils.GetServiceBaseURL("auth")
	if err != nil {
		return nil, err
	}

	client := authconnect.NewUserServiceClient(httpClient, baseURL)
	return &UserClient{client: client}, nil
}

// GetUserDevices calls the GetUserDevices method on the user service
func (u *UserClient) GetUserDevices(ctx context.Context, req *connect.Request[auth.GetUserDevicesRequest]) (*connect.Response[auth.GetUserDevicesResponse], error) {
	return u.client.GetUserDevices(ctx, req)
}
