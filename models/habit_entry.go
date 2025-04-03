package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type HabitEntry struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	HabitID   primitive.ObjectID `json:"habit_id" bson:"habit_id" binding:"required"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	EntryDate string             `json:"entry_date" bson:"entry_date"`
	CreatedAt string             `json:"created_at" bson:"created_at"`
	UpdatedAt string             `json:"updated_at" bson:"updated_at"`
}