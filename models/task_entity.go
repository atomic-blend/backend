package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// TaskEntity represents a task
type TaskEntity struct {
	ID          string                `json:"id" bson:"_id"`
	Title       string                `json:"title" bson:"title" binding:"required"`
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

// Equals compares this task with another task to check if they have the same content
// (excluding CreatedAt and UpdatedAt timestamps)
func (t *TaskEntity) Equals(other *TaskEntity) bool {
	if other == nil {
		return false
	}

	// Compare basic fields
	if t.Title != other.Title ||
		t.User != other.User ||
		t.ID != other.ID {
		return false
	}

	// Compare nullable string fields
	if !stringPtrEqual(t.Description, other.Description) {
		return false
	}

	// Compare nullable bool fields
	if !boolPtrEqual(t.Completed, other.Completed) {
		return false
	}

	// Compare nullable int fields
	if !intPtrEqual(t.Priority, other.Priority) {
		return false
	}

	// Compare nullable ObjectID fields
	if !objectIDPtrEqual(t.FolderID, other.FolderID) {
		return false
	}

	// Compare DateTime fields
	if !dateTimePtrEqual(t.StartDate, other.StartDate) {
		return false
	}

	if !dateTimePtrEqual(t.EndDate, other.EndDate) {
		return false
	}

	// Compare Tags
	if !tagsEqual(t.Tags, other.Tags) {
		return false
	}

	// Compare Reminders
	if !remindersEqual(t.Reminders, other.Reminders) {
		return false
	}

	return true
}

// Helper functions for comparing nullable fields
func stringPtrEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func boolPtrEqual(a, b *bool) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func intPtrEqual(a, b *int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func objectIDPtrEqual(a, b *primitive.ObjectID) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func dateTimePtrEqual(a, b *primitive.DateTime) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Time().Equal(b.Time())
}

func tagsEqual(a, b *[]*Tag) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(*a) != len(*b) {
		return false
	}

	// Create maps for easier comparison
	mapA := make(map[string]bool)
	mapB := make(map[string]bool)

	for _, tag := range *a {
		if tag.ID != nil {
			mapA[tag.ID.Hex()] = true
		}
	}

	for _, tag := range *b {
		if tag.ID != nil {
			mapB[tag.ID.Hex()] = true
		}
	}

	// Check if all tags in mapA exist in mapB
	for id := range mapA {
		if !mapB[id] {
			return false
		}
	}

	return true
}

func remindersEqual(a, b []*primitive.DateTime) bool {
	if len(a) != len(b) {
		return false
	}

	// Convert to maps for easier comparison
	mapA := make(map[int64]bool)
	mapB := make(map[int64]bool)

	for _, reminder := range a {
		if reminder != nil {
			mapA[reminder.Time().Unix()] = true
		}
	}

	for _, reminder := range b {
		if reminder != nil {
			mapB[reminder.Time().Unix()] = true
		}
	}

	// Check if all reminders in mapA exist in mapB
	for timestamp := range mapA {
		if !mapB[timestamp] {
			return false
		}
	}

	return true
}
