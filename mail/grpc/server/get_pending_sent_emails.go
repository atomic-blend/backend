package grpserver

import (
	"context"

	"connectrpc.com/connect"
	mailv1 "github.com/atomic-blend/backend/grpc/gen/mail/v1"
)

// GetPendingSentEmails retrieves pending or retrying sent emails for a user
func (s *GrpcServer) GetPendingSentEmails(ctx context.Context, req *connect.Request[mailv1.GetPendingSentEmailsRequest]) (*connect.Response[mailv1.GetPendingSentEmailsResponse], error) {
	pendingEmails, err := s.sendMailRepo.ClaimPendingSentEmails(ctx)
	if err != nil {
		return nil, err
	}

	// convert pending emails to the response format
	var emails []*mailv1.Mail
	for _, email := range pendingEmails {
		// Create headers map for the Mail proto
		headers := make(map[string]string)
		if email.Mail != nil && email.Mail.Headers != nil {
			if headersMap, ok := email.Mail.Headers.(map[string]interface{}); ok {
				for k, v := range headersMap {
					if strVal, ok := v.(string); ok {
						headers[k] = strVal
					}
				}
			}
		}

		// Use HTMLContent if available, otherwise TextContent
		htmlContent := ""
		textContent := ""
		if email.Mail != nil {
			htmlContent = email.Mail.HTMLContent
			textContent = email.Mail.TextContent
		}

		emails = append(emails, &mailv1.Mail{
			Id:          email.ID.Hex(),
			Headers:     headers,
			HtmlContent: htmlContent,
			TextContent: textContent,
			//TODO: get the files from s3
			Attachments: []*mailv1.MailAttachment{},
		})
	}

	return connect.NewResponse(&mailv1.GetPendingSentEmailsResponse{
		Emails: emails,
	}), nil
}
