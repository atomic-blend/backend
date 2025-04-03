package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	FrequencyDaily     = "daily"
	FrequencyWeekly    = "weekly"
	FrequencyMonthly   = "monthly"
	FrequencyRepeating = "repeatition" // every x days
)

// used in validators.go
var ValidFrequencies = []string{FrequencyDaily, FrequencyWeekly, FrequencyMonthly}

type Habit struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID 	  primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name          *string            `bson:"name" json:"name" binding:"required"`
	Emoji         *string            `bson:"emoji" json:"emoji"`
	Frequency     *string            `bson:"frequency" json:"frequency" binding:"required,validFrequency"`
	NumberOfTimes *int               `bson:"number_of_times" json:"number_of_times"`
	DaysOfWeek    *[]int             `bson:"days_of_week" json:"days_of_week"`
	StartDate     *string            `bson:"start_date" json:"start_date" binding:"required"`
	EndDate       *string            `bson:"end_date" json:"end_date"`
	CreatedAt     *string            `bson:"created_at" json:"created_at"`
	UpdatedAt     *string            `bson:"updated_at" json:"updated_at"`
	Reminders     []string           `bson:"reminders" json:"reminders"`
	Citation      *string            `bson:"citation" json:"citation"`
}
