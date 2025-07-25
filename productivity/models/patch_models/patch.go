package patchmodels

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Patch represents a patch operation on an item in the system.
type Patch struct {
	ID        primitive.ObjectID  `json:"id" bson:"_id" binding:"required"`
	Action    string              `json:"action" bson:"action" binding:"required"`
	ItemType  string              `json:"itemType" bson:"item_type" binding:"required"`
	ItemID    *primitive.ObjectID `json:"itemId" bson:"item_id"`
	Changes   []PatchChange       `json:"changes" bson:"changes" binding:"required"`
	PatchDate *primitive.DateTime `json:"patchDate" bson:"patch_date" binding:"required"`
	Force     *bool               `json:"force,omitempty" bson:"force,omitempty"`
	CreatedAt *primitive.DateTime `json:"createdAt" bson:"created_at"`
	UpdatedAt *primitive.DateTime `json:"updatedAt" bson:"updated_at"`
}
