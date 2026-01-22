package cron

import (
	"errors"
	"testing"

	"github.com/atomic-blend/backend/auth/tests/mocks"
	"github.com/atomic-blend/backend/shared/models"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeleteInactiveSubscriptionUsers(t *testing.T) {
	// Set Gin to test mode (though not used here, following pattern)
	// Actually, not needed, but keeping similar

	testCases := []struct {
		name          string
		setupMocks    func(*mocks.MockUserRepository, *mocks.MockProductivityClient)
		expectedCalls func(*mocks.MockUserRepository, *mocks.MockProductivityClient)
	}{
		{
			name: "Successfully delete multiple users",
			setupMocks: func(userRepo *mocks.MockUserRepository, productivityClient *mocks.MockProductivityClient) {
				userID1 := primitive.NewObjectID()
				userID2 := primitive.NewObjectID()
				users := []*models.UserEntity{
					{ID: &userID1},
					{ID: &userID2},
				}

				userRepo.On("FindInactiveSubscriptionUsers", mock.Anything, 7).Return(users, nil)

				// Mock successful deletion for each user
				productivityClient.On("DeleteUserData", mock.Anything, mock.Anything).Return(nil, nil).Times(2)
				userRepo.On("Delete", mock.Anything, userID1.Hex()).Return(nil)
				userRepo.On("Delete", mock.Anything, userID2.Hex()).Return(nil)
			},
			expectedCalls: func(userRepo *mocks.MockUserRepository, productivityClient *mocks.MockProductivityClient) {
				userRepo.AssertExpectations(t)
				productivityClient.AssertExpectations(t)
			},
		},
		{
			name: "No users found",
			setupMocks: func(userRepo *mocks.MockUserRepository, productivityClient *mocks.MockProductivityClient) {
				userRepo.On("FindInactiveSubscriptionUsers", mock.Anything, 7).Return([]*models.UserEntity{}, nil)
				// No deletion calls expected
			},
			expectedCalls: func(userRepo *mocks.MockUserRepository, productivityClient *mocks.MockProductivityClient) {
				userRepo.AssertExpectations(t)
				productivityClient.AssertExpectations(t)
			},
		},
		{
			name: "Error finding inactive users",
			setupMocks: func(userRepo *mocks.MockUserRepository, productivityClient *mocks.MockProductivityClient) {
				userRepo.On("FindInactiveSubscriptionUsers", mock.Anything, 7).Return(nil, errors.New("db error"))
				// No deletion calls expected
			},
			expectedCalls: func(userRepo *mocks.MockUserRepository, productivityClient *mocks.MockProductivityClient) {
				userRepo.AssertExpectations(t)
				productivityClient.AssertExpectations(t)
			},
		},
		{
			name: "Error deleting one user, continues to next",
			setupMocks: func(userRepo *mocks.MockUserRepository, productivityClient *mocks.MockProductivityClient) {
				userID1 := primitive.NewObjectID()
				userID2 := primitive.NewObjectID()
				users := []*models.UserEntity{
					{ID: &userID1},
					{ID: &userID2},
				}

				userRepo.On("FindInactiveSubscriptionUsers", mock.Anything, 7).Return(users, nil)

				// Mock failure for first user
				productivityClient.On("DeleteUserData", mock.Anything, mock.Anything).Return(nil, errors.New("grpc error")).Once()
				// No delete call for first user

				// Mock success for second user
				productivityClient.On("DeleteUserData", mock.Anything, mock.Anything).Return(nil, nil).Once()
				userRepo.On("Delete", mock.Anything, userID2.Hex()).Return(nil)
			},
			expectedCalls: func(userRepo *mocks.MockUserRepository, productivityClient *mocks.MockProductivityClient) {
				userRepo.AssertExpectations(t)
				productivityClient.AssertExpectations(t)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			mockUserRepo := new(mocks.MockUserRepository)
			mockProductivityClient := new(mocks.MockProductivityClient)

			// Setup mocks
			tc.setupMocks(mockUserRepo, mockProductivityClient)

			// Call the function
			DeleteInactiveSubscriptionUsers(mockUserRepo, mockProductivityClient)

			// Check expectations
			tc.expectedCalls(mockUserRepo, mockProductivityClient)
		})
	}
}
