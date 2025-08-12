package rspamdclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// DefaultConfig returns default configuration with environment variable support
func DefaultConfig() *Config {
	config := &Config{
		BaseURL:    "http://localhost:11333",
		Password:   "",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
	}

	// Override with environment variables if set
	if baseURL := os.Getenv("RSPAMD_BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}

	if password := os.Getenv("RSPAMD_PASSWORD"); password != "" {
		config.Password = password
	}

	if timeoutStr := os.Getenv("RSPAMD_TIMEOUT_SECONDS"); timeoutStr != "" {
		if timeout, err := strconv.Atoi(timeoutStr); err == nil {
			config.Timeout = time.Duration(timeout) * time.Second
		}
	}

	if retriesStr := os.Getenv("RSPAMD_MAX_RETRIES"); retriesStr != "" {
		if retries, err := strconv.Atoi(retriesStr); err == nil {
			config.MaxRetries = retries
		}
	}

	return config
}

// NewClient creates a new Rspamd client with the given configuration
func NewClient(config *Config) *Client {
	if config == nil {
		config = DefaultConfig()
	}

	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    strings.TrimSuffix(config.BaseURL, "/"),
		password:   config.Password,
	}
}

// CheckMessage sends a message to Rspamd for spam checking
func (c *Client) CheckMessage(req *CheckRequest) (*CheckResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.baseURL+"/checkv2", bytes.NewReader(req.Message))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set required headers
	httpReq.Header.Set("Content-Type", "application/octet-stream")
	httpReq.Header.Set("Content-Length", strconv.Itoa(len(req.Message)))

	// Set optional headers based on request
	if req.DeliverTo != "" {
		httpReq.Header.Set("Deliver-To", req.DeliverTo)
	}
	if req.IP != "" {
		httpReq.Header.Set("IP", req.IP)
	}
	if req.Helo != "" {
		httpReq.Header.Set("Helo", req.Helo)
	}
	if req.Hostname != "" {
		httpReq.Header.Set("Hostname", req.Hostname)
	}
	if len(req.Flags) > 0 {
		httpReq.Header.Set("Flags", strings.Join(req.Flags, ","))
	}
	if req.From != "" {
		httpReq.Header.Set("From", req.From)
	}
	if req.QueueID != "" {
		httpReq.Header.Set("Queue-Id", req.QueueID)
	}
	if req.Raw {
		httpReq.Header.Set("Raw", "yes")
	}
	for _, rcpt := range req.Rcpt {
		httpReq.Header.Add("Rcpt", rcpt)
	}
	if req.Pass != "" {
		httpReq.Header.Set("Pass", req.Pass)
	}
	if req.Subject != "" {
		httpReq.Header.Set("Subject", req.Subject)
	}
	if req.User != "" {
		httpReq.Header.Set("User", req.User)
	}
	if c.password != "" {
		httpReq.Header.Set("Password", c.password)
	}

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var checkResp CheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&checkResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &checkResp, nil
}

// Ping sends a ping request to check if Rspamd is available
func (c *Client) Ping() error {
	resp, err := c.httpClient.Get(c.baseURL + "/ping")
	if err != nil {
		return fmt.Errorf("failed to ping Rspamd: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected ping status code: %d", resp.StatusCode)
	}

	return nil
}

// IsSpam checks if the response indicates spam
func (r *CheckResponse) IsSpam() bool {
	return r.Action == "reject" || r.Action == "soft reject" || r.Action == "add header"
}

// GetScore returns the spam score
func (r *CheckResponse) GetScore() float64 {
	return r.Score
}

// GetAction returns the action taken by Rspamd
func (r *CheckResponse) GetAction() string {
	return r.Action
}
