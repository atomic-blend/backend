package tasks

import "atomic_blend_api/models"

// BulkTaskResponse represents the response for bulk task operations
type BulkTaskResponse struct {
	Updated   []*models.TaskEntity     `json:"updated"`
	Skipped   []*models.TaskEntity     `json:"skipped"`
	Conflicts []*models.ConflictedItem `json:"conflicts"`
}
