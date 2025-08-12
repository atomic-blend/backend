package grpcclientutils

import (
	"fmt"
	"os"
)

// GetServiceBaseURL retrieves the base URL for a given service from environment variables.
func GetServiceBaseURL(service string) (string, error) {
	var baseURL string
	switch service {
	case "auth":
		envValue := os.Getenv("AUTH_GRPC_HOST")
		if envValue != "" {
			baseURL = envValue
		} else {
			baseURL = "http://auth:50051" // Default value if not set
		}
	default:
		return "", fmt.Errorf("unknown service: %s", service)
	}
	return baseURL, nil
}
