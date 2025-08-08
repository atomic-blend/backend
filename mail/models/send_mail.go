package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SendStatus represents the status of a sent mail
type SendStatus string

const (
	SendStatusPending SendStatus = "pending"
	SendStatusSent    SendStatus = "sent"
	SendStatusFailed  SendStatus = "failed"
	SendStatusRetry   SendStatus = "retry"
)

// SendMail represents a mail that has been sent or is queued to be sent
type SendMail struct {
	ID           primitive.ObjectID  `bson:"_id" json:"id"`
	Mail         *Mail               `bson:"mail" json:"mail"`
	SendStatus   SendStatus          `bson:"send_status" json:"send_status"`
	RetryCounter *int                `bson:"retry_counter" json:"retry_counter"`
	Trashed      bool                `bson:"trashed" json:"trashed"`
	CreatedAt    *primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt    *primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
