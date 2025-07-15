// filepath: /Users/brandonguigo/workspace/atomic-blend/backend/controllers/users/utils_test.go
package users

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//TODO: replace the tests with mocks of gRPC calls to delete personal data
func TestDeletePersonalData_Success(t *testing.T) {
	// Setup
	controller := &UserController{}
	ctx := &gin.Context{}
	userID := primitive.NewObjectID()

	// Create some mock tasks
	// taskID1 := primitive.NewObjectID().Hex()
	// taskID2 := primitive.NewObjectID().Hex()
	// mockTasks := []*models.TaskEntity{
	// 	{ID: taskID1},
	// 	{ID: taskID2},
	// }

	// Expectations
	// mockTaskRepo.On("GetAll", mock.Anything, &userID).Return(mockTasks, nil)
	// for _, task := range mockTasks {
	// 	mockTaskRepo.On("Delete", mock.Anything, task.ID).Return(nil)
	// }

	// Call the function
	err := controller.DeletePersonalData(ctx, userID)

	// Assertions
	assert.NoError(t, err)
	// mockTaskRepo.AssertExpectations(t)
}

func TestDeletePersonalData_GetAllError(t *testing.T) {
	// Setup
	controller := &UserController{}
	ctx := &gin.Context{}
	userID := primitive.NewObjectID()
	expectedErr := errors.New("database error")

	// Expectations
	// mockTaskRepo.On("GetAll", mock.Anything, &userID).Return(nil, expectedErr)

	// Call the function
	err := controller.DeletePersonalData(ctx, userID)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	// mockTaskRepo.AssertExpectations(t)
}

func TestDeletePersonalData_DeleteError(t *testing.T) {
	// Setup
	controller := &UserController{}
	ctx := &gin.Context{}
	userID := primitive.NewObjectID()
	expectedErr := errors.New("delete error")

	// Create some mock tasks
	// taskID1 := primitive.NewObjectID().Hex()
	// taskID2 := primitive.NewObjectID().Hex()
	// mockTasks := []*models.TaskEntity{
	// 	{ID: taskID1},
	// 	{ID: taskID2},
	// }

	// Expectations
	// mockTaskRepo.On("GetAll", mock.Anything, &userID).Return(mockTasks, nil)
	// mockTaskRepo.On("Delete", mock.Anything, mockTasks[0].ID).Return(nil)         // First delete succeeds
	// mockTaskRepo.On("Delete", mock.Anything, mockTasks[1].ID).Return(expectedErr) // Second delete fails

	// Call the function
	err := controller.DeletePersonalData(ctx, userID)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	// mockTaskRepo.AssertExpectations(t)
}
