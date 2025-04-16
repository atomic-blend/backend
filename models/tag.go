package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tag struct {
	ID        *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    *primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	Name      string              `json:"name" bson:"name" binding:"required"`
	Color     *string             `json:"color" bson:"color"`
	CreatedAt *primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt *primitive.DateTime `json:"updated_at" bson:"updated_at"`
}
