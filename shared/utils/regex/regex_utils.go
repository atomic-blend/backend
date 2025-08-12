package regexutils

import (
	"regexp"
	"strings"
)

// SanitizeString takes a string input and returns a sanitized version of it.
func SanitizeString(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}