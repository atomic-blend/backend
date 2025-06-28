package models

// ConflictedItem represents a conflicted item during bulk operations
type ConflictedItem struct {
	Type    string      `json:"type"`
	OldItem interface{} `json:"old_item"`
	NewItem interface{} `json:"new_item"`
}
