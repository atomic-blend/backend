package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Tag represents a tag associated with a task
type Tag struct {
	ID        *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    *primitive.ObjectID `json:"userId" bson:"user_id,omitempty"`
	Name      string              `json:"name" bson:"name" binding:"required"`
	Color     *string             `json:"color" bson:"color"`
	CreatedAt *primitive.DateTime `json:"createdAt" bson:"created_at"`
	UpdatedAt *primitive.DateTime `json:"updatedAt" bson:"updated_at"`
}
