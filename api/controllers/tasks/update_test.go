package tasks

import (
	"atomic_blend_api/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func TestUpdateTask(t *testing.T) {
	router, mockRepo := setupTest()

	t.Run("successful update task", func(t *testing.T) {
		taskID := primitive.NewObjectID().Hex()
		existingTask := createTestTask()
		existingTask.ID = taskID

		updatedTask := createTestTask()
		updatedTask.ID = taskID
		updatedTask.Title = "Updated Task"

		mockRepo.On("GetByID", mock.Anything, taskID).Return(existingTask, nil)
		mockRepo.On("Update", mock.Anything, taskID, mock.AnythingOfType("*models.TaskEntity")).Return(updatedTask, nil)

		taskJSON, _ := json.Marshal(updatedTask)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/"+taskID, bytes.NewBuffer(taskJSON))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.TaskEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, updatedTask.Title, response.Title)
		assert.Equal(t, updatedTask.StartDate, response.StartDate)
		assert.Equal(t, updatedTask.EndDate, response.EndDate)
	})
}