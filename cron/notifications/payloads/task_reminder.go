package payloads

type TaskReminderPayload struct {
	Type string `json:"type"`
	Title string `json:"title"`
	DueDate string `json:"dueDate"`
}

func NewTaskReminderPayload(title string, dueDate string) *TaskReminderPayload {
	return &TaskReminderPayload{
		Type: "TASK_REMINDER",
		Title: title,
		DueDate: dueDate,
	}
}

func (p *TaskReminderPayload) GetType() string {
	return p.Type
}

func (p *TaskReminderPayload) GetData() map[string]string {
	return map[string]string{
		"type": p.Type,
		"title": p.Title,
		"dueDate": p.DueDate,
	}
}