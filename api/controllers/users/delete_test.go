package users

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
	"atomic_blend_api/repositories"
	"atomic_blend_api/tests/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeleteAccount(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Save original factory to restore later
	originalFactory := defaultTaskRepositoryFactory
	defer func() {
		// Restore original factory after test
		defaultTaskRepositoryFactory = originalFactory
	}()

	// Test cases
	testCases := []struct {
		name           string
		setupAuth      func(*gin.Context)
		setupMocks     func(*mocks.MockUserRepository, *mocks.MockUserRoleRepository, *mocks.MockTaskRepository, primitive.ObjectID)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful Account Deletion",
			setupAuth: func(c *gin.Context) {
				// Set up user authentication
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{
					UserID: userID,
				})
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockUserRoleRepo *mocks.MockUserRoleRepository, mockTaskRepo *mocks.MockTaskRepository, userID primitive.ObjectID) {
				// Setup user repo mock to return a user
				email := "test@example.com"
				user := &models.UserEntity{
					ID:    &userID,
					Email: &email,
				}
				mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil)

				// Setup mock for user's tasks
				tasks := []*models.TaskEntity{
					{
						ID:    primitive.NewObjectID().Hex(),
						Title: "Test Task 1",
						User:  userID,
					},
					{
						ID:    primitive.NewObjectID().Hex(),
						Title: "Test Task 2",
						User:  userID,
					},
				}

				// Mock task operations
				mockTaskRepo.On("GetAll", mock.Anything, &userID).Return(tasks, nil)
				for _, task := range tasks {
					mockTaskRepo.On("Delete", mock.Anything, task.ID).Return(nil)
				}

				// Setup delete to succeed
				mockUserRepo.On("Delete", mock.Anything, userID.Hex()).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Account successfully deleted",
		},
		{
			name: "No Authentication",
			setupAuth: func(c *gin.Context) {
				// Don't set any auth info
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockUserRoleRepo *mocks.MockUserRoleRepository, mockTaskRepo *mocks.MockTaskRepository, userID primitive.ObjectID) {
				// No mocks need to be set up because we should fail before repo is called
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Authentication required",
		},
		{
			name: "User Not Found",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{
					UserID: userID,
				})
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockUserRoleRepo *mocks.MockUserRoleRepository, mockTaskRepo *mocks.MockTaskRepository, userID primitive.ObjectID) {
				// Setup user repo mock to return nil (user not found)
				mockUserRepo.On("FindByID", mock.Anything, userID).Return(nil, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "User not found",
		},
		{
			name: "Database Error on FindByID",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{
					UserID: userID,
				})
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockUserRoleRepo *mocks.MockUserRoleRepository, mockTaskRepo *mocks.MockTaskRepository, userID primitive.ObjectID) {
				// Setup user repo mock to return an error
				mockUserRepo.On("FindByID", mock.Anything, userID).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to verify user",
		},
		{
			name: "Error Getting Tasks",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{
					UserID: userID,
				})
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockUserRoleRepo *mocks.MockUserRoleRepository, mockTaskRepo *mocks.MockTaskRepository, userID primitive.ObjectID) {
				// Setup user repo mock to return a user
				email := "test@example.com"
				user := &models.UserEntity{
					ID:    &userID,
					Email: &email,
				}
				mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil)

				// Setup error for GetAll tasks
				mockTaskRepo.On("GetAll", mock.Anything, &userID).Return(nil, errors.New("error fetching tasks"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to delete personal data",
		},
		{
			name: "Error Deleting Task",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{
					UserID: userID,
				})
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockUserRoleRepo *mocks.MockUserRoleRepository, mockTaskRepo *mocks.MockTaskRepository, userID primitive.ObjectID) {
				// Setup user repo mock to return a user
				email := "test@example.com"
				user := &models.UserEntity{
					ID:    &userID,
					Email: &email,
				}
				mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil)

				// Setup tasks but make one fail to delete
				taskID1 := primitive.NewObjectID().Hex()
				taskID2 := primitive.NewObjectID().Hex()
				tasks := []*models.TaskEntity{
					{
						ID:    taskID1,
						Title: "Test Task 1",
						User:  userID,
					},
					{
						ID:    taskID2,
						Title: "Test Task 2",
						User:  userID,
					},
				}

				mockTaskRepo.On("GetAll", mock.Anything, &userID).Return(tasks, nil)
				mockTaskRepo.On("Delete", mock.Anything, taskID1).Return(nil)
				mockTaskRepo.On("Delete", mock.Anything, taskID2).Return(errors.New("error deleting task"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to delete personal data",
		},
		{
			name: "Database Error on Delete User",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{
					UserID: userID,
				})
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockUserRoleRepo *mocks.MockUserRoleRepository, mockTaskRepo *mocks.MockTaskRepository, userID primitive.ObjectID) {
				// Setup user repo mock to return a user but fail on delete
				email := "test@example.com"
				user := &models.UserEntity{
					ID:    &userID,
					Email: &email,
				}
				mockUserRepo.On("FindByID", mock.Anything, userID).Return(user, nil)

				// No tasks to delete
				mockTaskRepo.On("GetAll", mock.Anything, &userID).Return([]*models.TaskEntity{}, nil)

				// Fail on user deletion
				mockUserRepo.On("Delete", mock.Anything, userID.Hex()).Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to delete account",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock repositories
			mockUserRepo := new(mocks.MockUserRepository)
			mockUserRoleRepo := new(mocks.MockUserRoleRepository)
			mockTaskRepo := new(mocks.MockTaskRepository)

			// Override the task repository factory to return our mock
			defaultTaskRepositoryFactory = func() repositories.TaskRepositoryInterface {
				return mockTaskRepo
			}

			// Create test context
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			req := httptest.NewRequest("DELETE", "/users/me", nil)
			ctx.Request = req

			// Setup authentication in context
			tc.setupAuth(ctx)

			// Get userID from context if available
			var userID primitive.ObjectID
			if authUser, exists := ctx.Get("authUser"); exists {
				if userAuthInfo, ok := authUser.(*auth.UserAuthInfo); ok {
					userID = userAuthInfo.UserID
				}
			}

			// Setup mocks
			tc.setupMocks(mockUserRepo, mockUserRoleRepo, mockTaskRepo, userID)

			// Create controller
			controller := NewUserController(mockUserRepo, mockUserRoleRepo)

			// Execute the function under test
			controller.DeleteAccount(ctx)

			// Assert expectations
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)
			mockUserRepo.AssertExpectations(t)
			mockTaskRepo.AssertExpectations(t)
		})
	}
}
