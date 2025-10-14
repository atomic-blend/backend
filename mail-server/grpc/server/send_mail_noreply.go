package server

import (
	"context"

	"connectrpc.com/connect"
	mailserverv1 "github.com/atomic-blend/backend/grpc/gen/mail-server/v1"
)

func (s *GrpcServer) SendMailNoReply(ctx context.Context, req *connect.Request[mailserverv1.SendMailNoReplyRequest]) (*connect.Response[mailserverv1.SendMailNoReplyResponse], error) {
	//TODO: send the email as noreply@atomic-blend.com
	return nil, nil
}
