package auth

import (
	"atomic_blend_api/models"
	"atomic_blend_api/tests/mocks"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshTokenTestSuite struct {
	suite.Suite
	userRepo     *mocks.MockUserRepository
	userRoleRepo *mocks.MockUserRoleRepository
	controller   *Controller
	router       *gin.Engine
}

func (suite *RefreshTokenTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.userRepo = new(mocks.MockUserRepository)
	suite.userRoleRepo = new(mocks.MockUserRoleRepository)
	suite.controller = NewController(
		suite.userRepo,
		suite.userRoleRepo,
	)
	suite.router = gin.New()
	authGroup := suite.router.Group("/auth")
	authGroup.POST("/refresh", suite.controller.RefreshToken)
}

func (suite *RefreshTokenTestSuite) TestRefreshToken() {
	// User ID for testing
	userID := primitive.NewObjectID()
	email := "test@example.com"
	keySet := models.EncryptionKey{
		UserKey:   "testUserKey123",
		BackupKey: "testBackupKey123",
		UserSalt:  "testUserSalt123",
	}
	now := primitive.NewDateTimeFromTime(time.Now())
	roleID := primitive.NewObjectID()
	role := &models.UserRoleEntity{
		ID:   &roleID,
		Name: "user",
	}

	// Test cases
	testCases := []struct {
		name               string
		setupMocks         func()
		authHeader         string
		expectedStatusCode int
		validateResponse   func(*httptest.ResponseRecorder)
	}{
		{
			name: "Successful_Token_Refresh",
			setupMocks: func() {
				user := &models.UserEntity{
					ID:        &userID,
					Email:     &email,
					KeySet:    &keySet,
					CreatedAt: &now,
					UpdatedAt: &now,
				}

				suite.userRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
				suite.userRoleRepo.On("PopulateRoles", mock.Anything, user).Return(nil).Run(func(args mock.Arguments) {
					user := args.Get(1).(*models.UserEntity)
					user.Roles = []*models.UserRoleEntity{role}
				})
			},
			authHeader:         "Bearer test_refresh_token_" + userID.Hex(),
			expectedStatusCode: http.StatusOK,
			validateResponse: func(recorder *httptest.ResponseRecorder) {
				var response Response
				err := json.Unmarshal(recorder.Body.Bytes(), &response)

				assert.Nil(suite.T(), err)
				assert.NotEmpty(suite.T(), response.AccessToken)
				assert.NotEmpty(suite.T(), response.RefreshToken)
				assert.NotZero(suite.T(), response.ExpiresAt)
				assert.NotNil(suite.T(), response.User)
				assert.Equal(suite.T(), *response.User.ID, userID)
				assert.Equal(suite.T(), *response.User.Email, email)
				
				// Verify KeySet
				assert.NotNil(suite.T(), response.User.KeySet)
				assert.Equal(suite.T(), "testUserKey123", response.User.KeySet.UserKey)
				assert.Equal(suite.T(), "testBackupKey123", response.User.KeySet.BackupKey)
				assert.Equal(suite.T(), "testUserSalt123", response.User.KeySet.UserSalt)
				
				assert.Equal(suite.T(), response.User.Roles[0], role)
			},
		},
		{
			name:               "Missing_Auth_Header",
			setupMocks:         func() {},
			authHeader:         "",
			expectedStatusCode: http.StatusUnauthorized,
			validateResponse: func(recorder *httptest.ResponseRecorder) {
				var response gin.H
				err := json.Unmarshal(recorder.Body.Bytes(), &response)

				assert.Nil(suite.T(), err)
				assert.Contains(suite.T(), response, "error")
				assert.Equal(suite.T(), "Invalid or missing token", response["error"])
			},
		},
		{
			name:               "Invalid_Token_Format",
			setupMocks:         func() {},
			authHeader:         "Invalid token",
			expectedStatusCode: http.StatusUnauthorized,
			validateResponse: func(recorder *httptest.ResponseRecorder) {
				var response gin.H
				err := json.Unmarshal(recorder.Body.Bytes(), &response)

				assert.Nil(suite.T(), err)
				assert.Contains(suite.T(), response, "error")
				assert.Equal(suite.T(), "Invalid or missing token", response["error"])
			},
		},
		{
			name:               "Invalid_Test_Token_ID",
			setupMocks:         func() {},
			authHeader:         "Bearer test_refresh_token_invalid",
			expectedStatusCode: http.StatusUnauthorized,
			validateResponse: func(recorder *httptest.ResponseRecorder) {
				var response gin.H
				err := json.Unmarshal(recorder.Body.Bytes(), &response)

				assert.Nil(suite.T(), err)
				assert.Contains(suite.T(), response, "error")
				assert.Equal(suite.T(), "Invalid test token", response["error"])
			},
		},
		{
			name: "User_Not_Found",
			setupMocks: func() {
				suite.userRepo.On("FindByID", mock.Anything, userID).Return(nil, assert.AnError)
			},
			authHeader:         "Bearer test_refresh_token_" + userID.Hex(),
			expectedStatusCode: http.StatusUnauthorized,
			validateResponse: func(recorder *httptest.ResponseRecorder) {
				var response gin.H
				err := json.Unmarshal(recorder.Body.Bytes(), &response)

				assert.Nil(suite.T(), err)
				assert.Contains(suite.T(), response, "error")
				assert.Equal(suite.T(), "User not found", response["error"])
			},
		},
		{
			name: "Role_Population_Error",
			setupMocks: func() {
				user := &models.UserEntity{
					ID:        &userID,
					Email:     &email,
					KeySet:    &keySet,
					CreatedAt: &now,
					UpdatedAt: &now,
				}

				suite.userRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
				suite.userRoleRepo.On("PopulateRoles", mock.Anything, user).Return(assert.AnError)
			},
			authHeader:         "Bearer test_refresh_token_" + userID.Hex(),
			expectedStatusCode: http.StatusInternalServerError,
			validateResponse: func(recorder *httptest.ResponseRecorder) {
				var response gin.H
				err := json.Unmarshal(recorder.Body.Bytes(), &response)

				assert.Nil(suite.T(), err)
				assert.Contains(suite.T(), response, "error")
				assert.Equal(suite.T(), "Failed to populate user roles", response["error"])
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Reset mocks before each test to clear previous expectations
			suite.userRepo = new(mocks.MockUserRepository)
			suite.userRoleRepo = new(mocks.MockUserRoleRepository)
			suite.controller = NewController(
				suite.userRepo,
				suite.userRoleRepo,
			)
			suite.router = gin.New()
			authGroup := suite.router.Group("/auth")
			authGroup.POST("/refresh", suite.controller.RefreshToken)

			// Setup mocks
			tc.setupMocks()

			// Create request
			req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			recorder := httptest.NewRecorder()

			// Perform request
			suite.router.ServeHTTP(recorder, req)

			// Assert status code
			assert.Equal(t, tc.expectedStatusCode, recorder.Code)

			// Validate response
			tc.validateResponse(recorder)

			// Clear mocks
			suite.userRepo.AssertExpectations(t)
			suite.userRoleRepo.AssertExpectations(t)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	suite.Run(t, new(RefreshTokenTestSuite))
}
