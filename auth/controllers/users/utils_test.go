// filepath: /Users/brandonguigo/workspace/atomic-blend/backend/controllers/users/utils_test.go
package users

import (
	"errors"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/atomic-blend/backend/auth/tests/mocks"
	"github.com/atomic-blend/backend/grpc/gen/productivity/v1"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeletePersonalData_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(mocks.MockUserRepository)
	mockUserRoleRepo := new(mocks.MockUserRoleRepository)
	mockProductivityClient := new(mocks.MockProductivityClient)

	controller := NewUserController(mockUserRepo, mockUserRoleRepo, mockProductivityClient)

	// Create gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("DELETE", "/", nil)
	ctx.Request = req

	userID := primitive.NewObjectID()

	// Mock successful gRPC call
	mockResponse := &connect.Response[productivityv1.DeleteUserDataResponse]{}
	mockProductivityClient.On("DeleteUserData", mock.Anything, mock.Anything).Return(mockResponse, nil)

	// Call the function
	err := controller.DeletePersonalData(ctx, userID)

	// Assertions
	assert.NoError(t, err)
	mockProductivityClient.AssertExpectations(t)
}

func TestDeletePersonalData_gRPCError(t *testing.T) {
	// Setup
	mockUserRepo := new(mocks.MockUserRepository)
	mockUserRoleRepo := new(mocks.MockUserRoleRepository)
	mockProductivityClient := new(mocks.MockProductivityClient)

	controller := NewUserController(mockUserRepo, mockUserRoleRepo, mockProductivityClient)

	// Create gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("DELETE", "/", nil)
	ctx.Request = req

	userID := primitive.NewObjectID()
	expectedErr := errors.New("gRPC connection error")

	// Mock failing gRPC call
	mockProductivityClient.On("DeleteUserData", mock.Anything, mock.Anything).Return(nil, expectedErr)

	// Call the function
	err := controller.DeletePersonalData(ctx, userID)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockProductivityClient.AssertExpectations(t)
}

func TestDeletePersonalData_CorrectRequestData(t *testing.T) {
	// Setup
	mockUserRepo := new(mocks.MockUserRepository)
	mockUserRoleRepo := new(mocks.MockUserRoleRepository)
	mockProductivityClient := new(mocks.MockProductivityClient)

	controller := NewUserController(mockUserRepo, mockUserRoleRepo, mockProductivityClient)

	// Create gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("DELETE", "/", nil)
	ctx.Request = req

	userID := primitive.NewObjectID()

	// Mock with basic validation that the request is properly formatted
	mockResponse := &connect.Response[productivityv1.DeleteUserDataResponse]{}
	mockProductivityClient.On("DeleteUserData",
		mock.Anything,
		mock.MatchedBy(func(req *connect.Request[productivityv1.DeleteUserDataRequest]) bool {
			// Validate that the request message is not nil
			return req.Msg != nil
		})).Return(mockResponse, nil)

	// Call the function
	err := controller.DeletePersonalData(ctx, userID)

	// Assertions
	assert.NoError(t, err)
	mockProductivityClient.AssertExpectations(t)
}
