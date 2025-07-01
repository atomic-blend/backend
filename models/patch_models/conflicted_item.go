package patchmodels


// ConflictedItem represents a conflicted item during bulk operations
type ConflictedItem struct {
	PatchID      string      `json:"patchId"`
	RemoteObject interface{} `json:"remoteObject"`
}
