package rspamdclient

import (
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	// Test default config without environment variables
	config := DefaultConfig()

	if config.BaseURL != "http://localhost:11333" {
		t.Errorf("Expected default BaseURL to be 'http://localhost:11333', got '%s'", config.BaseURL)
	}

	if config.Timeout != 30*time.Second {
		t.Errorf("Expected default Timeout to be 30 seconds, got %v", config.Timeout)
	}

	if config.MaxRetries != 3 {
		t.Errorf("Expected default MaxRetries to be 3, got %d", config.MaxRetries)
	}
}

func TestConfigWithEnvironmentVariables(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("RSPAMD_BASE_URL", "http://rspamd.example.com:11333")
	os.Setenv("RSPAMD_PASSWORD", "testpassword")
	os.Setenv("RSPAMD_TIMEOUT_SECONDS", "60")
	os.Setenv("RSPAMD_MAX_RETRIES", "5")

	// Clean up after test
	defer func() {
		os.Unsetenv("RSPAMD_BASE_URL")
		os.Unsetenv("RSPAMD_PASSWORD")
		os.Unsetenv("RSPAMD_TIMEOUT_SECONDS")
		os.Unsetenv("RSPAMD_MAX_RETRIES")
	}()

	config := DefaultConfig()

	if config.BaseURL != "http://rspamd.example.com:11333" {
		t.Errorf("Expected BaseURL to be 'http://rspamd.example.com:11333', got '%s'", config.BaseURL)
	}

	if config.Password != "testpassword" {
		t.Errorf("Expected Password to be 'testpassword', got '%s'", config.Password)
	}

	if config.Timeout != 60*time.Second {
		t.Errorf("Expected Timeout to be 60 seconds, got %v", config.Timeout)
	}

	if config.MaxRetries != 5 {
		t.Errorf("Expected MaxRetries to be 5, got %d", config.MaxRetries)
	}
}

func TestNewClient(t *testing.T) {
	config := &Config{
		BaseURL:    "http://localhost:11333",
		Password:   "testpass",
		Timeout:    10 * time.Second,
		MaxRetries: 2,
	}

	client := NewClient(config)

	if client.baseURL != "http://localhost:11333" {
		t.Errorf("Expected baseURL to be 'http://localhost:11333', got '%s'", client.baseURL)
	}

	if client.password != "testpass" {
		t.Errorf("Expected password to be 'testpass', got '%s'", client.password)
	}

	if client.httpClient.Timeout != 10*time.Second {
		t.Errorf("Expected timeout to be 10 seconds, got %v", client.httpClient.Timeout)
	}
}

func TestNewClientWithNilConfig(t *testing.T) {
	client := NewClient(nil)

	if client.baseURL != "http://localhost:11333" {
		t.Errorf("Expected default baseURL to be 'http://localhost:11333', got '%s'", client.baseURL)
	}
}

func TestCheckResponseMethods(t *testing.T) {
	// Test spam response
	spamResp := &CheckResponse{
		Action: "reject",
		Score:  15.5,
	}

	if !spamResp.IsSpam() {
		t.Error("Expected IsSpam() to return true for 'reject' action")
	}

	if spamResp.GetScore() != 15.5 {
		t.Errorf("Expected GetScore() to return 15.5, got %f", spamResp.GetScore())
	}

	if spamResp.GetAction() != "reject" {
		t.Errorf("Expected GetAction() to return 'reject', got '%s'", spamResp.GetAction())
	}

	// Test non-spam response
	hamResp := &CheckResponse{
		Action: "no action",
		Score:  2.1,
	}

	if hamResp.IsSpam() {
		t.Error("Expected IsSpam() to return false for 'no action' action")
	}
}

func TestCheckRequestValidation(t *testing.T) {
	client := NewClient(nil)

	// Test with nil request
	_, err := client.CheckMessage(nil)
	if err == nil {
		t.Error("Expected error when passing nil request")
	}

	expectedErr := "request cannot be nil"
	if err.Error() != expectedErr {
		t.Errorf("Expected error message '%s', got '%s'", expectedErr, err.Error())
	}
}
