// Package client contains the logic for creating mail status requests
package client

import (
	"time"

	"connectrpc.com/connect"
	mailv1 "github.com/atomic-blend/backend/grpc/gen/mail/v1"
)

// CreateUpdateMailStatusRequest creates a new UpdateMailStatusRequest with the given parameters
func CreateUpdateMailStatusRequest(emailID, status string, failureReason *string, retryCounter *int32) *connect.Request[mailv1.UpdateMailStatusRequest] {
	req := &mailv1.UpdateMailStatusRequest{
		EmailId: emailID,
		Status:  status,
	}

	// Only set FailedAt when there's a failure reason
	if failureReason != nil {
		req.FailureReason = failureReason
		failedAt := time.Now().Format(time.RFC3339)
		req.FailedAt = &failedAt
	}

	if retryCounter != nil {
		req.RetryCounter = retryCounter
	}

	return connect.NewRequest(req)
}

// CreateSuccessStatusRequest creates a request for successful mail delivery
func CreateSuccessStatusRequest(emailID string) *connect.Request[mailv1.UpdateMailStatusRequest] {
	return CreateUpdateMailStatusRequest(emailID, "sent", nil, nil)
}

// CreateFailureStatusRequest creates a request for failed mail delivery
func CreateFailureStatusRequest(emailID, failureReason string, retryCounter int32) *connect.Request[mailv1.UpdateMailStatusRequest] {
	return CreateUpdateMailStatusRequest(emailID, "failed", &failureReason, &retryCounter)
}

// CreateRetryStatusRequest creates a request for mail that needs to be retried
func CreateRetryStatusRequest(emailID, failureReason string, retryCounter int32) *connect.Request[mailv1.UpdateMailStatusRequest] {
	return CreateUpdateMailStatusRequest(emailID, "retry", &failureReason, &retryCounter)
}
