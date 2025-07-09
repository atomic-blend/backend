package tasks

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
	patchmodels "atomic_blend_api/models/patch_models"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PatchResponse struct {
	Success   []string                     `json:"success"`
	Errors    []patchmodels.PatchError     `json:"errors"`
	Conflicts []patchmodels.ConflictedItem `json:"conflicts"`
	Date      time.Time                    `json:"date"`
}

func (c *TaskController) Patch(ctx *gin.Context) {
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

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

		if patch.Action != patchmodels.PatchActionCreate {
			if patch.ItemID == nil {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(), ErrorCode: "item_id_required"})
				continue
			}
			// get the task by ID
			task, err := c.taskRepo.GetByID(ct, patch.ItemID.Hex())
			if err != nil {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(), ErrorCode: "task_not_found"})
				continue
			}

			// check that the patch is dated after the last update of the task
			if patch.Action != patchmodels.PatchActionCreate && patch.PatchDate.Time().Before(task.UpdatedAt.Time()) {
				if patch.Force == nil || !*patch.Force {
					conflicts = append(conflicts, patchmodels.ConflictedItem{Type: "task", PatchID: patch.ID.Hex(), RemoteObject: task})
					continue
				}
			}
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
			err := c.taskRepo.Delete(ct, patch.ItemID.Hex())
			if err != nil {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(),
					ErrorCode: "delete_failed"})
			} else {
				success = append(success, patch.ID.Hex())
			}
			continue
		case patchmodels.PatchActionCreate:
			// the content of the task is under the data key of the first change
			newTask := &models.TaskEntity{}

			// Convert map to json string
			jsonStr, err := json.Marshal(patch.Changes[0].Value)
			if err != nil {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(),
					ErrorCode: "invalid_task_data"})
				continue
			}
			if err := json.Unmarshal(jsonStr, newTask); err != nil {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(),
					ErrorCode: "invalid_task_data"})
				continue
			}

			validate := validator.New()
			if err := validate.Struct(newTask); err != nil {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(),
					ErrorCode: "invalid_task_data"})
				continue
			}

			newTask.User = authUser.UserID
			newTask.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
			newTask.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

			_, err = c.taskRepo.Create(ct, newTask)
			if err != nil {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(),
					ErrorCode: "create_failed"})
			} else {
				success = append(success, patch.ID.Hex())
			}
			continue
		}
	}

	ctx.JSON(200, PatchResponse{
		Success:   success,
		Errors:    errors,
		Conflicts: conflicts,
		Date:      time.Now(),
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
