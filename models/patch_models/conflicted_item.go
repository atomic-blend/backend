package patchmodels

// ConflictedItem represents a conflicted item during bulk operations
type ConflictedItem struct {
	Type         string      `json:"type"`
	PatchID      string      `json:"patchId"`
	RemoteObject interface{} `json:"remoteObject"`
	LocalObject  interface{} `json:"localObject"`
}
