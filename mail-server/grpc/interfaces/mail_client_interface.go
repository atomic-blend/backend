// Package interfaces contains the interfaces for the mail-related gRPC operations
package interfaces

import (
	"context"

	"connectrpc.com/connect"
	mailv1 "github.com/atomic-blend/backend/grpc/gen/mail/v1"
)

// MailClientInterface defines the methods for mail-related gRPC operations
type MailClientInterface interface {
	UpdateMailStatus(context.Context, *connect.Request[mailv1.UpdateMailStatusRequest]) (*connect.Response[mailv1.UpdateMailStatusResponse], error)
}
