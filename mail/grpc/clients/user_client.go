package clients

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	userv1 "github.com/atomic-blend/backend/grpc/gen/user/v1"
	"github.com/atomic-blend/backend/grpc/gen/user/v1/userv1connect"
	grpcclientutils "github.com/atomic-blend/backend/mail/utils/grpc_client_utils"
)

// UserClient is the client for user-related gRPC operations
type UserClient struct {
	client userv1connect.UserServiceClient
}

// NewUserClient creates a new user client
func NewUserClient() (*UserClient, error) {
	httpClient := &http.Client{}
	baseURL, err := grpcclientutils.GetServiceBaseURL("auth")
	if err != nil {
		return nil, err
	}

	client := userv1connect.NewUserServiceClient(httpClient, baseURL)
	return &UserClient{client: client}, nil
}

// GetUserDevices calls the GetUserDevices method on the user service
func (u *UserClient) GetUserDevices(ctx context.Context, req *connect.Request[userv1.GetUserDevicesRequest]) (*connect.Response[userv1.GetUserDevicesResponse], error) {
	return u.client.GetUserDevices(ctx, req)
}
