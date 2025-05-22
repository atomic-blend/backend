package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Folder struct {
	ID        *primitive.ObjectID `bson:"_id" json:"id"`
	Name      string              `bson:"name" json:"name" binding:"required"`
	Color     *string             `bson:"color" json:"color"`
	ParentID  *primitive.ObjectID `bson:"parent_id" json:"parent_id"`
	UserID    primitive.ObjectID  `bson:"user_id" json:"user_id"`
	Emoji     *string             `bson:"emoji" json:"emoji"`
	CreatedAt *primitive.DateTime `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt *primitive.DateTime `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
