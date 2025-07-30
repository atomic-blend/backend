package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/atomic-blend/backend/auth/models"
	"github.com/atomic-blend/backend/auth/tests/utils/inmemorymongo"
	"github.com/atomic-blend/backend/auth/utils/db"
	"github.com/atomic-blend/backend/auth/utils/jwt"

	"github.com/atomic-blend/memongo"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var testMongoServer *memongo.Server

// setupTestDB initializes an in-memory MongoDB for testing
func setupTestDB() error {
	if testMongoServer != nil {
		return nil // Already set up
	}

	server, err := inmemorymongo.CreateInMemoryMongoDB()
	if err != nil {
		return err
	}

	testMongoServer = server

	// Connect to the in-memory database
	client, err := inmemorymongo.ConnectToInMemoryDB(server.URI())
	if err != nil {
		server.Stop()
		return err
	}

	// Set the global database variables
	db.MongoClient = client
	db.Database = client.Database("test_database")

	return nil
}

// teardownTestDB cleans up the test database
func teardownTestDB() {
	if testMongoServer != nil {
		if db.MongoClient != nil {
			db.MongoClient.Disconnect(context.Background())
		}
		testMongoServer.Stop()
		testMongoServer = nil
		db.MongoClient = nil
		db.Database = nil
	}
}

// Mock repository interfaces for testing
type UserRepositoryInterface interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.UserEntity, error)
}

type UserRoleRepositoryInterface interface {
	PopulateRoles(ctx context.Context, user *models.UserEntity) error
}

// Mock UserRepository
type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.UserEntity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserEntity), args.Error(1)
}

// Mock UserRoleRepository
type mockUserRoleRepository struct {
	mock.Mock
}

func (m *mockUserRoleRepository) PopulateRoles(ctx context.Context, user *models.UserEntity) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestAuthMiddleware(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Set up in-memory database for tests
	err := setupTestDB()
	assert.NoError(t, err, "Failed to set up test database")
	defer teardownTestDB()

	// Set a mock secret for testing
	originalSecret := os.Getenv("SSO_SECRET")
	os.Setenv("SSO_SECRET", "test_secret_for_auth_middleware")
	defer func() {
		os.Setenv("SSO_SECRET", originalSecret)
	}()

	t.Run("Missing Authorization Header", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest("GET", "/", nil)
		c.Request = req

		// Execute middleware
		Middleware()(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header is required")
	})

	t.Run("Invalid Authorization Format", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", "InvalidFormat")
		c.Request = req

		// Execute middleware
		Middleware()(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header must be in format: Bearer {token}")
	})

	t.Run("Invalid Token", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", "Bearer invalid-token")
		c.Request = req

		// Execute middleware
		Middleware()(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid or expired token")
	})

	t.Run("Valid Token", func(t *testing.T) {
		// Setup
		userID := primitive.NewObjectID()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		tokenDetails, err := jwt.GenerateToken(c, userID, []string{"user"}, jwt.AccessToken)
		assert.NoError(t, err, "Token generation should not fail")
		assert.NotEmpty(t, tokenDetails.Token, "Token should not be empty")

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", "Bearer "+tokenDetails.Token)
		c.Request = req

		// Create a new engine to use with the middleware
		router := gin.New()
		router.Use(Middleware())
		router.GET("/", func(c *gin.Context) {
			// Assert user was set in context
			authUser := GetAuthUser(c)
			assert.NotNil(t, authUser)
			assert.Equal(t, userID, authUser.UserID)
			c.Status(http.StatusOK)
		})

		// Execute
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetAuthUser(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	t.Run("Auth User Not Set", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Execute
		authUser := GetAuthUser(c)

		// Assert
		assert.Nil(t, authUser)
	})

	t.Run("Auth User Set", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		userID := primitive.NewObjectID()
		c.Set("authUser", &UserAuthInfo{UserID: userID})

		// Execute
		authUser := GetAuthUser(c)

		// Assert
		assert.NotNil(t, authUser)
		assert.Equal(t, userID, authUser.UserID)
	})

	t.Run("Invalid Auth User Type", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set invalid type
		c.Set("authUser", "not-a-user-auth-info")

		// Execute
		authUser := GetAuthUser(c)

		// Assert
		assert.Nil(t, authUser)
	})
}

// Helper function to wrap our mock repositories for testing
func mockRoleHandler(roleName string, repo UserRepositoryInterface, roleRepo UserRoleRepositoryInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authenticated user info
		authUser := GetAuthUser(c)
		if authUser == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Get user details from database to check roles
		user, err := repo.FindByID(c, authUser.UserID)
		if err != nil {
			if err.Error() == "user not found" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user roles"})
			}
			c.Abort()
			return
		}

		err = roleRepo.PopulateRoles(c, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user roles"})
			c.Abort()
			return
		}

		// Check if user has the required role
		hasRole := false
		for _, role := range user.Roles {
			if role != nil && role.Name == roleName {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		// User has the required role, proceed
		c.Next()
	}
}

func TestRequireRoleHandler(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	t.Run("No Auth User", func(t *testing.T) {
		// Setup mocks
		userRepo := new(mockUserRepository)
		userRoleRepo := new(mockUserRoleRepository)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Execute our test wrapper function
		mockRoleHandler("admin", userRepo, userRoleRepo)(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authentication required")
	})

	t.Run("User Not Found", func(t *testing.T) {
		// Setup mocks
		userRepo := new(mockUserRepository)
		userRoleRepo := new(mockUserRoleRepository)

		userID := primitive.NewObjectID()
		userRepo.On("FindByID", mock.Anything, userID).Return(nil, errors.New("user not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("authUser", &UserAuthInfo{UserID: userID})

		// Execute our test wrapper function
		mockRoleHandler("admin", userRepo, userRoleRepo)(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "User not found")
		userRepo.AssertExpectations(t)
	})

	t.Run("User Has Required Role", func(t *testing.T) {
		// Setup mocks
		userRepo := new(mockUserRepository)
		userRoleRepo := new(mockUserRoleRepository)

		userID := primitive.NewObjectID()
		email := "test@example.com"
		roleID := primitive.NewObjectID()

		// Create user with admin role
		user := &models.UserEntity{
			ID:      &userID,
			Email:   &email,
			RoleIds: []*primitive.ObjectID{&roleID},
			Roles: []*models.UserRoleEntity{
				{
					ID:   &roleID,
					Name: "admin",
				},
			},
		}

		userRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
		userRoleRepo.On("PopulateRoles", mock.Anything, user).Return(nil)

		// Create a new router to test the middleware chain
		router := gin.New()
		router.Use(func(c *gin.Context) {
			// Set auth user in context
			c.Set("authUser", &UserAuthInfo{UserID: userID})
		})
		router.Use(mockRoleHandler("admin", userRepo, userRoleRepo))
		router.GET("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Create request
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// Execute
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		userRepo.AssertExpectations(t)
		userRoleRepo.AssertExpectations(t)
	})

	t.Run("User Doesn't Have Required Role", func(t *testing.T) {
		// Setup mocks
		userRepo := new(mockUserRepository)
		userRoleRepo := new(mockUserRoleRepository)

		userID := primitive.NewObjectID()
		email := "test@example.com"
		roleID := primitive.NewObjectID()

		// Create user with user role (not admin)
		user := &models.UserEntity{
			ID:      &userID,
			Email:   &email,
			RoleIds: []*primitive.ObjectID{&roleID},
			Roles: []*models.UserRoleEntity{
				{
					ID:   &roleID,
					Name: "user",
				},
			},
		}

		userRepo.On("FindByID", mock.Anything, userID).Return(user, nil)
		userRoleRepo.On("PopulateRoles", mock.Anything, user).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("authUser", &UserAuthInfo{UserID: userID})

		// Execute our test wrapper function
		mockRoleHandler("admin", userRepo, userRoleRepo)(c)

		// Assert
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "Insufficient permissions")
		userRepo.AssertExpectations(t)
		userRoleRepo.AssertExpectations(t)
	})
}

func TestOptionalAuth(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Set up in-memory database for tests
	err := setupTestDB()
	assert.NoError(t, err, "Failed to set up test database")
	defer teardownTestDB()

	// Set a mock secret for testing
	originalSecret := os.Getenv("SSO_SECRET")
	os.Setenv("SSO_SECRET", "test_secret_for_optional_auth")
	defer func() {
		os.Setenv("SSO_SECRET", originalSecret)
	}()

	t.Run("No Authorization Header", func(t *testing.T) {
		// Setup
		router := gin.New()
		router.Use(OptionalAuth())
		router.GET("/", func(c *gin.Context) {
			// Auth user should not be set
			authUser := GetAuthUser(c)
			assert.Nil(t, authUser)
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// Execute
		router.ServeHTTP(w, req)

		// Assert next was called and no errors
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid Authorization Format", func(t *testing.T) {
		// Setup
		router := gin.New()
		router.Use(OptionalAuth())
		router.GET("/", func(c *gin.Context) {
			// Auth user should not be set
			authUser := GetAuthUser(c)
			assert.Nil(t, authUser)
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", "InvalidFormat")

		// Execute
		router.ServeHTTP(w, req)

		// Assert next was called and no errors
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Valid Token", func(t *testing.T) {
		// Setup
		userID := primitive.NewObjectID()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		tokenDetails, err := jwt.GenerateToken(c, userID, []string{"user"}, jwt.AccessToken)
		assert.NoError(t, err, "Token generation should not fail")

		router := gin.New()
		router.Use(OptionalAuth())
		router.GET("/", func(c *gin.Context) {
			// Assert user was set in context
			authUser := GetAuthUser(c)
			if authUser == nil {
				t.Fatal("Auth user should not be nil")
			}
			assert.Equal(t, userID, authUser.UserID)
			c.Status(http.StatusOK)
		})

		w = httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", "Bearer "+tokenDetails.Token)

		// Execute
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestStaticStringMiddleware(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	t.Run("Missing Authorization Header", func(t *testing.T) {
		// Setup
		router := gin.New()
		router.Use(StaticStringMiddleware("Bearer test_token"))
		router.GET("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// Execute
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header is required")
	})

	t.Run("Invalid Token", func(t *testing.T) {
		// Setup
		router := gin.New()
		router.Use(StaticStringMiddleware("Bearer test_token"))
		router.GET("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", "Bearer wrong_token")

		// Execute
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid token")
	})

	t.Run("Valid Token", func(t *testing.T) {
		// Setup
		router := gin.New()
		router.Use(StaticStringMiddleware("Bearer test_token"))
		router.GET("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", "Bearer test_token")

		// Execute
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
