package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// UserEntity represents a user in the system
// @Summary User entity
// @Description Represents a user in the system
type UserEntity struct {
	ID       *primitive.ObjectID `json:"id" bson:"_id"`
	Email    *string             `json:"email" bson:"email" binding:"required"`
	Password *string             `json:"password,omitempty" bson:"password" binding:"required"`
	Roles    []*UserRoleEntity   `json:"roles" bson:"roles"`
}
