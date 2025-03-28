package payloads

// TaskReminderPayload represents the payload for task reminder notifications.
type TaskReminderPayload struct {
	Type string `json:"type"`
	Title string `json:"title"`
	DueDate string `json:"dueDate"`
}

// NewTaskReminderPayload creates a new TaskReminderPayload with the given title and due date.
func NewTaskReminderPayload(title string, dueDate string) *TaskReminderPayload {
	return &TaskReminderPayload{
		Type: "TASK_REMINDER",
		Title: title,
		DueDate: dueDate,
	}
}

// GetType returns the type of the payload.
func (p *TaskReminderPayload) GetType() string {
	return p.Type
}

// GetData returns the ready to send data for the payload.
func (p *TaskReminderPayload) GetData() map[string]string {
	return map[string]string{
		"type": p.Type,
		"title": p.Title,
		"dueDate": p.DueDate,
	}
}