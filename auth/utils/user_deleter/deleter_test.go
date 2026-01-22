// filepath: /Users/brandonguigo/workspace/atomic-blend/backend/controllers/users/utils_test.go
package userdeleter

import (
	"errors"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/atomic-blend/backend/auth/tests/mocks"
	productivityv1 "github.com/atomic-blend/backend/grpc/gen/productivity/v1"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeletePersonalData_Success(t *testing.T) {
	// Setup
	mockProductivityClient := new(mocks.MockProductivityClient)

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
	mockUserRepo := new(mocks.MockUserRepository)
	mockUserRepo.On("Delete", mock.Anything, userID.Hex()).Return(nil)

	// Call the function
	err := DeletePersonalDataAndUser(userID, mockProductivityClient, mockUserRepo)

	// Assertions
	assert.NoError(t, err)
	mockProductivityClient.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestDeletePersonalData_gRPCError(t *testing.T) {
	// Setup
	mockProductivityClient := new(mocks.MockProductivityClient)

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
	mockUserRepo := new(mocks.MockUserRepository)

	// Call the function
	err := DeletePersonalDataAndUser(userID, mockProductivityClient, mockUserRepo)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockProductivityClient.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestDeletePersonalData_CorrectRequestData(t *testing.T) {
	// Setup
	mockProductivityClient := new(mocks.MockProductivityClient)

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

	// Mock user repository delete
	mockUserRepo := new(mocks.MockUserRepository)
	mockUserRepo.On("Delete", mock.Anything, userID.Hex()).Return(nil)

	// Call the function
	err := DeletePersonalDataAndUser(userID, mockProductivityClient, mockUserRepo)

	// Assertions
	assert.NoError(t, err)
	mockProductivityClient.AssertExpectations(t)
}
