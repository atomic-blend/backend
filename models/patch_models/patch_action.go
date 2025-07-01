package patchmodels

const (
	// PatchActionCreate represents a create action in patch operations
	PatchActionCreate = "create"
	// PatchActionUpdate represents an update action in patch operations
	PatchActionUpdate = "update"
	// PatchActionDelete represents a delete action in patch operations
	PatchActionDelete = "delete"
)

var ValidPatchActions = []string{PatchActionCreate, PatchActionUpdate, PatchActionDelete}
