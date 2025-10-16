package mailserver

import (
	"context"

	"connectrpc.com/connect"
	mailserverv1 "github.com/atomic-blend/backend/grpc/gen/mailserver/v1"
)

// Interface defines the methods for mail-server-related gRPC operations
type Interface interface {
	SendMailInternal(ctx context.Context, req *connect.Request[mailserverv1.SendMailInternalRequest]) (*connect.Response[mailserverv1.SendMailInternalResponse], error)
}