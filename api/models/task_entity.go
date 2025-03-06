package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TaskEntity struct {
	ID		  string `json:"id" bson:"_id"`
	Title    string `json:"title" bson:"title" binding:"required"`
	User primitive.ObjectID `json:"user" bson:"user"`
	Description *string `json:"description" bson:"description"`
	StartDate *primitive.DateTime `json:"start_date" bson:"start_date" binding:"required"`
	EndDate *primitive.DateTime `json:"end_date" bson:"end_date" binding:"required"`
	Completed *bool `json:"completed" bson:"completed"`
	CreatedAt string `json:"created_at" bson:"created_at"`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
}