package patchmodels

const (
	// ItemTypeTask is the type for tasks
	ItemTypeTask = "task"
	// ItemTypeNote is the type for notes
	ItemTypeNote = "note"
)

// ValidItemTypes contains the valid item types for patch operations
var ValidItemTypes = []string{ItemTypeTask, ItemTypeNote}
