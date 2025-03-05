package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserEntity struct {
	ID       *primitive.ObjectID `json:"id" bson:"_id"`
	Email    *string             `json:"email" bson:"email" binding:"required"`
	Password *string             `json:"password,omitempty" bson:"password" binding:"required"`
}
