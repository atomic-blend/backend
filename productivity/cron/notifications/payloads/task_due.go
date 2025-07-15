package payloads

// TaskDuePayload represents the payload for task due notifications.
type TaskDuePayload struct {
	Type string `json:"type"`
	Title string `json:"title"`
}

// NewTaskDuePayload creates a new TaskDuePayload with the given title.
func NewTaskDuePayload(title string) *TaskDuePayload {
	return &TaskDuePayload{
		Type: "TASK_DUE",
		Title: title,
	}
}

// GetType returns the type of the payload.
func (p *TaskDuePayload) GetType() string {
	return p.Type
}


// GetData returns the ready to send data for the payload.
func (p *TaskDuePayload) GetData() map[string]string {
	return map[string]string{
		"type": p.Type,
		"title": p.Title,
	}
}