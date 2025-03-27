package payloads

type TaskStartingPayload struct {
	Type string `json:"type"`
	Title string `json:"title"`
}

func NewTaskStartingPayload(title string) *TaskStartingPayload {
	return &TaskStartingPayload{
		Type: "TASK_STARTING",
		Title: title,
	}
}

func (p *TaskStartingPayload) GetType() string {
	return p.Type
}

func (p *TaskStartingPayload) GetData() map[string]string {
	return map[string]string{
		"type": p.Type,
		"title": p.Title,
	}
}