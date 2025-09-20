package config

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestController_AvailableAccountDomain(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		envValue        string
		expectedStatus  int
		expectedDomains []string
	}{
		{
			name:            "Success with multiple domains",
			envValue:        "example.com,test.com,domain.org",
			expectedStatus:  http.StatusOK,
			expectedDomains: []string{"example.com", "test.com", "domain.org"},
		},
		{
			name:            "Success with single domain",
			envValue:        "example.com",
			expectedStatus:  http.StatusOK,
			expectedDomains: []string{"example.com"},
		},
		{
			name:            "Success with empty environment variable",
			envValue:        "",
			expectedStatus:  http.StatusOK,
			expectedDomains: []string{},
		},
		{
			name:            "Success with domains containing spaces",
			envValue:        "example.com, test.com , domain.org",
			expectedStatus:  http.StatusOK,
			expectedDomains: []string{"example.com", " test.com ", " domain.org"},
		},
		{
			name:            "Success with empty string in comma-separated list",
			envValue:        "example.com,,test.com",
			expectedStatus:  http.StatusOK,
			expectedDomains: []string{"example.com", "", "test.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original environment variable value
			originalValue := os.Getenv("ACCOUNT_DOMAINS")
			defer func() {
				// Restore original value or unset if it wasn't set
				if originalValue == "" {
					os.Unsetenv("ACCOUNT_DOMAINS")
				} else {
					os.Setenv("ACCOUNT_DOMAINS", originalValue)
				}
			}()

			// Set test environment variable
			if tt.envValue == "" {
				os.Unsetenv("ACCOUNT_DOMAINS")
			} else {
				os.Setenv("ACCOUNT_DOMAINS", tt.envValue)
			}

			// Create controller
			controller := NewConfigController()

			// Create test router
			router := gin.New()
			router.GET("/config", controller.AvailableAccountDomain)

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

			// Check that domains key exists
			assert.Contains(t, response, "domains")

			// Extract domains array and convert to []string
			domainsInterface, ok := response["domains"].([]interface{})
			assert.True(t, ok, "domains should be an array")

			domains := make([]string, len(domainsInterface))
			for i, domain := range domainsInterface {
				domains[i] = domain.(string)
			}

			// Check the response values
			assert.Equal(t, tt.expectedDomains, domains)
		})
	}
}

func TestController_AvailableAccountDomain_EnvironmentVariableNotSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Store original environment variable value
	originalValue := os.Getenv("ACCOUNT_DOMAINS")
	defer func() {
		// Restore original value or unset if it wasn't set
		if originalValue == "" {
			os.Unsetenv("ACCOUNT_DOMAINS")
		} else {
			os.Setenv("ACCOUNT_DOMAINS", originalValue)
		}
	}()

	// Ensure environment variable is not set
	os.Unsetenv("ACCOUNT_DOMAINS")

	// Create controller
	controller := NewConfigController()

	// Create test router
	router := gin.New()
	router.GET("/config", controller.AvailableAccountDomain)

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

	// Check that domains key exists
	assert.Contains(t, response, "domains")

	// Extract domains array and convert to []string
	domainsInterface, ok := response["domains"].([]interface{})
	assert.True(t, ok, "domains should be an array")

	domains := make([]string, len(domainsInterface))
	for i, domain := range domainsInterface {
		domains[i] = domain.(string)
	}

	// Check the response values
	expectedDomains := []string{}
	assert.Equal(t, expectedDomains, domains)
}
