package payloads

// TaskStartingPayload represents the payload for task starting notifications.
type TaskStartingPayload struct {
	Type  string `json:"type"`
	Title string `json:"title"`
}

// NewTaskStartingPayload creates a new TaskStartingPayload with the given title.
func NewTaskStartingPayload(title string) *TaskStartingPayload {
	return &TaskStartingPayload{
		Type:  "TASK_STARTING",
		Title: title,
	}
}

// GetType returns the type of the payload.
func (p *TaskStartingPayload) GetType() string {
	return p.Type
}

// GetData returns the ready to send data for the payload.
func (p *TaskStartingPayload) GetData() map[string]string {
	return map[string]string{
		"type":  p.Type,
		"title": p.Title,
	}
}
