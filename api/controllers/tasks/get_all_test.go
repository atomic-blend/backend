package tasks

import (
	"atomic_blend_api/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetAllTasks(t *testing.T) {
	router, mockRepo := setupTest()

	t.Run("successful get all tasks", func(t *testing.T) {
		task := createTestTask()
		task.ID = primitive.NewObjectID().Hex()
		tasks := []*models.TaskEntity{task}

		// Updated the mock to include both context and userID parameters
		mockRepo.On("GetAll", mock.Anything, mock.AnythingOfType("*primitive.ObjectID")).Return(tasks, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []*models.TaskEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, tasks[0].Title, response[0].Title)
		assert.NotNil(t, response[0].StartDate)
		assert.NotNil(t, response[0].EndDate)
	})
}
