package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// NoteEntity represents a note in the system
// @Summary Note entity
// @Description Represents a note in the system
type NoteEntity struct {
	ID        *primitive.ObjectID `json:"id" bson:"_id"`
	Title     *string             `json:"title" bson:"title"`
	Content   *string             `json:"content" bson:"content"`
	User      primitive.ObjectID  `json:"user" bson:"user"`
	Deleted   *bool               `json:"deleted,omitempty" bson:"deleted,omitempty"`
	CreatedAt primitive.DateTime  `json:"createdAt" bson:"created_at"`
	UpdatedAt primitive.DateTime  `json:"updatedAt" bson:"updated_at"`
}
