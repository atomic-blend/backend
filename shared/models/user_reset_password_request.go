package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// UserResetPassword represents a user's reset password request
// It contains the user ID, reset code, and timestamps for creation and update
type UserResetPassword struct {
	UserID    *primitive.ObjectID `json:"id" bson:"_id"`
	ResetCode string              `json:"reset_code" bson:"reset_code"`
	CreatedAt primitive.DateTime  `json:"created_at" bson:"created_at"`
	UpdatedAt primitive.DateTime  `json:"updated_at" bson:"updated_at"`
}
