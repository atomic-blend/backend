package payloads

// HabitReminderPayload represents the payload for task starting notifications.
type HabitReminderPayload struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Citation string `json:"citation"`
	Emoji *string `json:"emoji"`
}

// NewHabitReminderPayload creates a new HabitReminderPayload with the given title.
func NewHabitReminderPayload(title string, citation string, emoji *string) *HabitReminderPayload {
	return &HabitReminderPayload{
		Type:  "HABIT_REMINDER",
		Title: title,
		Citation: citation,
		Emoji: emoji,
	}
}

// GetType returns the type of the payload.
func (p *HabitReminderPayload) GetType() string {
	return p.Type
}

// GetData returns the ready to send data for the payload.
func (p *HabitReminderPayload) GetData() map[string]string {
	return map[string]string{
		"type":  p.Type,
		"title": p.Title,
		"citation": p.Citation,
	}
}
