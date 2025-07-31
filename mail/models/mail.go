package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Mail struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	From      primitive.ObjectID `bson:"from" json:"from"`
	To        primitive.ObjectID `bson:"to" json:"to"`
	Subject   string `bson:"subject" json:"subject"`
	Content   string `bson:"content" json:"content"`
	CreatedAt *primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt *primitive.DateTime `bson:"updated_at" json:"updated_at"`
}