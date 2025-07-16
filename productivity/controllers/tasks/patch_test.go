package tasks

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
	_, mockTaskRepo, mockTagRepo := setupTest()

	t.Run("successful patch update", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()

		// Create existing task
		existingTask := createTestTask()
		existingTask.ID = taskID.Hex()
		existingTask.User = userID
		existingTask.UpdatedAt = primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Hour)) // Task was updated 1 hour ago

		// Create patch for update action
		patchDate := primitive.NewDateTimeFromTime(time.Now()) // Patch is newer than task
		patch := patchmodels.Patch{
			ID:       patchID,
			Action:   patchmodels.PatchActionUpdate,
			ItemType: patchmodels.ItemTypeTask,
			ItemID:   &taskID,
			Changes: []patchmodels.PatchChange{
				{
					Key:   "title",
					Value: "Updated Title",
				},
			},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}

		updatedTask := createTestTask()
		updatedTask.ID = taskID.Hex()
		updatedTask.Title = "Updated Title"

		mockTaskRepo.On("GetByID", mock.Anything, taskID.Hex()).Return(existingTask, nil)
		mockTaskRepo.On("UpdatePatch", mock.Anything, &patch).Return(updatedTask, nil)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
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
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()

		// Create existing task
		existingTask := createTestTask()
		existingTask.ID = taskID.Hex()
		existingTask.User = userID
		existingTask.UpdatedAt = primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Hour))

		// Create patch for delete action
		patchDate := primitive.NewDateTimeFromTime(time.Now())
		patch := patchmodels.Patch{
			ID:        patchID,
			Action:    patchmodels.PatchActionDelete,
			ItemType:  patchmodels.ItemTypeTask,
			ItemID:    &taskID,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}

		mockTaskRepo.On("GetByID", mock.Anything, taskID.Hex()).Return(existingTask, nil)
		mockTaskRepo.On("Delete", mock.Anything, taskID.Hex()).Return(nil)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
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
		newTaskData := map[string]interface{}{
			"title":       "New Task",
			"description": "New task description",
			"completed":   false,
		}

		patch := patchmodels.Patch{
			ID:       patchID,
			Action:   patchmodels.PatchActionCreate,
			ItemType: patchmodels.ItemTypeTask,
			ItemID:   nil, // No ItemID for create action
			Changes: []patchmodels.PatchChange{
				{
					Key:   "data",
					Value: newTaskData,
				},
			},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}

		createdTask := createTestTask()
		createdTask.Title = "New Task"
		createdTask.User = userID

		mockTaskRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.TaskEntity")).Return(createdTask, nil)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
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
			ItemType:  patchmodels.ItemTypeTask,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}
		patchJSON, _ := json.Marshal(patches)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context without auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		userID := primitive.NewObjectID()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks/patch", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
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
		req, _ := http.NewRequest("POST", "/tasks/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
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
			ItemType:  patchmodels.ItemTypeTask,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}
		patchJSON, _ := json.Marshal(patches)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
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

	t.Run("task not found", func(t *testing.T) {
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()
		patchDate := primitive.NewDateTimeFromTime(time.Now())

		patch := patchmodels.Patch{
			ID:        patchID,
			Action:    patchmodels.PatchActionUpdate,
			ItemType:  patchmodels.ItemTypeTask,
			ItemID:    &taskID,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}

		mockTaskRepo.On("GetByID", mock.Anything, taskID.Hex()).Return(nil, assert.AnError)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PatchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Success, 0)
		assert.Len(t, response.Errors, 1)
		assert.Equal(t, "task_not_found", response.Errors[0].ErrorCode)
		assert.Equal(t, patchID.Hex(), response.Errors[0].PatchID)
	})

	t.Run("patch conflict - patch is older than task", func(t *testing.T) {
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()

		// Create existing task that was updated recently
		existingTask := createTestTask()
		existingTask.ID = taskID.Hex()
		existingTask.User = userID
		existingTask.UpdatedAt = primitive.NewDateTimeFromTime(time.Now()) // Task was updated now

		// Create patch that is older than the task update
		patchDate := primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Hour)) // Patch is 1 hour old
		patch := patchmodels.Patch{
			ID:        patchID,
			Action:    patchmodels.PatchActionUpdate,
			ItemType:  patchmodels.ItemTypeTask,
			ItemID:    &taskID,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}

		mockTaskRepo.On("GetByID", mock.Anything, taskID.Hex()).Return(existingTask, nil)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
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
			ItemType:  patchmodels.ItemTypeTask,
			ItemID:    nil, // Missing ItemID for update
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}
		patchJSON, _ := json.Marshal(patches)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
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
			ItemType:  patchmodels.ItemTypeTask,
			ItemID:    nil, // Missing ItemID for delete
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}
		patchJSON, _ := json.Marshal(patches)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
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
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()

		existingTask := createTestTask()
		existingTask.ID = taskID.Hex()
		existingTask.User = userID
		existingTask.UpdatedAt = primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Hour))

		patchDate := primitive.NewDateTimeFromTime(time.Now())
		patch := patchmodels.Patch{
			ID:        patchID,
			Action:    patchmodels.PatchActionUpdate,
			ItemType:  patchmodels.ItemTypeTask,
			ItemID:    &taskID,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}

		mockTaskRepo.On("GetByID", mock.Anything, taskID.Hex()).Return(existingTask, nil)
		mockTaskRepo.On("UpdatePatch", mock.Anything, &patch).Return(nil, assert.AnError)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
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

	t.Run("delete failed", func(t *testing.T) {
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()

		existingTask := createTestTask()
		existingTask.ID = taskID.Hex()
		existingTask.User = userID
		existingTask.UpdatedAt = primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Hour))

		patchDate := primitive.NewDateTimeFromTime(time.Now())
		patch := patchmodels.Patch{
			ID:        patchID,
			Action:    patchmodels.PatchActionDelete,
			ItemType:  patchmodels.ItemTypeTask,
			ItemID:    &taskID,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}

		mockTaskRepo.On("GetByID", mock.Anything, taskID.Hex()).Return(existingTask, nil)
		mockTaskRepo.On("Delete", mock.Anything, taskID.Hex()).Return(assert.AnError)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
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

	t.Run("create with invalid task data", func(t *testing.T) {
		userID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()
		patchDate := primitive.NewDateTimeFromTime(time.Now())

		// Create patch with invalid data that can't be unmarshaled into TaskEntity
		invalidData := make(map[string]interface {
		})

		patch := patchmodels.Patch{
			ID:       patchID,
			Action:   patchmodels.PatchActionCreate,
			ItemType: patchmodels.ItemTypeTask,
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
		req, _ := http.NewRequest("POST", "/tasks/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.Patch(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PatchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Success, 0)
		assert.Len(t, response.Errors, 1)
		assert.Equal(t, "invalid_task_data", response.Errors[0].ErrorCode)
		assert.Equal(t, patchID.Hex(), response.Errors[0].PatchID)
	})

	t.Run("create failed", func(t *testing.T) {
		mockTaskRepo := new(mocks.MockTaskRepository)
		mockTagRepo := new(mocks.MockTagRepository)
		userID := primitive.NewObjectID()
		patchID := primitive.NewObjectID()
		patchDate := primitive.NewDateTimeFromTime(time.Now())

		newTaskData := map[string]interface{}{
			"title":       "New Task",
			"description": "New task description",
			"completed":   false,
		}

		patch := patchmodels.Patch{
			ID:       patchID,
			Action:   patchmodels.PatchActionCreate,
			ItemType: patchmodels.ItemTypeTask,
			ItemID:   nil,
			Changes: []patchmodels.PatchChange{
				{
					Key:   "data",
					Value: newTaskData,
				},
			},
			PatchDate: &patchDate,
		}

		patches := []patchmodels.Patch{patch}

		mockTaskRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.TaskEntity")).Return(nil, assert.AnError)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
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
		userID := primitive.NewObjectID()

		// Patch 1: Successful update
		taskID1 := primitive.NewObjectID()
		patchID1 := primitive.NewObjectID()
		existingTask1 := createTestTask()
		existingTask1.ID = taskID1.Hex()
		existingTask1.User = userID
		existingTask1.UpdatedAt = primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Hour))

		patchDate1 := primitive.NewDateTimeFromTime(time.Now())
		patch1 := patchmodels.Patch{
			ID:        patchID1,
			Action:    patchmodels.PatchActionUpdate,
			ItemType:  patchmodels.ItemTypeTask,
			ItemID:    &taskID1,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate1,
		}

		// Patch 2: Conflict (patch older than task)
		taskID2 := primitive.NewObjectID()
		patchID2 := primitive.NewObjectID()
		existingTask2 := createTestTask()
		existingTask2.ID = taskID2.Hex()
		existingTask2.User = userID
		existingTask2.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

		patchDate2 := primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Hour))
		patch2 := patchmodels.Patch{
			ID:        patchID2,
			Action:    patchmodels.PatchActionUpdate,
			ItemType:  patchmodels.ItemTypeTask,
			ItemID:    &taskID2,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate2,
		}

		// Patch 3: Error (invalid action)
		patchID3 := primitive.NewObjectID()
		patchDate3 := primitive.NewDateTimeFromTime(time.Now())
		patch3 := patchmodels.Patch{
			ID:        patchID3,
			Action:    "invalid_action",
			ItemType:  patchmodels.ItemTypeTask,
			Changes:   []patchmodels.PatchChange{},
			PatchDate: &patchDate3,
		}

		patches := []patchmodels.Patch{patch1, patch2, patch3}

		updatedTask1 := createTestTask()
		updatedTask1.ID = taskID1.Hex()

		mockTaskRepo.On("GetByID", mock.Anything, taskID1.Hex()).Return(existingTask1, nil)
		mockTaskRepo.On("GetByID", mock.Anything, taskID2.Hex()).Return(existingTask2, nil)
		mockTaskRepo.On("UpdatePatch", mock.Anything, &patch1).Return(updatedTask1, nil)

		patchJSON, _ := json.Marshal(patches)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks/patch", bytes.NewBuffer(patchJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
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
