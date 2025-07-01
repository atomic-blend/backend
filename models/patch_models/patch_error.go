package patchmodels

type PatchError struct {
	PatchID   string `json:"patchId"`
	ErrorCode string `json:"errorCode"`
}
