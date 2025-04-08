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
	ID            primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID  `bson:"user_id" json:"userId"`
	Name          *string             `bson:"name" json:"name" binding:"required"`
	Emoji         *string             `bson:"emoji" json:"emoji"`
	Frequency     *string             `bson:"frequency" json:"frequency" binding:"required,validFrequency"`
	NumberOfTimes *int                `bson:"number_of_times" json:"numberOfTimes"`
	Duration      *int                `bson:"duration" json:"duration"`
	DaysOfWeek    *[]int              `bson:"days_of_week" json:"daysOfWeek"`
	StartDate     *primitive.DateTime `bson:"start_date" json:"startDate" binding:"required"`
	EndDate       *primitive.DateTime `bson:"end_date" json:"endDate"`
	CreatedAt     *string             `bson:"created_at" json:"createdAt"`
	UpdatedAt     *string             `bson:"updated_at" json:"updatedAt"`
	Reminders     []string            `bson:"reminders" json:"reminders"`
	Citation      *string             `bson:"citation" json:"citation"`
	Entries       []HabitEntry        `json:"entries"`
}
