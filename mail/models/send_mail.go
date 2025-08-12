package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SendStatus represents the status of a sent mail
type SendStatus string

const (
	// SendStatusPending represents a mail that is pending to be sent
	SendStatusPending SendStatus = "pending"
	// SendStatusSent represents a mail that has been sent
	SendStatusSent    SendStatus = "sent"
	// SendStatusFailed represents a mail that has failed to be sent
	SendStatusFailed  SendStatus = "failed"
	// SendStatusRetry represents a mail that has failed to be sent and is being retried
	SendStatusRetry   SendStatus = "retry"
)

// SendMail represents a mail that has been sent or is queued to be sent
type SendMail struct {
	ID           primitive.ObjectID  `bson:"_id" json:"id"`
	Mail         *Mail               `bson:"mail" json:"mail"`
	SendStatus   SendStatus          `bson:"send_status" json:"send_status"`
	RetryCounter *int                `bson:"retry_counter,omitempty" json:"retry_counter,omitempty"`
	FailureReason *string             `bson:"failure_reason,omitempty" json:"failure_reason,omitempty"`
	FailedAt      *primitive.DateTime `bson:"failed_at,omitempty" json:"failed_at,omitempty"`
	Trashed      bool                `bson:"trashed" json:"trashed"`
	CreatedAt    *primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt    *primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
