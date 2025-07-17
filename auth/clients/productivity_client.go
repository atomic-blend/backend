package clients

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	grpcclientutils "github.com/atomic-blend/backend/auth/utils/grpc_client_utils"
	"github.com/atomic-blend/backend/grpc/gen/productivity"
	"github.com/atomic-blend/backend/grpc/gen/productivity/productivityconnect"
)

// ProductivityClient wraps the real gRPC productivity client
type ProductivityClient struct {
	client productivityconnect.ProductivityServiceClient
}

// NewProductivityClient creates a new productivity client
func NewProductivityClient() (*ProductivityClient, error) {
	httpClient := &http.Client{}
	baseURL, err := grpcclientutils.GetServiceBaseURL("productivity")
	if err != nil {
		return nil, err
	}

	client := productivityconnect.NewProductivityServiceClient(httpClient, baseURL)
	return &ProductivityClient{client: client}, nil
}

// DeleteUserData calls the DeleteUserData method on the productivity service
func (p *ProductivityClient) DeleteUserData(ctx context.Context, req *connect.Request[productivity.DeleteUserDataRequest]) (*connect.Response[productivity.DeleteUserDataResponse], error) {
	return p.client.DeleteUserData(ctx, req)
}
