package server

import (
	"context"
	"fmt"
	"os"

	"connectrpc.com/connect"
	mailserverv1 "github.com/atomic-blend/backend/grpc/gen/mail-server/v1"
	mailsender "github.com/atomic-blend/backend/mail-server/utils/mail-sender"
	"github.com/atomic-blend/backend/mail/models"
)

// SendMailNoReply sends an email to the given recipients
// The email is sent as noreply@atomic-blend.com
// The email is sent to the given recipients
// The email is sent with the given subject, text content, and html content
func (s *GrpcServer) SendMailNoReply(ctx context.Context, req *connect.Request[mailserverv1.SendMailNoReplyRequest]) (*connect.Response[mailserverv1.SendMailNoReplyResponse], error) {
	noReplyEmail := "noreply@atomic-blend.com"
	if os.Getenv("NO_REPLY_EMAIL") != "" {
		noReplyEmail = os.Getenv("NO_REPLY_EMAIL")
	}
	mail := models.RawMail{
		Headers: map[string]interface{}{
			"To": req.Msg.To,
			"From": noReplyEmail,
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
	return nil, nil
}
