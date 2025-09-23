// Package models is a package that contains the models for the microservice
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Folder represents a folder in the application
type Folder struct {
	ID        *primitive.ObjectID `bson:"_id" json:"id"`
	Name      string              `bson:"name" json:"name" binding:"required"`
	Color     *string             `bson:"color" json:"color"`
	ParentID  *primitive.ObjectID `bson:"parent_id" json:"parentId"`
	UserID    primitive.ObjectID  `bson:"user_id" json:"userId"`
	Emoji     *string             `bson:"emoji" json:"emoji"`
	CreatedAt *primitive.DateTime `bson:"created_at,omitempty" json:"createdAt,omitempty"`
	UpdatedAt *primitive.DateTime `bson:"updated_at,omitempty" json:"updatedAt,omitempty"`
}
