package regexutils

import (
	"strings"
)

// SanitizeString takes a string input and returns a sanitized version of it.
func SanitizeString(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}
