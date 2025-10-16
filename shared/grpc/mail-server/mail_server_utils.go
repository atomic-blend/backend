package mailserver

import (
	"connectrpc.com/connect"
	mailserverv1 "github.com/atomic-blend/backend/grpc/gen/mailserver/v1"
)

// CreateSendMailInternalRequest creates a new SendMailInternalRequest with the given parameters
func CreateSendMailInternalRequest(to []string, from, subject, htmlContent, textContent string) *connect.Request[mailserverv1.SendMailInternalRequest] {
	req := &mailserverv1.SendMailInternalRequest{
		To:            to,
		From:          from,
		Subject:       subject,
		HtmlContent:   htmlContent,
		TextContent:   textContent,
	}

	return connect.NewRequest(req)
}