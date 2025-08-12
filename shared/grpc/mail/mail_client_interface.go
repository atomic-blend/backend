// Package mailclient contains the interfaces for the mail-related gRPC operations
package mailclient

import (
	"context"

	"connectrpc.com/connect"
	mailv1 "github.com/atomic-blend/backend/grpc/gen/mail/v1"
)

// Interface defines the methods for mail-related gRPC operations
type Interface interface {
	UpdateMailStatus(context.Context, *connect.Request[mailv1.UpdateMailStatusRequest]) (*connect.Response[mailv1.UpdateMailStatusResponse], error)
}
