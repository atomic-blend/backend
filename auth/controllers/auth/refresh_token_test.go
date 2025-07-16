package auth

import (
	"atomic-blend/backend/auth/models"
	"atomic-blend/backend/auth/repositories"
	"atomic-blend/backend/auth/tests/utils/inmemorymongo"
	"atomic-blend/backend/auth/utils/db"
	"atomic-blend/backend/auth/utils/jwt"
	"context"
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
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// setupTestDB creates an in-memory MongoDB for testing
func setupTestDB(t *testing.T) (repositories.UserRepositoryInterface, repositories.UserRoleRepositoryInterface, func()) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Start in-memory MongoDB server
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	// Get MongoDB connection URI
	mongoURI := mongoServer.URI()

	// Connect to the in-memory MongoDB
	client, err := inmemorymongo.ConnectToInMemoryDB(mongoURI)
	require.NoError(t, err)

	// Get database reference and create repositories
	database := client.Database("test_db")
	userRepo := repositories.NewUserRepository(database)
	userRoleRepo := repositories.NewUserRoleRepository(database)

	// Set the global database for the subscription function to use
	db.Database = database

	// Return cleanup function
	cleanup := func() {
		// Reset global database
		db.Database = nil
		client.Disconnect(context.Background())
		mongoServer.Stop()
	}

	return userRepo, userRoleRepo, cleanup
}

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
	accessToken, err := jwt.GenerateToken(ctx, userID, jwt.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := jwt.GenerateToken(ctx, userID, jwt.RefreshToken)
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
	userRepo repositories.UserRepositoryInterface,
	userRoleRepo repositories.UserRoleRepositoryInterface,
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
	userRepo     repositories.UserRepositoryInterface
	userRoleRepo repositories.UserRoleRepositoryInterface
	controller   *TestController
	router       *gin.Engine
	cleanup      func()
}

// createTestUser creates a test user with optional purchase data
func createTestUser(t *testing.T, repo repositories.UserRepositoryInterface, purchases []*models.PurchaseEntity) *models.UserEntity {
	email := "test@example.com"
	password := "testpassword"
	keySet := models.EncryptionKey{
		UserKey:      "testUserKey123",
		BackupKey:    "testBackupKey123",
		Salt:         "testSalt123",
		MnemonicSalt: "testMnemonicSalt123",
	}

	user := &models.UserEntity{
		Email:     &email,
		Password:  &password,
		KeySet:    &keySet,
		Purchases: purchases,
	}

	created, err := repo.Create(context.Background(), user)
	require.NoError(t, err)
	require.NotNil(t, created.ID)

	return created
}

// createTestRole creates a test role
func createTestRole(t *testing.T, repo repositories.UserRoleRepositoryInterface) *models.UserRoleEntity {
	role := &models.UserRoleEntity{
		Name: "user",
	}

	created, err := repo.Create(context.Background(), role)
	require.NoError(t, err)
	require.NotNil(t, created.ID)

	return created
}

func (suite *RefreshTokenTestSuite) SetupTest() {
	os.Setenv("SSO_SECRET", "test-secret-key")
	gin.SetMode(gin.TestMode)

	// Setup real database for testing
	var cleanup func()
	suite.userRepo, suite.userRoleRepo, cleanup = setupTestDB(suite.T())
	suite.cleanup = cleanup
}

func (suite *RefreshTokenTestSuite) TearDownTest() {
	if suite.cleanup != nil {
		suite.cleanup()
	}
}

func (suite *RefreshTokenTestSuite) TestRefreshToken() {
	// Test cases
	testCases := []struct {
		name               string
		mockJWTValidator   func(string, jwt.TokenType) (*jwtlib.MapClaims, error)
		setupTestData      func() (primitive.ObjectID, *models.UserEntity)
		authHeader         string
		expectedStatusCode int
		validateResponse   func(*httptest.ResponseRecorder, primitive.ObjectID, *models.UserEntity)
	}{
		{
			name: "Successful_Token_Refresh_With_Active_Subscription",
			mockJWTValidator: func(tokenString string, tokenType jwt.TokenType) (*jwtlib.MapClaims, error) {
				// This will be set dynamically in setupTestData
				return nil, nil
			},
			setupTestData: func() (primitive.ObjectID, *models.UserEntity) {
				// Create purchase with active subscription
				futureTime := time.Now().Add(24 * time.Hour).UnixMilli() // 1 day in future
				purchaseType := "REVENUE_CAT"
				purchases := []*models.PurchaseEntity{
					{
						ID:   primitive.NewObjectID(),
						Type: &purchaseType,
						PurchaseData: models.RevenueCatPurchaseData{
							ExpirationAtMs: futureTime,
							ProductID:      "test_product",
							AppUserID:      "test_user",
						},
						CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
						UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
					},
				}

				user := createTestUser(suite.T(), suite.userRepo, purchases)
				role := createTestRole(suite.T(), suite.userRoleRepo)

				// Assign role to user
				user.Roles = []*models.UserRoleEntity{role}

				return *user.ID, user
			},
			authHeader:         "Bearer valid_refresh_token",
			expectedStatusCode: http.StatusOK,
			validateResponse: func(recorder *httptest.ResponseRecorder, userID primitive.ObjectID, user *models.UserEntity) {
				var response Response
				err := json.NewDecoder(recorder.Body).Decode(&response)

				assert.NoError(suite.T(), err)
				assert.NotEmpty(suite.T(), response.AccessToken)
				assert.NotEmpty(suite.T(), response.RefreshToken)
				assert.NotZero(suite.T(), response.ExpiresAt)
				assert.NotNil(suite.T(), response.User)

				if response.User != nil {
					assert.Equal(suite.T(), userID, *response.User.ID)
					assert.Equal(suite.T(), *user.Email, *response.User.Email)
					assert.NotNil(suite.T(), response.User.KeySet)
					if response.User.KeySet != nil {
						assert.Equal(suite.T(), user.KeySet.UserKey, response.User.KeySet.UserKey)
						assert.Equal(suite.T(), user.KeySet.BackupKey, response.User.KeySet.BackupKey)
						assert.Equal(suite.T(), user.KeySet.Salt, response.User.KeySet.Salt)
						assert.Equal(suite.T(), user.KeySet.MnemonicSalt, response.User.KeySet.MnemonicSalt)
					}
				}
			},
		},
		{
			name: "Successful_Token_Refresh_Without_Subscription",
			mockJWTValidator: func(tokenString string, tokenType jwt.TokenType) (*jwtlib.MapClaims, error) {
				return nil, nil
			},
			setupTestData: func() (primitive.ObjectID, *models.UserEntity) {
				// Create user without purchases
				user := createTestUser(suite.T(), suite.userRepo, nil)
				role := createTestRole(suite.T(), suite.userRoleRepo)

				// Assign role to user
				user.Roles = []*models.UserRoleEntity{role}

				return *user.ID, user
			},
			authHeader:         "Bearer valid_refresh_token",
			expectedStatusCode: http.StatusOK,
			validateResponse: func(recorder *httptest.ResponseRecorder, userID primitive.ObjectID, user *models.UserEntity) {
				var response Response
				err := json.NewDecoder(recorder.Body).Decode(&response)

				assert.NoError(suite.T(), err)
				assert.NotEmpty(suite.T(), response.AccessToken)
				assert.NotEmpty(suite.T(), response.RefreshToken)
				assert.NotZero(suite.T(), response.ExpiresAt)
				assert.NotNil(suite.T(), response.User)
			},
		},
		{
			name:               "Missing_Auth_Header",
			mockJWTValidator:   nil,
			setupTestData:      func() (primitive.ObjectID, *models.UserEntity) { return primitive.ObjectID{}, nil },
			authHeader:         "",
			expectedStatusCode: http.StatusUnauthorized,
			validateResponse: func(recorder *httptest.ResponseRecorder, userID primitive.ObjectID, user *models.UserEntity) {
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
			setupTestData:      func() (primitive.ObjectID, *models.UserEntity) { return primitive.ObjectID{}, nil },
			authHeader:         "Invalid token",
			expectedStatusCode: http.StatusUnauthorized,
			validateResponse: func(recorder *httptest.ResponseRecorder, userID primitive.ObjectID, user *models.UserEntity) {
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
			setupTestData:      func() (primitive.ObjectID, *models.UserEntity) { return primitive.ObjectID{}, nil },
			authHeader:         "Bearer invalid_token",
			expectedStatusCode: http.StatusUnauthorized,
			validateResponse: func(recorder *httptest.ResponseRecorder, userID primitive.ObjectID, user *models.UserEntity) {
				var response gin.H
				err := json.Unmarshal(recorder.Body.Bytes(), &response)

				assert.Nil(suite.T(), err)
				assert.Contains(suite.T(), response, "error")
				assert.Equal(suite.T(), "Invalid refresh token", response["error"])
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Setup test data
			userID, user := tc.setupTestData()

			// Update the JWT validator to use the actual userID if needed
			if tc.mockJWTValidator != nil && !userID.IsZero() {
				tc.mockJWTValidator = func(tokenString string, tokenType jwt.TokenType) (*jwtlib.MapClaims, error) {
					claims := jwtlib.MapClaims{
						"user_id": userID.Hex(),
						"type":    string(jwt.RefreshToken),
					}
					return &claims, nil
				}
			}

			// Create a new test controller with the mock JWT validator
			suite.controller = NewTestController(
				suite.userRepo,
				suite.userRoleRepo,
				tc.mockJWTValidator,
			)

			suite.router = gin.New()
			authGroup := suite.router.Group("/auth")
			authGroup.POST("/refresh", suite.controller.RefreshToken)

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
			tc.validateResponse(recorder, userID, user)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	os.Setenv("SSO_SECRET", "test-secret-key")
	defer os.Unsetenv("SSO_SECRET")
	suite.Run(t, new(RefreshTokenTestSuite))
}
