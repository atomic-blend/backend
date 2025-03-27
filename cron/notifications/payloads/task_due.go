package payloads

type TaskDuePayload struct {
	Type string `json:"type"`
	Title string `json:"title"`
}

func NewTaskDuePayload(title string) *TaskDuePayload {
	return &TaskDuePayload{
		Type: "TASK_DUE",
		Title: title,
	}
}

func (p *TaskDuePayload) GetType() string {
	return p.Type
}

func (p *TaskDuePayload) GetData() map[string]string {
	return map[string]string{
		"type": p.Type,
		"title": p.Title,
	}
}