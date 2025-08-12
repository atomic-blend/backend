package mailclient

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	mailv1 "github.com/atomic-blend/backend/grpc/gen/mail/v1"
	"github.com/atomic-blend/backend/grpc/gen/mail/v1/mailv1connect"
	grpcclientutils "github.com/atomic-blend/backend/shared/utils/grpc_client_utils"
)

// MailClient is the client for mail-related gRPC operations
type MailClient struct {
	client mailv1connect.MailServiceClient
}

var _ Interface = (*MailClient)(nil)

// NewMailClient creates a new mail client
func NewMailClient() (*MailClient, error) {
	httpClient := &http.Client{}
	baseURL, err := grpcclientutils.GetServiceBaseURL("mail")
	if err != nil {
		return nil, err
	}

	client := mailv1connect.NewMailServiceClient(httpClient, baseURL)
	return &MailClient{client: client}, nil
}

// UpdateMailStatus calls the UpdateMailStatus method on the mail service
func (m *MailClient) UpdateMailStatus(ctx context.Context, req *connect.Request[mailv1.UpdateMailStatusRequest]) (*connect.Response[mailv1.UpdateMailStatusResponse], error) {
	return m.client.UpdateMailStatus(ctx, req)
}
