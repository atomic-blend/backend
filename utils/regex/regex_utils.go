package regexutils

import (
	"regexp"
	"strings"
)

func SanitizeString(input string) string {
	return regexp.QuoteMeta(strings.TrimSpace(strings.ToLower(input)))
}
