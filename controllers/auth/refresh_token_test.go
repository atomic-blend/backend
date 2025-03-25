package auth

import (
	"atomic_blend_api/models"
	"atomic_blend_api/tests/mocks"
	"atomic_blend_api/utils/jwt"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create a test-specific controller that can use our test JWT validator
type TestController struct {
	Controller
	mockJWTValidator func(string, jwt.TokenType) (*jwtlib.MapClaims, error)
}

// Override the RefreshToken method to use our mock JWT validator
func (c *TestController) RefreshToken(ctx *gin.Context) {
	// Get token from Authorization header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
		return
	}

	tokenString := authHeader[7:]

	// Use our test JWT validator instead of the real one
	var claims *jwtlib.MapClaims
	var err error

	if c.mockJWTValidator != nil {
		claims, err = c.mockJWTValidator(tokenString, jwt.RefreshToken)
	} else {
		claims, err = jwt.ValidateToken(tokenString, jwt.RefreshToken)
	}

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Extract user ID from claims
	userIDStr, ok := (*claims)["user_id"].(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	// Find user to ensure they still exist and get their current data
	user, err := c.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Populate user roles
	if err = c.userRoleRepo.PopulateRoles(ctx, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to populate user roles"})
		return
	}

	// Generate new tokens
	accessToken, err := jwt.GenerateToken(userID, jwt.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := jwt.GenerateToken(userID, jwt.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Return user and new tokens
	responseSafeUser := &models.UserEntity{
		ID:        user.ID,
		Email:     user.Email,
		KeySet:    user.KeySet,
		Roles:     user.Roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, Response{
		User:         responseSafeUser,
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
		ExpiresAt:    accessToken.ExpiresAt.Unix(),
	})
}

// Create a new TestController
func NewTestController(
	userRepo *mocks.MockUserRepository,
	userRoleRepo *mocks.MockUserRoleRepository,
	mockJWTValidator func(string, jwt.TokenType) (*jwtlib.MapClaims, error),
) *TestController {
	return &TestController{
		Controller: Controller{
			userRepo:     userRepo,
			userRoleRepo: userRoleRepo,
		},
		mockJWTValidator: mockJWTValidator,
	}
}

type RefreshTokenTestSuite struct {
	suite.Suite
	userRepo     *mocks.MockUserRepository
	userRoleRepo *mocks.MockUserRoleRepository
	controller   *TestController
	router       *gin.Engine
}

func (suite *RefreshTokenTestSuite) SetupTest() {
	os.Setenv("SSO_SECRET", "test-secret-key")
	gin.SetMode(gin.TestMode)
	suite.userRepo = new(mocks.MockUserRepository)
	suite.userRoleRepo = new(mocks.MockUserRoleRepository)
}

func (suite *RefreshTokenTestSuite) TestRefreshToken() {
	// User ID for testing
	userID := primitive.NewObjectID()
	email := "test@example.com"
	keySet := models.EncryptionKey{
		UserKey:      "testUserKey123",
		BackupKey:    "testBackupKey123",
		Salt:         "testSalt123",
		MnemonicSalt: "testMnemonicSalt123",
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
		mockJWTValidator   func(string, jwt.TokenType) (*jwtlib.MapClaims, error)
		setupMocks         func()
		authHeader         string
		expectedStatusCode int
		validateResponse   func(*httptest.ResponseRecorder)
	}{
		{
			name: "Successful_Token_Refresh",
			mockJWTValidator: func(tokenString string, tokenType jwt.TokenType) (*jwtlib.MapClaims, error) {
				claims := jwtlib.MapClaims{
					"user_id": userID.Hex(),
					"type":    string(jwt.RefreshToken),
				}
				return &claims, nil
			},
			setupMocks: func() {
				user := &models.UserEntity{
					ID:        &userID,
					Email:     &email,
					KeySet:    &keySet,
					CreatedAt: &now,
					UpdatedAt: &now,
					Roles:     []*models.UserRoleEntity{role},
				}

				suite.userRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
				suite.userRoleRepo.On("PopulateRoles", mock.Anything, user).Return(nil)
			},
			authHeader:         "Bearer valid_refresh_token", // The actual token content doesn't matter since we mock validation
			expectedStatusCode: http.StatusOK,
			validateResponse: func(recorder *httptest.ResponseRecorder) {
				var response Response
				err := json.NewDecoder(recorder.Body).Decode(&response)

				assert.NoError(suite.T(), err)
				assert.NotEmpty(suite.T(), response.AccessToken)
				assert.NotEmpty(suite.T(), response.RefreshToken)
				assert.NotZero(suite.T(), response.ExpiresAt)
				assert.NotNil(suite.T(), response.User)

				if response.User != nil {
					assert.Equal(suite.T(), userID, *response.User.ID)
					assert.Equal(suite.T(), email, *response.User.Email)

					// Verify KeySet
					assert.NotNil(suite.T(), response.User.KeySet)
					if response.User.KeySet != nil {
						assert.Equal(suite.T(), keySet.UserKey, response.User.KeySet.UserKey)
						assert.Equal(suite.T(), keySet.BackupKey, response.User.KeySet.BackupKey)
						assert.Equal(suite.T(), keySet.Salt, response.User.KeySet.Salt)
						assert.Equal(suite.T(), keySet.MnemonicSalt, response.User.KeySet.MnemonicSalt)
					}

					assert.NotEmpty(suite.T(), response.User.Roles)
					if len(response.User.Roles) > 0 {
						assert.Equal(suite.T(), role, response.User.Roles[0])
					}
				}
			},
		},
		{
			name:               "Missing_Auth_Header",
			mockJWTValidator:   nil,
			setupMocks:         func() {},
			authHeader:         "",
			expectedStatusCode: http.StatusUnauthorized,
			validateResponse: func(recorder *httptest.ResponseRecorder) {
				var response gin.H
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(suite.T(), err)
				assert.Contains(suite.T(), response, "error")
				assert.Equal(suite.T(), "Invalid or missing token", response["error"])
			},
		},
		{
			name:               "Invalid_Token_Format",
			mockJWTValidator:   nil,
			setupMocks:         func() {},
			authHeader:         "Invalid token",
			expectedStatusCode: http.StatusUnauthorized,
			validateResponse: func(recorder *httptest.ResponseRecorder) {
				var response gin.H
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(suite.T(), err)
				assert.Contains(suite.T(), response, "error")
				assert.Equal(suite.T(), "Invalid or missing token", response["error"])
			},
		},
		{
			name: "Invalid_Test_Token_ID",
			mockJWTValidator: func(tokenString string, tokenType jwt.TokenType) (*jwtlib.MapClaims, error) {
				return nil, errors.New("Invalid token")
			},
			setupMocks:         func() {},
			authHeader:         "Bearer invalid_token",
			expectedStatusCode: http.StatusUnauthorized,
			validateResponse: func(recorder *httptest.ResponseRecorder) {
				var response gin.H
				err := json.Unmarshal(recorder.Body.Bytes(), &response)

				assert.Nil(suite.T(), err)
				assert.Contains(suite.T(), response, "error")
				assert.Equal(suite.T(), "Invalid refresh token", response["error"])
			},
		},
		{
			name: "User_Not_Found",
			mockJWTValidator: func(tokenString string, tokenType jwt.TokenType) (*jwtlib.MapClaims, error) {
				claims := jwtlib.MapClaims{
					"user_id": userID.Hex(),
					"type":    string(jwt.RefreshToken),
				}
				return &claims, nil
			},
			setupMocks: func() {
				suite.userRepo.On("FindByID", mock.Anything, userID).Return(nil, errors.New("User not found"))
			},
			authHeader:         "Bearer valid_token_but_user_not_found",
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
			mockJWTValidator: func(tokenString string, tokenType jwt.TokenType) (*jwtlib.MapClaims, error) {
				claims := jwtlib.MapClaims{
					"user_id": userID.Hex(),
					"type":    string(jwt.RefreshToken),
				}
				return &claims, nil
			},
			setupMocks: func() {
				user := &models.UserEntity{
					ID:        &userID,
					Email:     &email,
					KeySet:    &keySet,
					CreatedAt: &now,
					UpdatedAt: &now,
				}

				suite.userRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
				suite.userRoleRepo.On("PopulateRoles", mock.Anything, user).Return(errors.New("Failed to populate user roles"))
			},
			authHeader:         "Bearer valid_token_but_role_error",
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

			// Create a new test controller with the mock JWT validator
			suite.controller = NewTestController(
				suite.userRepo,
				suite.userRoleRepo,
				tc.mockJWTValidator,
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

			// Assert status code and validate response
			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			tc.validateResponse(recorder)

			// Clear mocks
			suite.userRepo.AssertExpectations(t)
			suite.userRoleRepo.AssertExpectations(t)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	os.Setenv("SSO_SECRET", "test-secret-key")
	defer os.Unsetenv("SSO_SECRET")
	suite.Run(t, new(RefreshTokenTestSuite))
}
