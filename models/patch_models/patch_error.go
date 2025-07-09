package patchmodels

// PatchError represents an error that occurred during a patch operation.
type PatchError struct {
	PatchID   string `json:"patchId"`
	ErrorCode string `json:"errorCode"`
}
