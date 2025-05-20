package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TimeEntry struct {
	ID        *primitive.ObjectID `json:"id" bson:"_id"`
	StartDate string              `json:"startDate" bson:"start_date" binding:"required"`
	EndDate   string              `json:"endDate" bson:"end_date" binding:"required"`
	CreatedAt string              `json:"createdAt" bson:"created_at"`
	UpdatedAt string              `json:"updatedAt" bson:"updated_at"`
}
