package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// HabitEntry represents a single entry for a habit.
type HabitEntry struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	HabitID   primitive.ObjectID `json:"habitId" bson:"habit_id" binding:"required"`
	UserID    primitive.ObjectID `json:"userId" bson:"user_id"`
	EntryDate primitive.DateTime `json:"entryDate" bson:"entry_date"`
	CreatedAt string             `json:"createdAt" bson:"created_at"`
	UpdatedAt string             `json:"updatedAt" bson:"updated_at"`
}
