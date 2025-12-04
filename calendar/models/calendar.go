package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Calendar represents an RFC 5545 VCALENDAR object.
type Calendar struct {
	ID        *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    *primitive.ObjectID `json:"userId" bson:"user_id,omitempty"`
	ProdID    string              `json:"prodId,omitempty" bson:"prod_id,omitempty"` // PRODID (required)
	Version   string              `json:"version" bson:"version"`                    // VERSION (must be "2.0")
	Method    string              `json:"method,omitempty" bson:"method,omitempty"`  // METHOD (optional, used for iTIP: REQUEST/PUBLISH/etc.)
	Events    []Event             `json:"events,omitempty" bson:"events,omitempty"`  // VEVENT components
	Timezones []VTimezone         `json:"timezones,omitempty" bson:"timezones,omitempty"`
	CreatedAt *primitive.DateTime `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *primitive.DateTime `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}
