package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserResetPasswordRequest struct {
	UserID    *primitive.ObjectID `json:"id" bson:"_id"`
	ResetCode string              `json:"reset_code" bson:"reset_code"`
	CreatedAt string              `json:"created_at" bson:"created_at"`
	UpdatedAt string              `json:"updated_at" bson:"updated_at"`
}
