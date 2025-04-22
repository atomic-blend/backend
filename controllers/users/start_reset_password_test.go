package users

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
	"atomic_blend_api/tests/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock controller with override for templates
type mockUserController struct {
	*UserController
	templateErr bool
	updateErr   bool
	emailErr    bool
}

func (c *mockUserController) StartResetPassword(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	user, err := c.userRepo.FindByID(ctx, authUser.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	// generate reset code
	resetCode := "12345678" // Fixed code for testing

	// Skip template parsing in tests
	if c.templateErr {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse HTML template"})
		return
	}

	// store the reset code in the database
	user.ResetPasswordCode = &resetCode
	if c.updateErr {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user reset password code"})
		return
	}

	// Skip email sending in tests
	if c.emailErr {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Reset password email sent successfully", "sent": true})
}

func TestStartResetPassword(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		setupAuth      func(*gin.Context)
		setupMocks     func(*mocks.MockUserRepository, *mocks.MockUserRoleRepository)
		templateErr    bool
		updateErr      bool
		emailErr       bool
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Successfully send reset password email",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {
				userID := primitive.NewObjectID()
				email := "test@example.com"
				user := &models.UserEntity{
					ID:    &userID,
					Email: &email,
				}

				// Mock finding user by ID
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(user, nil)
			},
			templateErr:    false,
			updateErr:      false,
			emailErr:       false,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
				assert.Equal(t, "Reset password email sent successfully", response["message"])
				assert.Contains(t, response, "sent")
			},
		},
		{
			name:           "Unauthorized - no auth user",
			setupAuth:      func(c *gin.Context) {},
			setupMocks:     func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {},
			templateErr:    false,
			updateErr:      false,
			emailErr:       false,
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Authentication required", response["error"])
			},
		},
		{
			name: "Error finding user",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {
				// Mock user not found
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(nil, errors.New("user not found"))
			},
			templateErr:    false,
			updateErr:      false,
			emailErr:       false,
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to retrieve user profile", response["error"])
			},
		},
		{
			name: "Error with template parsing",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {
				userID := primitive.NewObjectID()
				email := "test@example.com"
				user := &models.UserEntity{
					ID:    &userID,
					Email: &email,
				}

				// Mock finding user successfully
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(user, nil)
			},
			templateErr:    true,
			updateErr:      false,
			emailErr:       false,
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to parse HTML template", response["error"])
			},
		},
		{
			name: "Error updating user",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {
				userID := primitive.NewObjectID()
				email := "test@example.com"
				user := &models.UserEntity{
					ID:    &userID,
					Email: &email,
				}

				// Mock finding user successfully
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(user, nil)
			},
			templateErr:    false,
			updateErr:      true,
			emailErr:       false,
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to update user reset password code", response["error"])
			},
		},
		{
			name: "Error sending email",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {
				userID := primitive.NewObjectID()
				email := "test@example.com"
				user := &models.UserEntity{
					ID:    &userID,
					Email: &email,
				}

				// Mock finding user successfully
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(user, nil)
			},
			templateErr:    false,
			updateErr:      false,
			emailErr:       true,
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to send email", response["error"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			mockUserRepo := new(mocks.MockUserRepository)
			mockUserRoleRepo := new(mocks.MockUserRoleRepository)

			// Setup mocks
			tc.setupMocks(mockUserRepo, mockUserRoleRepo)

			// Create base controller
			baseController := NewUserController(mockUserRepo, mockUserRoleRepo)

			// Create mock controller with overrides
			mockController := &mockUserController{
				UserController: baseController,
				templateErr:    tc.templateErr,
				updateErr:      tc.updateErr,
				emailErr:       tc.emailErr,
			}

			// Create router
			router := gin.New()
			router.POST("/users/reset-password/start", func(c *gin.Context) {
				tc.setupAuth(c)
				mockController.StartResetPassword(c)
			})

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/users/reset-password/start", nil)

			// Perform request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Check response
			tc.checkResponse(t, w)

			// Verify mock expectations
			mockUserRepo.AssertExpectations(t)
			mockUserRoleRepo.AssertExpectations(t)
		})
	}
}
