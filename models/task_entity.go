package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// TaskEntity represents a task
type TaskEntity struct {
	ID          string                `json:"id" bson:"_id"`
	Title       string                `json:"title" bson:"title" binding:"required" validate:"required"`
	User        primitive.ObjectID    `json:"user" bson:"user"`
	Description *string               `json:"description" bson:"description"`
	StartDate   *primitive.DateTime   `json:"startDate" bson:"start_date"`
	EndDate     *primitive.DateTime   `json:"endDate,omitempty" bson:"end_date"`
	Reminders   []*primitive.DateTime `json:"reminders,omitempty" bson:"reminders"`
	Completed   *bool                 `json:"completed" bson:"completed"`
	Tags        *[]*Tag               `json:"tags" bson:"tags"`
	Priority    *int                  `json:"priority" bson:"priority"`
	FolderID    *primitive.ObjectID   `json:"folderId" bson:"folder_id"`
	// TimeEntries []*TimeEntry          `json:"timeEntries" bson:"time_entries"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"created_at"`
	UpdatedAt primitive.DateTime `json:"updatedAt" bson:"updated_at"`
}
