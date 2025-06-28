package tasks

import "atomic_blend_api/models"

// BulkTaskRequest represents the request payload for bulk task operations
type BulkTaskRequest struct {
	Tasks []*models.TaskEntity `json:"tasks" binding:"required"`
}

// BulkTaskResponse represents the response for bulk task operations
type BulkTaskResponse struct {
	Updated   []*models.TaskEntity     `json:"updated"`
	Conflicts []*models.ConflictedItem `json:"conflicts"`
}
