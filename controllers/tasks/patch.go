package tasks

import (
	patchmodels "atomic_blend_api/models/patch_models"
	"context"

	"github.com/gin-gonic/gin"
)

type PatchResponse struct {
	Success   []string                     `json:"success"`
	Errors    []patchmodels.PatchError     `json:"errors"`
	Conflicts []patchmodels.ConflictedItem `json:"conflicts"`
}

func (c *TaskController) Patch(ctx *gin.Context) {
	var ct = context.TODO()
	// parse the request body into a Patch struct
	var patchs []patchmodels.Patch
	if err := ctx.ShouldBindJSON(&patchs); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var success = make([]string, 0)
	var errors = make([]patchmodels.PatchError, 0)
	var conflicts = make([]patchmodels.ConflictedItem, 0)

	for _, patch := range patchs {
		//check if patch type is task
		if patch.ItemType != patchmodels.ItemTypeTask {
			errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(), ErrorCode: "item_type_not_supported"})
			continue
		}

		// check if action is contained in validPatchActions
		if !containString(patchmodels.ValidPatchActions, patch.Action) {
			errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(), ErrorCode: "invalid_action"})
			continue
		}

		// get the task by ID
		task, err := c.taskRepo.GetByID(ct, patch.ItemID.Hex())
		if err != nil {
			errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(), ErrorCode: "task_not_found"})
			continue
		}

		// check that the patch is dated after the last update of the task
		if patch.UpdatedAt.Time().Before(task.UpdatedAt.Time()) {
			conflicts = append(conflicts, patchmodels.ConflictedItem{PatchID: patch.ID.Hex(), RemoteObject: task})
			continue
		}

		switch patch.Action {
		case patchmodels.PatchActionUpdate:
			_, err := c.taskRepo.UpdatePatch(ct, &patch)
			if err != nil {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(), ErrorCode: "update_failed"})
			} else {
				success = append(success, patch.ID.Hex())
			}
			continue
		case patchmodels.PatchActionDelete:
			//TODO:
			errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(), ErrorCode: "delete_not_supported"})
			continue
		case patchmodels.PatchActionCreate:
			//TODO:
			errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(), ErrorCode: "create_not_supported"})
			continue
		}
	}

	ctx.JSON(200, PatchResponse{
		Success:   success,
		Errors:    errors,
		Conflicts: conflicts,
	})
}

func containString(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
