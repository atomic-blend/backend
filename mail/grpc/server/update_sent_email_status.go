package grpserver

import (
	"context"
	"fmt"

	"errors"

	"connectrpc.com/connect"
	mailv1 "github.com/atomic-blend/backend/grpc/gen/mail/v1"
	"github.com/atomic-blend/backend/mail/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *GrpcServer) UpdateSentEmailStatus(ctx context.Context, req *connect.Request[mailv1.UpdateSentEmailStatusRequest]) (*connect.Response[mailv1.UpdateSentEmailStatusResponse], error) {
	//TODO: update the status of the email
	//TODO: remove the claimed data if set
	idHex := req.Msg.GetEmailId()
	if idHex == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("id is required"))
	}

	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid id format: %w", err))
	}

	status := req.Msg.GetStatus()
	if status == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("status is required"))
	}

	sendStatus := models.SendStatus(status)

	_, err = s.sendMailRepo.UpdateStatus(ctx, id, sendStatus)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update email status: %w", err))
	}

	return connect.NewResponse(&mailv1.UpdateSentEmailStatusResponse{
		Success: true,
	}), nil
}
