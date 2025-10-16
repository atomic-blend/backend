package server

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	mailserverv1 "github.com/atomic-blend/backend/grpc/gen/mailserver/v1"
	mailsender "github.com/atomic-blend/backend/mail-server/utils/mail-sender"
	"github.com/atomic-blend/backend/mail/models"
)

// SendMailInternal sends an email to the given recipients
// The email is sent as noreply@atomic-blend.com
// The email is sent to the given recipients
// The email is sent with the given subject, text content, and html content
func (s *GrpcServer) SendMailInternal(ctx context.Context, req *connect.Request[mailserverv1.SendMailInternalRequest]) (*connect.Response[mailserverv1.SendMailInternalResponse], error) {
	mail := models.RawMail{
		Headers: map[string]any{
			"To":      req.Msg.To,
			"From":    req.Msg.From,
			"Subject": req.Msg.Subject,
		},
		TextContent: req.Msg.TextContent,
		HTMLContent: req.Msg.HtmlContent,
	}
	recipientsToRetry, err := mailsender.SendEmail(mail, nil)
	if err != nil {
		return nil, err
	}
	if len(recipientsToRetry) > 0 {
		return nil, fmt.Errorf("failed to send email to all recipients")
	}
	return connect.NewResponse(&mailserverv1.SendMailInternalResponse{
		Success: true,
	}), nil
}
