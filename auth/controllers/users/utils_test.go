// filepath: /Users/brandonguigo/workspace/atomic-blend/backend/controllers/users/utils_test.go
package users

import (
	"errors"
	"testing"

	"atomic_blend_api/models"
	"atomic_blend_api/repositories"
	"atomic_blend_api/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock task repository factory for testing
func setupMockTaskRepository() *mocks.MockTaskRepository {
	mockTaskRepo := new(mocks.MockTaskRepository)
	// Override the default factory for testing
	defaultTaskRepositoryFactory = func() repositories.TaskRepositoryInterface {
		return mockTaskRepo
	}
	return mockTaskRepo
}

func TestDeletePersonalData_Success(t *testing.T) {
	// Setup
	mockTaskRepo := setupMockTaskRepository()
	controller := &UserController{}
	ctx := &gin.Context{}
	userID := primitive.NewObjectID()

	// Create some mock tasks
	taskID1 := primitive.NewObjectID().Hex()
	taskID2 := primitive.NewObjectID().Hex()
	mockTasks := []*models.TaskEntity{
		{ID: taskID1},
		{ID: taskID2},
	}

	// Expectations
	mockTaskRepo.On("GetAll", mock.Anything, &userID).Return(mockTasks, nil)
	for _, task := range mockTasks {
		mockTaskRepo.On("Delete", mock.Anything, task.ID).Return(nil)
	}

	// Call the function
	err := controller.DeletePersonalData(ctx, userID)

	// Assertions
	assert.NoError(t, err)
	mockTaskRepo.AssertExpectations(t)
}

func TestDeletePersonalData_GetAllError(t *testing.T) {
	// Setup
	mockTaskRepo := setupMockTaskRepository()
	controller := &UserController{}
	ctx := &gin.Context{}
	userID := primitive.NewObjectID()
	expectedErr := errors.New("database error")

	// Expectations
	mockTaskRepo.On("GetAll", mock.Anything, &userID).Return(nil, expectedErr)

	// Call the function
	err := controller.DeletePersonalData(ctx, userID)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockTaskRepo.AssertExpectations(t)
}

func TestDeletePersonalData_DeleteError(t *testing.T) {
	// Setup
	mockTaskRepo := setupMockTaskRepository()
	controller := &UserController{}
	ctx := &gin.Context{}
	userID := primitive.NewObjectID()
	expectedErr := errors.New("delete error")

	// Create some mock tasks
	taskID1 := primitive.NewObjectID().Hex()
	taskID2 := primitive.NewObjectID().Hex()
	mockTasks := []*models.TaskEntity{
		{ID: taskID1},
		{ID: taskID2},
	}

	// Expectations
	mockTaskRepo.On("GetAll", mock.Anything, &userID).Return(mockTasks, nil)
	mockTaskRepo.On("Delete", mock.Anything, mockTasks[0].ID).Return(nil)         // First delete succeeds
	mockTaskRepo.On("Delete", mock.Anything, mockTasks[1].ID).Return(expectedErr) // Second delete fails

	// Call the function
	err := controller.DeletePersonalData(ctx, userID)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockTaskRepo.AssertExpectations(t)
}
