package patchmodels

// PatchChange represents a change in a patch operation.
type PatchChange struct {
	Key   string      `json:"key" bson:"key" binding:"required"`
	Value interface{} `json:"value" bson:"value" binding:"required"`
}
