package tasks

import (
	"atomic_blend_api/models"
	"atomic_blend_api/tests/mocks"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createTestTask() *models.TaskEntity {
	desc := "Test Description"
	completed := false
	now := primitive.NewDateTimeFromTime(time.Now())
	end := primitive.NewDateTimeFromTime(time.Now().Add(24 * time.Hour))

	return &models.TaskEntity{
		Title:       "Test Task",
		Description: &desc,
		Completed:   &completed,
		StartDate:   &now,
		EndDate:     &end,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
}

func setupTest() (*gin.Engine, *mocks.MockTaskRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockRepo := new(mocks.MockTaskRepository)
	controller := NewTaskController(mockRepo)
	controller.SetupRoutes(router.Group("/api"))
	return router, mockRepo
}

func TestGetAllTasks(t *testing.T) {
	router, mockRepo := setupTest()

	t.Run("successful get all tasks", func(t *testing.T) {
		task := createTestTask()
		task.ID = primitive.NewObjectID().Hex()
		tasks := []*models.TaskEntity{task}

		mockRepo.On("GetAll", mock.Anything, (*primitive.ObjectID)(nil)).Return(tasks, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/tasks", nil)
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

func TestGetTaskByID(t *testing.T) {
	router, mockRepo := setupTest()

	t.Run("successful get task by id", func(t *testing.T) {
		taskID := primitive.NewObjectID().Hex()
		task := createTestTask()
		task.ID = taskID

		mockRepo.On("GetByID", mock.Anything, taskID).Return(task, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/tasks/"+taskID, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.TaskEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, task.Title, response.Title)
		assert.Equal(t, task.StartDate, response.StartDate)
		assert.Equal(t, task.EndDate, response.EndDate)
	})

	t.Run("task not found", func(t *testing.T) {
		taskID := primitive.NewObjectID().Hex()
		mockRepo.On("GetByID", mock.Anything, taskID).Return(nil, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/tasks/"+taskID, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCreateTask(t *testing.T) {
	router, mockRepo := setupTest()

	t.Run("successful create task", func(t *testing.T) {
		task := createTestTask()

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.TaskEntity")).Return(task, nil)

		taskJSON, _ := json.Marshal(task)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(taskJSON))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response models.TaskEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, task.Title, response.Title)
		assert.Equal(t, task.StartDate, response.StartDate)
		assert.Equal(t, task.EndDate, response.EndDate)
	})
}

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
		req, _ := http.NewRequest("PUT", "/api/tasks/"+taskID, bytes.NewBuffer(taskJSON))
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

func TestDeleteTask(t *testing.T) {
	router, mockRepo := setupTest()

	t.Run("successful delete task", func(t *testing.T) {
		taskID := primitive.NewObjectID().Hex()
		existingTask := createTestTask()
		existingTask.ID = taskID

		mockRepo.On("GetByID", mock.Anything, taskID).Return(existingTask, nil)
		mockRepo.On("Delete", mock.Anything, taskID).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/tasks/"+taskID, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
