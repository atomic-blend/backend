// Package mailserver provides the client for the mail-server service
package mailserver

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	mailserverv1 "github.com/atomic-blend/backend/grpc/gen/mailserver/v1"
	"github.com/atomic-blend/backend/grpc/gen/mailserver/v1/mailserverv1connect"
	grpcclientutils "github.com/atomic-blend/backend/shared/utils/grpc_client_utils"
)

// MailServerClient is the client for mail-server-related gRPC operations
type MailServerClient struct {
	client mailserverv1connect.MailServerServiceClient
}

var _ Interface = (*MailServerClient)(nil)

// NewMailServerClient creates a new mail-server client
func NewMailServerClient() (*MailServerClient, error) {
	httpClient := &http.Client{}
	baseURL, err := grpcclientutils.GetServiceBaseURL("mail-server")
	if err != nil {
		return nil, err
	}

	client := mailserverv1connect.NewMailServerServiceClient(httpClient, baseURL)
	return &MailServerClient{client: client}, nil
}

// SendMailInternal calls the SendMailInternal method on the mail-server service
func (m *MailServerClient) SendMailInternal(ctx context.Context, req *connect.Request[mailserverv1.SendMailInternalRequest]) (*connect.Response[mailserverv1.SendMailInternalResponse], error) {
	return m.client.SendMailInternal(ctx, req)
}