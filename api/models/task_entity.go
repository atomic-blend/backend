package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TaskEntity struct {
	ID		  string `json:"id" bson:"_id"`
	Title    string `json:"title" bson:"title" binding:"required"`
	User primitive.ObjectID `json:"user" bson:"user"`
	Description *string `json:"description" bson:"description"`
	StartDate *primitive.DateTime `json:"startDate" bson:"start_date" binding:"required"`
	EndDate *primitive.DateTime `json:"endDate" bson:"end_date" binding:"required"`
	Completed *bool `json:"completed" bson:"completed"`
	CreatedAt string `json:"createdAt" bson:"created_at"`
	UpdatedAt string `json:"updatedAt" bson:"updated_at"`
}