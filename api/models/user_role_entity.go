package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// UserRoleEntity represents a user role in the system
// @Summary User role entity
// @Description Represents a user role in the system
type UserRoleEntity struct {
	ID        *primitive.ObjectID `json:"id" bson:"_id"`
	Name      string              `json:"name" bson:"name" binding:"required"`
	CreatedAt *primitive.DateTime `json:"createdAt" bson:"created_at"`
	UpdatedAt *primitive.DateTime `json:"updatedAt" bson:"updated_at"`
}
