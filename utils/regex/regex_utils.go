package regexutils

import (
	"regexp"
	"strings"
)

// SanitizeString takes a string input and returns a sanitized version of it.
func SanitizeString(input string) string {
	return regexp.QuoteMeta(strings.TrimSpace(strings.ToLower(input)))
}
