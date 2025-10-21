package config

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/atomic-blend/backend/auth/tests/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestController_GetConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name                   string
		envValue               string
		maxUsersEnv            string
		userCount              int64
		expectedStatus         int
		expectedDomains        []string
		expectedRemainingSpots int64
	}{
		{
			name:                   "Success with multiple domains and user count",
			envValue:               "example.com,test.com,domain.org",
			maxUsersEnv:            "100",
			userCount:              25,
			expectedStatus:         http.StatusOK,
			expectedDomains:        []string{"example.com", "test.com", "domain.org"},
			expectedRemainingSpots: 75, // 100 - 25
		},
		{
			name:                   "Success with single domain and default max users",
			envValue:               "example.com",
			maxUsersEnv:            "",
			userCount:              0,
			expectedStatus:         http.StatusOK,
			expectedDomains:        []string{"example.com"},
			expectedRemainingSpots: 1, // default 1 - 0
		},
		{
			name:                   "Success with empty environment variable",
			envValue:               "",
			maxUsersEnv:            "50",
			userCount:              10,
			expectedStatus:         http.StatusOK,
			expectedDomains:        []string{},
			expectedRemainingSpots: 40, // 50 - 10
		},
		{
			name:                   "Success with domains containing spaces",
			envValue:               "example.com, test.com , domain.org",
			maxUsersEnv:            "200",
			userCount:              150,
			expectedStatus:         http.StatusOK,
			expectedDomains:        []string{"example.com", " test.com ", " domain.org"},
			expectedRemainingSpots: 50, // 200 - 150
		},
		{
			name:                   "Success with empty string in comma-separated list",
			envValue:               "example.com,,test.com",
			maxUsersEnv:            "10",
			userCount:              5,
			expectedStatus:         http.StatusOK,
			expectedDomains:        []string{"example.com", "", "test.com"},
			expectedRemainingSpots: 5, // 10 - 5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original environment variable values
			originalDomainsValue := os.Getenv("ACCOUNT_DOMAINS")
			originalMaxUsersValue := os.Getenv("AUTH_MAX_NB_USER")
			defer func() {
				// Restore original values or unset if they weren't set
				if originalDomainsValue == "" {
					os.Unsetenv("ACCOUNT_DOMAINS")
				} else {
					os.Setenv("ACCOUNT_DOMAINS", originalDomainsValue)
				}
				if originalMaxUsersValue == "" {
					os.Unsetenv("AUTH_MAX_NB_USER")
				} else {
					os.Setenv("AUTH_MAX_NB_USER", originalMaxUsersValue)
				}
			}()

			// Set test environment variables
			if tt.envValue == "" {
				os.Unsetenv("ACCOUNT_DOMAINS")
			} else {
				os.Setenv("ACCOUNT_DOMAINS", tt.envValue)
			}
			if tt.maxUsersEnv == "" {
				os.Unsetenv("AUTH_MAX_NB_USER")
			} else {
				os.Setenv("AUTH_MAX_NB_USER", tt.maxUsersEnv)
			}

			// Create mock user repository
			mockUserRepo := &mocks.MockUserRepository{}
			mockUserRepo.On("Count", mock.Anything).Return(tt.userCount, nil)

			// Create mock waiting list repository
			mockWaitingListRepo := &mocks.MockWaitingListRepository{}
			mockWaitingListRepo.On("CountWithCode", mock.Anything).Return(int64(0), nil)

			// Create controller
			controller := NewConfigController(mockUserRepo, mockWaitingListRepo)

			// Create test router
			router := gin.New()
			router.GET("/config", controller.GetConfig)

			// Create request
			req, _ := http.NewRequest("GET", "/config", nil)
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response body
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Check that domains and remainingSpots keys exist
			assert.Contains(t, response, "domains")
			assert.Contains(t, response, "remainingSpots")

			// Extract domains array and convert to []string
			domainsInterface, ok := response["domains"].([]interface{})
			assert.True(t, ok, "domains should be an array")

			domains := make([]string, len(domainsInterface))
			for i, domain := range domainsInterface {
				domains[i] = domain.(string)
			}

			// Extract remainingSpots value
			remainingSpotsFloat, ok := response["remainingSpots"].(float64)
			assert.True(t, ok, "remainingSpots should be a number")
			remainingSpots := int64(remainingSpotsFloat)

			// Check the response values
			assert.Equal(t, tt.expectedDomains, domains)
			assert.Equal(t, tt.expectedRemainingSpots, remainingSpots)

			// Verify mock was called
			mockUserRepo.AssertExpectations(t)
			mockWaitingListRepo.AssertExpectations(t)
		})
	}
}

func TestController_GetConfig_EnvironmentVariableNotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Store original environment variable values
	originalDomainsValue := os.Getenv("ACCOUNT_DOMAINS")
	originalMaxUsersValue := os.Getenv("AUTH_MAX_NB_USER")
	defer func() {
		// Restore original values or unset if they weren't set
		if originalDomainsValue == "" {
			os.Unsetenv("ACCOUNT_DOMAINS")
		} else {
			os.Setenv("ACCOUNT_DOMAINS", originalDomainsValue)
		}
		if originalMaxUsersValue == "" {
			os.Unsetenv("AUTH_MAX_NB_USER")
		} else {
			os.Setenv("AUTH_MAX_NB_USER", originalMaxUsersValue)
		}
	}()

	// Ensure environment variables are not set
	os.Unsetenv("ACCOUNT_DOMAINS")
	os.Unsetenv("AUTH_MAX_NB_USER")

	// Create mock user repository
	mockUserRepo := &mocks.MockUserRepository{}
	mockUserRepo.On("Count", mock.Anything).Return(int64(0), nil)

	// Create mock waiting list repository
	mockWaitingListRepo := &mocks.MockWaitingListRepository{}
	mockWaitingListRepo.On("CountWithCode", mock.Anything).Return(int64(0), nil)

	// Create controller
	controller := NewConfigController(mockUserRepo, mockWaitingListRepo)

	// Create test router
	router := gin.New()
	router.GET("/config", controller.GetConfig)

	// Create request
	req, _ := http.NewRequest("GET", "/config", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response body
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check that domains and remainingSpots keys exist
	assert.Contains(t, response, "domains")
	assert.Contains(t, response, "remainingSpots")

	// Extract domains array and convert to []string
	domainsInterface, ok := response["domains"].([]interface{})
	assert.True(t, ok, "domains should be an array")

	domains := make([]string, len(domainsInterface))
	for i, domain := range domainsInterface {
		domains[i] = domain.(string)
	}

	// Extract remainingSpots value
	remainingSpotsFloat, ok := response["remainingSpots"].(float64)
	assert.True(t, ok, "remainingSpots should be a number")
	remainingSpots := int64(remainingSpotsFloat)

	// Check the response values
	expectedDomains := []string{}
	expectedRemainingSpots := int64(1) // default 1 - 0
	assert.Equal(t, expectedDomains, domains)
	assert.Equal(t, expectedRemainingSpots, remainingSpots)

	// Verify mock was called
	mockUserRepo.AssertExpectations(t)
	mockWaitingListRepo.AssertExpectations(t)
}

func TestController_GetConfig_InvalidMaxUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Store original environment variable values
	originalDomainsValue := os.Getenv("ACCOUNT_DOMAINS")
	originalMaxUsersValue := os.Getenv("AUTH_MAX_NB_USER")
	defer func() {
		// Restore original values or unset if they weren't set
		if originalDomainsValue == "" {
			os.Unsetenv("ACCOUNT_DOMAINS")
		} else {
			os.Setenv("ACCOUNT_DOMAINS", originalDomainsValue)
		}
		if originalMaxUsersValue == "" {
			os.Unsetenv("AUTH_MAX_NB_USER")
		} else {
			os.Setenv("AUTH_MAX_NB_USER", originalMaxUsersValue)
		}
	}()

	// Set invalid max users environment variable
	os.Setenv("ACCOUNT_DOMAINS", "example.com")
	os.Setenv("AUTH_MAX_NB_USER", "invalid")

	// Create mock user repository
	mockUserRepo := &mocks.MockUserRepository{}
	mockUserRepo.On("Count", mock.Anything).Return(int64(0), nil)

	// Create mock waiting list repository
	mockWaitingListRepo := &mocks.MockWaitingListRepository{}
	mockWaitingListRepo.On("CountWithCode", mock.Anything).Return(int64(0), nil)

	// Create controller
	controller := NewConfigController(mockUserRepo, mockWaitingListRepo)

	// Create test router
	router := gin.New()
	router.GET("/config", controller.GetConfig)

	// Create request
	req, _ := http.NewRequest("GET", "/config", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check status code - should return 500 for invalid max users
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Parse response body
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check that error message exists
	assert.Contains(t, response, "error")
	assert.Equal(t, "Failed to parse AUTH_MAX_NB_USER", response["error"])
}

func TestController_GetConfig_UserCountError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Store original environment variable values
	originalDomainsValue := os.Getenv("ACCOUNT_DOMAINS")
	originalMaxUsersValue := os.Getenv("AUTH_MAX_NB_USER")
	defer func() {
		// Restore original values or unset if they weren't set
		if originalDomainsValue == "" {
			os.Unsetenv("ACCOUNT_DOMAINS")
		} else {
			os.Setenv("ACCOUNT_DOMAINS", originalDomainsValue)
		}
		if originalMaxUsersValue == "" {
			os.Unsetenv("AUTH_MAX_NB_USER")
		} else {
			os.Setenv("AUTH_MAX_NB_USER", originalMaxUsersValue)
		}
	}()

	// Set environment variables
	os.Setenv("ACCOUNT_DOMAINS", "example.com")
	os.Setenv("AUTH_MAX_NB_USER", "100")

	// Create mock user repository that returns an error
	mockUserRepo := &mocks.MockUserRepository{}
	mockUserRepo.On("Count", mock.Anything).Return(int64(0), assert.AnError)

	// Create mock waiting list repository
	mockWaitingListRepo := &mocks.MockWaitingListRepository{}
	mockWaitingListRepo.On("CountWithCode", mock.Anything).Return(int64(0), nil)

	// Create controller
	controller := NewConfigController(mockUserRepo, mockWaitingListRepo)

	// Create test router
	router := gin.New()
	router.GET("/config", controller.GetConfig)

	// Create request
	req, _ := http.NewRequest("GET", "/config", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check status code - should return 500 for user count error
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Parse response body
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check that error message exists
	assert.Contains(t, response, "error")
	assert.Equal(t, "Failed to get current user count", response["error"])

	// Verify mock was called
	mockUserRepo.AssertExpectations(t)
}

func TestController_GetConfig_WithWaitingListCodes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Store original environment variable values
	originalDomainsValue := os.Getenv("ACCOUNT_DOMAINS")
	originalMaxUsersValue := os.Getenv("AUTH_MAX_NB_USER")
	defer func() {
		// Restore original values or unset if they weren't set
		if originalDomainsValue == "" {
			os.Unsetenv("ACCOUNT_DOMAINS")
		} else {
			os.Setenv("ACCOUNT_DOMAINS", originalDomainsValue)
		}
		if originalMaxUsersValue == "" {
			os.Unsetenv("AUTH_MAX_NB_USER")
		} else {
			os.Setenv("AUTH_MAX_NB_USER", originalMaxUsersValue)
		}
	}()

	// Set environment variables
	os.Setenv("ACCOUNT_DOMAINS", "example.com")
	os.Setenv("AUTH_MAX_NB_USER", "100")

	// Create mock user repository
	mockUserRepo := &mocks.MockUserRepository{}
	mockUserRepo.On("Count", mock.Anything).Return(int64(25), nil)

	// Create mock waiting list repository with 10 users having codes
	mockWaitingListRepo := &mocks.MockWaitingListRepository{}
	mockWaitingListRepo.On("CountWithCode", mock.Anything).Return(int64(10), nil)

	// Create controller
	controller := NewConfigController(mockUserRepo, mockWaitingListRepo)

	// Create test router
	router := gin.New()
	router.GET("/config", controller.GetConfig)

	// Create request
	req, _ := http.NewRequest("GET", "/config", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response body
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Extract remainingSpots value
	remainingSpotsFloat, ok := response["remainingSpots"].(float64)
	assert.True(t, ok, "remainingSpots should be a number")
	remainingSpots := int64(remainingSpotsFloat)

	// Check that remaining spots accounts for both users and waiting list with codes
	// 100 max - 25 users - 10 waiting list with codes = 65
	expectedRemainingSpots := int64(65)
	assert.Equal(t, expectedRemainingSpots, remainingSpots)

	// Verify mocks were called
	mockUserRepo.AssertExpectations(t)
	mockWaitingListRepo.AssertExpectations(t)
}
