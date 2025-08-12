package grpcclientutils

import (
	"fmt"
	"os"
)

// GetServiceBaseURL retrieves the base URL for a given service from environment variables.
func GetServiceBaseURL(service string) (string, error) {
	var baseURL string
	switch service {
	case "mail":
		envValue := os.Getenv("MAIL_SERVICE_BASE_URL")
		if envValue != "" {
			baseURL = envValue
		} else {
			baseURL = "http://mail:50051" // Default value if not set
		}
	default:
		return "", fmt.Errorf("unknown service: %s", service)
	}
	return baseURL, nil
}
