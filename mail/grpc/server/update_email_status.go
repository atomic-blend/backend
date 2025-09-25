package global

import (
	"context"
	"errors"
	"time"

	"connectrpc.com/connect"
	mailv1 "github.com/atomic-blend/backend/grpc/gen/mail/v1"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateMailStatus updates the status of an email
func (s *GrpcServer) UpdateMailStatus(ctx context.Context, req *connect.Request[mailv1.UpdateMailStatusRequest]) (*connect.Response[mailv1.UpdateMailStatusResponse], error) {
	emailID := req.Msg.GetEmailId()
	status := req.Msg.GetStatus()
	failureReason := req.Msg.GetFailureReason()
	failedAt := req.Msg.GetFailedAt()
	retryCounter := req.Msg.GetRetryCounter()

	if emailID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("email_id is required"))
	}

	if status == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("status is required"))
	}

	if status == "failed" && failureReason == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("failure_reason is required when status is failed"))
	}

	if status == "failed" && failedAt == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("failed_at is required when status is failed"))
	}

	if status == "failed" && retryCounter <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("retry_counter must be greater than 0 when status is failed"))
	}

	//TODO: update the email status in the database
	sendEmailID, err := primitive.ObjectIDFromHex(emailID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid email_id"))
	}

	log.Debug().Interface("sendEmailId", sendEmailID).Msg("Updating email status")

	update := bson.M{
		"send_status": status,
	}

	// Only add optional fields if they are defined
	if failureReason != "" {
		update["failure_reason"] = failureReason
	}
	if failedAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, failedAt)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid failed_at format"))
		}
		update["failed_at"] = primitive.NewDateTimeFromTime(parsedTime)
	}
	if retryCounter > 0 {
		update["retry_counter"] = retryCounter
	}

	_, err = s.sendMailRepository.Update(ctx, sendEmailID, update)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	log.Debug().Msg("Updated email status successfully")

	return connect.NewResponse(&mailv1.UpdateMailStatusResponse{
		Success: true,
	}), nil
}
