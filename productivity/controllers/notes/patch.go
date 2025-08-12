package notes

import (
	"context"
	"encoding/json"
	"net/http"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/atomic-blend/backend/productivity/models"
	patchmodels "github.com/atomic-blend/backend/productivity/models/patch_models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PatchResponse represents the response structure for the patch operation
type PatchResponse struct {
	Success   []string                     `json:"success"`
	Errors    []patchmodels.PatchError     `json:"errors"`
	Conflicts []patchmodels.ConflictedItem `json:"conflicts"`
	Date      time.Time                    `json:"date"`
}

// Patch handles the patching of notes
func (c *NoteController) Patch(ctx *gin.Context) {
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
		//check if patch type is note
		if patch.ItemType != patchmodels.ItemTypeNote {
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
			// get the note by ID
			note, err := c.noteRepo.GetByID(ct, patch.ItemID.Hex())
			if err != nil {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(), ErrorCode: "note_not_found"})
				continue
			}

			// check that the patch is dated after the last update of the note
			if patch.Action != patchmodels.PatchActionCreate && patch.PatchDate.Time().Before(note.UpdatedAt.Time()) {
				if patch.Force == nil || !*patch.Force {
					conflicts = append(conflicts, patchmodels.ConflictedItem{Type: "note", PatchID: patch.ID.Hex(), RemoteObject: note})
					continue
				}
			}
		}

		switch patch.Action {
		case patchmodels.PatchActionUpdate:
			// get the original note and check ownership
			if patch.ItemID == nil {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(), ErrorCode: "item_id_required"})
				continue
			}
			note, err := c.noteRepo.GetByID(ct, patch.ItemID.Hex())
			if err != nil {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(), ErrorCode: "note_not_found"})
				continue
			}

			if note.User != authUser.UserID {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(), ErrorCode: "not_authorized"})
				continue
			}

			// apply the changes to the note
			_, err = c.noteRepo.UpdatePatch(ct, &patch)
			if err != nil {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(), ErrorCode: "update_failed"})
			} else {
				success = append(success, patch.ID.Hex())
			}
			continue
		case patchmodels.PatchActionDelete:
			// get the original note and check ownership
			if patch.ItemID == nil {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(), ErrorCode: "item_id_required"})
				continue
			}
			note, err := c.noteRepo.GetByID(ct, patch.ItemID.Hex())
			if err != nil {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(), ErrorCode: "note_not_found"})
				continue
			}

			if note.User != authUser.UserID {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(), ErrorCode: "not_authorized"})
				continue
			}

			// delete the note
			err = c.noteRepo.Delete(ct, patch.ItemID.Hex())
			if err != nil {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(),
					ErrorCode: "delete_failed"})
			} else {
				success = append(success, patch.ID.Hex())
			}
			continue
		case patchmodels.PatchActionCreate:
			// the content of the note is under the data key of the first change
			newNote := &models.NoteEntity{}

			// Convert map to json string
			jsonStr, err := json.Marshal(patch.Changes[0].Value)
			if err != nil {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(),
					ErrorCode: "invalid_note_data"})
				continue
			}
			if err := json.Unmarshal(jsonStr, newNote); err != nil {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(),
					ErrorCode: "invalid_note_data"})
				continue
			}

			validate := validator.New()
			if err := validate.Struct(newNote); err != nil {
				errors = append(errors, patchmodels.PatchError{PatchID: patch.ID.Hex(),
					ErrorCode: "invalid_note_data"})
				continue
			}

			newNote.User = authUser.UserID
			newNote.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
			newNote.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

			_, err = c.noteRepo.Create(ct, newNote)
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
