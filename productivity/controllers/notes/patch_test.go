package notes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"atomic-blend/backend/productivity/auth"
	patchmodels "atomic-blend/backend/productivity/models/patch_models"
	"atomic-blend/backend/productivity/tests/mocks"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPatch(t *testing.T) {
	_, mockNoteRepo := setupTest()

	t.Run("successful patch update", func(t *testing.T) {
		// Create authenticated user
		noteID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()

		// Create existing note
		existingNote := createTestNote()
		existingNote.UpdatedAt = primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Hour)) // Note was updated 1 hour ago

		// Create patch for update action
		patchDate := primitive.NewDateTimeFromTime(time.Now()) // Patch is newer than note
		patch := patchmodels.Patch{
			ID:       patchID,
			Action:   patchmodels.PatchActionUpdate,
			ItemType: patchmodels.ItemTypeNote,
			ItemID:   &noteID,
			Changes: []patchmodels.PatchChange{
				{
					Key:   "title",
					Value: "Updated Title",
				},
			},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}

		updatedNote := createTestNote()
		title := "Updated Title"
		updatedNote.Title = &title

		mockNoteRepo.On("GetByID", mock.Anything, noteID.Hex()).Return(existingNote, nil)
		mockNoteRepo.On("UpdatePatch", mock.Anything, &patch).Return(updatedNote, nil)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/notes/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: existingNote.User})

		// Call the handler directly
		controller := NewNoteController(mockNoteRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PatchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Success, 1)
		assert.Contains(t, response.Success, patchID.Hex())
		assert.Len(t, response.Errors, 0)
		assert.Len(t, response.Conflicts, 0)
	})

	t.Run("successful patch delete", func(t *testing.T) {
		// Create authenticated user
		noteID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()

		// Create existing note
		existingNote := createTestNote()
		existingNote.UpdatedAt = primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Hour))

		// Create patch for delete action
		patchDate := primitive.NewDateTimeFromTime(time.Now())
		patch := patchmodels.Patch{
			ID:        patchID,
			Action:    patchmodels.PatchActionDelete,
			ItemType:  patchmodels.ItemTypeNote,
			ItemID:    &noteID,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}

		mockNoteRepo.On("GetByID", mock.Anything, noteID.Hex()).Return(existingNote, nil)
		mockNoteRepo.On("Delete", mock.Anything, noteID.Hex()).Return(nil)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/notes/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: existingNote.User})

		// Call the handler directly
		controller := NewNoteController(mockNoteRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PatchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Success, 1)
		assert.Contains(t, response.Success, patchID.Hex())
		assert.Len(t, response.Errors, 0)
		assert.Len(t, response.Conflicts, 0)
	})

	t.Run("successful patch create", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()

		// Create patch for create action
		patchDate := primitive.NewDateTimeFromTime(time.Now())
		title := "New Note"
		content := "New note content"
		newNoteData := map[string]interface{}{
			"title":   title,
			"content": content,
		}

		patch := patchmodels.Patch{
			ID:       patchID,
			Action:   patchmodels.PatchActionCreate,
			ItemType: patchmodels.ItemTypeNote,
			ItemID:   nil, // No ItemID for create action
			Changes: []patchmodels.PatchChange{
				{
					Key:   "data",
					Value: newNoteData,
				},
			},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}

		createdNote := createTestNote()
		createdNote.Title = &title
		createdNote.Content = &content
		createdNote.User = userID

		mockNoteRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.NoteEntity")).Return(createdNote, nil)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/notes/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewNoteController(mockNoteRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PatchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Success, 1)
		assert.Contains(t, response.Success, patchID.Hex())
		assert.Len(t, response.Errors, 0)
		assert.Len(t, response.Conflicts, 0)
	})

	t.Run("unauthorized - no auth user", func(t *testing.T) {
		patchID := primitive.NewObjectID()
		patchDate := primitive.NewDateTimeFromTime(time.Now())
		patch := patchmodels.Patch{
			ID:        patchID,
			Action:    patchmodels.PatchActionUpdate,
			ItemType:  patchmodels.ItemTypeNote,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}
		patchJSON, _ := json.Marshal(patches)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/notes/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context without auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		// Call the handler directly
		controller := NewNoteController(mockNoteRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		userID := primitive.NewObjectID()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/notes/patch", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewNoteController(mockNoteRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("unsupported item type", func(t *testing.T) {
		userID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()
		patchDate := primitive.NewDateTimeFromTime(time.Now())

		patch := patchmodels.Patch{
			ID:        patchID,
			Action:    patchmodels.PatchActionUpdate,
			ItemType:  "unsupported_type",
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}
		patchJSON, _ := json.Marshal(patches)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/notes/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewNoteController(mockNoteRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PatchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Success, 0)
		assert.Len(t, response.Errors, 1)
		assert.Equal(t, "item_type_not_supported", response.Errors[0].ErrorCode)
		assert.Equal(t, patchID.Hex(), response.Errors[0].PatchID)
	})

	t.Run("invalid action", func(t *testing.T) {
		userID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()
		patchDate := primitive.NewDateTimeFromTime(time.Now())

		patch := patchmodels.Patch{
			ID:        patchID,
			Action:    "invalid_action",
			ItemType:  patchmodels.ItemTypeNote,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}
		patchJSON, _ := json.Marshal(patches)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/notes/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewNoteController(mockNoteRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PatchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Success, 0)
		assert.Len(t, response.Errors, 1)
		assert.Equal(t, "invalid_action", response.Errors[0].ErrorCode)
		assert.Equal(t, patchID.Hex(), response.Errors[0].PatchID)
	})

	t.Run("note not found", func(t *testing.T) {
		userID := primitive.NewObjectID()
		noteID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()
		patchDate := primitive.NewDateTimeFromTime(time.Now())

		patch := patchmodels.Patch{
			ID:        patchID,
			Action:    patchmodels.PatchActionUpdate,
			ItemType:  patchmodels.ItemTypeNote,
			ItemID:    &noteID,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}

		mockNoteRepo.On("GetByID", mock.Anything, noteID.Hex()).Return(nil, assert.AnError)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/notes/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewNoteController(mockNoteRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PatchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Success, 0)
		assert.Len(t, response.Errors, 1)
		assert.Equal(t, "note_not_found", response.Errors[0].ErrorCode)
		assert.Equal(t, patchID.Hex(), response.Errors[0].PatchID)
	})

	t.Run("patch conflict - patch is older than note", func(t *testing.T) {
		userID := primitive.NewObjectID()
		noteID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()

		// Create existing note that was updated recently
		existingNote := createTestNote()
		existingNote.UpdatedAt = primitive.NewDateTimeFromTime(time.Now()) // Note was updated now

		// Create patch that is older than the note update
		patchDate := primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Hour)) // Patch is 1 hour old
		patch := patchmodels.Patch{
			ID:        patchID,
			Action:    patchmodels.PatchActionUpdate,
			ItemType:  patchmodels.ItemTypeNote,
			ItemID:    &noteID,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}

		mockNoteRepo.On("GetByID", mock.Anything, noteID.Hex()).Return(existingNote, nil)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/notes/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewNoteController(mockNoteRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PatchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Success, 0)
		assert.Len(t, response.Errors, 0)
		assert.Len(t, response.Conflicts, 1)
		assert.Equal(t, patchID.Hex(), response.Conflicts[0].PatchID)
		assert.NotNil(t, response.Conflicts[0].RemoteObject)
	})

	t.Run("update without item ID", func(t *testing.T) {
		userID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()
		patchDate := primitive.NewDateTimeFromTime(time.Now())

		patch := patchmodels.Patch{
			ID:        patchID,
			Action:    patchmodels.PatchActionUpdate,
			ItemType:  patchmodels.ItemTypeNote,
			ItemID:    nil, // Missing ItemID for update
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}
		patchJSON, _ := json.Marshal(patches)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/notes/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewNoteController(mockNoteRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PatchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Success, 0)
		assert.Len(t, response.Errors, 1)
		assert.Equal(t, "item_id_required", response.Errors[0].ErrorCode)
		assert.Equal(t, patchID.Hex(), response.Errors[0].PatchID)
	})

	t.Run("delete without item ID", func(t *testing.T) {
		userID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()
		patchDate := primitive.NewDateTimeFromTime(time.Now())

		patch := patchmodels.Patch{
			ID:        patchID,
			Action:    patchmodels.PatchActionDelete,
			ItemType:  patchmodels.ItemTypeNote,
			ItemID:    nil, // Missing ItemID for delete
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}
		patchJSON, _ := json.Marshal(patches)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/notes/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewNoteController(mockNoteRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PatchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Success, 0)
		assert.Len(t, response.Errors, 1)
		assert.Equal(t, "item_id_required", response.Errors[0].ErrorCode)
		assert.Equal(t, patchID.Hex(), response.Errors[0].PatchID)
	})

	t.Run("update failed", func(t *testing.T) {
		noteID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()

		existingNote := createTestNote()
		existingNote.UpdatedAt = primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Hour))

		patchDate := primitive.NewDateTimeFromTime(time.Now())
		patch := patchmodels.Patch{
			ID:        patchID,
			Action:    patchmodels.PatchActionUpdate,
			ItemType:  patchmodels.ItemTypeNote,
			ItemID:    &noteID,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}

		mockNoteRepo.On("GetByID", mock.Anything, noteID.Hex()).Return(existingNote, nil)
		mockNoteRepo.On("UpdatePatch", mock.Anything, &patch).Return(nil, assert.AnError)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/notes/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: existingNote.User})

		// Call the handler directly
		controller := NewNoteController(mockNoteRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PatchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Success, 0)
		assert.Len(t, response.Errors, 1)
		assert.Equal(t, "update_failed", response.Errors[0].ErrorCode)
		assert.Equal(t, patchID.Hex(), response.Errors[0].PatchID)
	})

	t.Run("update unauthorized", func(t *testing.T) {
		userID := primitive.NewObjectID()
		noteID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()

		existingNote := createTestNote()
		existingNote.UpdatedAt = primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Hour))

		patchDate := primitive.NewDateTimeFromTime(time.Now())
		patch := patchmodels.Patch{
			ID:        patchID,
			Action:    patchmodels.PatchActionUpdate,
			ItemType:  patchmodels.ItemTypeNote,
			ItemID:    &noteID,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}

		mockNoteRepo.On("GetByID", mock.Anything, noteID.Hex()).Return(existingNote, nil)
		mockNoteRepo.On("UpdatePatch", mock.Anything, &patch).Return(nil, assert.AnError)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/notes/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewNoteController(mockNoteRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PatchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Success, 0)
		assert.Len(t, response.Errors, 1)
		assert.Equal(t, "not_authorized", response.Errors[0].ErrorCode)
		assert.Equal(t, patchID.Hex(), response.Errors[0].PatchID)
	})

	t.Run("delete failed", func(t *testing.T) {
		noteID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()

		existingNote := createTestNote()
		existingNote.UpdatedAt = primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Hour))

		patchDate := primitive.NewDateTimeFromTime(time.Now())
		patch := patchmodels.Patch{
			ID:        patchID,
			Action:    patchmodels.PatchActionDelete,
			ItemType:  patchmodels.ItemTypeNote,
			ItemID:    &noteID,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}

		mockNoteRepo.On("GetByID", mock.Anything, noteID.Hex()).Return(existingNote, nil)
		mockNoteRepo.On("Delete", mock.Anything, noteID.Hex()).Return(assert.AnError)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/notes/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: existingNote.User})

		// Call the handler directly
		controller := NewNoteController(mockNoteRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PatchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Success, 0)
		assert.Len(t, response.Errors, 1)
		assert.Equal(t, "delete_failed", response.Errors[0].ErrorCode)
		assert.Equal(t, patchID.Hex(), response.Errors[0].PatchID)
	})

	t.Run("delete unauthorized", func(t *testing.T) {
		userID := primitive.NewObjectID()
		noteID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()

		existingNote := createTestNote()
		existingNote.UpdatedAt = primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Hour))

		patchDate := primitive.NewDateTimeFromTime(time.Now())
		patch := patchmodels.Patch{
			ID:        patchID,
			Action:    patchmodels.PatchActionDelete,
			ItemType:  patchmodels.ItemTypeNote,
			ItemID:    &noteID,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}

		mockNoteRepo.On("GetByID", mock.Anything, noteID.Hex()).Return(existingNote, nil)
		mockNoteRepo.On("Delete", mock.Anything, noteID.Hex()).Return(assert.AnError)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/notes/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewNoteController(mockNoteRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PatchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Success, 0)
		assert.Len(t, response.Errors, 1)
		assert.Equal(t, "not_authorized", response.Errors[0].ErrorCode)
		assert.Equal(t, patchID.Hex(), response.Errors[0].PatchID)
	})

	t.Run("create with invalid note data", func(t *testing.T) {
		userID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()
		patchDate := primitive.NewDateTimeFromTime(time.Now())

		// Create patch with invalid data that can't be unmarshaled into NoteEntity
		invalidData := make(map[string]interface{})

		patch := patchmodels.Patch{
			ID:       patchID,
			Action:   patchmodels.PatchActionCreate,
			ItemType: patchmodels.ItemTypeNote,
			ItemID:   nil,
			Changes: []patchmodels.PatchChange{
				{
					Key:   "data",
					Value: invalidData,
				},
			},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}
		patchJSON, _ := json.Marshal(patches)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/notes/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewNoteController(mockNoteRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PatchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Success, 1)
		assert.Contains(t, response.Success, patchID.Hex())
		assert.Len(t, response.Errors, 0)
		assert.Len(t, response.Conflicts, 0)
	})

	t.Run("create failed", func(t *testing.T) {
		mockNoteRepo := new(mocks.MockNoteRepository)
		userID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()
		patchDate := primitive.NewDateTimeFromTime(time.Now())

		title := "New Note"
		content := "New note content"
		newNoteData := map[string]interface{}{
			"title":   title,
			"content": content,
		}

		patch := patchmodels.Patch{
			ID:       patchID,
			Action:   patchmodels.PatchActionCreate,
			ItemType: patchmodels.ItemTypeNote,
			ItemID:   nil,
			Changes: []patchmodels.PatchChange{
				{
					Key:   "data",
					Value: newNoteData,
				},
			},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}

		mockNoteRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.NoteEntity")).Return(nil, assert.AnError)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/notes/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewNoteController(mockNoteRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PatchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Success, 0)
		assert.Len(t, response.Errors, 1)
		assert.Equal(t, "create_failed", response.Errors[0].ErrorCode)
		assert.Equal(t, patchID.Hex(), response.Errors[0].PatchID)
	})

	t.Run("mixed operations with successes, errors, and conflicts", func(t *testing.T) {
		// Patch 1: Successful update
		noteID1 := primitive.NewObjectID()
		patchID1 := primitive.NewObjectID()
		existingNote1 := createTestNote()
		existingNote1.UpdatedAt = primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Hour))

		patchDate1 := primitive.NewDateTimeFromTime(time.Now())
		patch1 := patchmodels.Patch{
			ID:        patchID1,
			Action:    patchmodels.PatchActionUpdate,
			ItemType:  patchmodels.ItemTypeNote,
			ItemID:    &noteID1,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate1,
		}

		// Patch 2: Conflict (patch older than note)
		noteID2 := primitive.NewObjectID()
		patchID2 := primitive.NewObjectID()
		existingNote2 := createTestNote()
		existingNote2.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

		patchDate2 := primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Hour))
		patch2 := patchmodels.Patch{
			ID:        patchID2,
			Action:    patchmodels.PatchActionUpdate,
			ItemType:  patchmodels.ItemTypeNote,
			ItemID:    &noteID2,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate2,
		}

		// Patch 3: Error (invalid action)
		patchID3 := primitive.NewObjectID()
		patchDate3 := primitive.NewDateTimeFromTime(time.Now())
		patch3 := patchmodels.Patch{
			ID:        patchID3,
			Action:    "invalid_action",
			ItemType:  patchmodels.ItemTypeNote,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate3,
		}

		patches := []patchmodels.Patch{patch1, patch2, patch3}

		updatedNote1 := createTestNote()

		mockNoteRepo.On("GetByID", mock.Anything, noteID1.Hex()).Return(existingNote1, nil)
		mockNoteRepo.On("GetByID", mock.Anything, noteID2.Hex()).Return(existingNote2, nil)
		mockNoteRepo.On("UpdatePatch", mock.Anything, &patch1).Return(updatedNote1, nil)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/notes/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: existingNote1.User})

		// Call the handler directly
		controller := NewNoteController(mockNoteRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PatchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check success
		assert.Len(t, response.Success, 1)
		assert.Contains(t, response.Success, patchID1.Hex())

		// Check conflicts
		assert.Len(t, response.Conflicts, 1)
		assert.Equal(t, patchID2.Hex(), response.Conflicts[0].PatchID)

		// Check errors
		assert.Len(t, response.Errors, 1)
		assert.Equal(t, patchID3.Hex(), response.Errors[0].PatchID)
		assert.Equal(t, "invalid_action", response.Errors[0].ErrorCode)
	})
}
