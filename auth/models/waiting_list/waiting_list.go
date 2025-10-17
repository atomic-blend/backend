package waitinglist

import "go.mongodb.org/mongo-driver/bson/primitive"

// WaitingList represents a waiting list for a user
type WaitingList struct {
	ID        *primitive.ObjectID `bson:"_id" json:"id"`
	Email     string              `bson:"email" json:"email" binding:"required,email"`
	Code      *string             `bson:"code" json:"code"`
	CreatedAt *primitive.DateTime `bson:"created_at" json:"createdAt"`
	UpdatedAt *primitive.DateTime `bson:"updated_at" json:"updatedAt"`
}
