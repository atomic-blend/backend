package models

// ConflictedItem represents a conflicted item during bulk operations
type ConflictedItem struct {
	Type string      `json:"type"`
	Old  interface{} `json:"old"`
	New  interface{} `json:"new"`
}
