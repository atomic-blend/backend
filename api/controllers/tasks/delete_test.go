package tasks

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
func TestDeleteTask(t *testing.T) {
	router, mockRepo := setupTest()

	t.Run("successful delete task", func(t *testing.T) {
		taskID := primitive.NewObjectID().Hex()
		existingTask := createTestTask()
		existingTask.ID = taskID

		mockRepo.On("GetByID", mock.Anything, taskID).Return(existingTask, nil)
		mockRepo.On("Delete", mock.Anything, taskID).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tasks/"+taskID, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}