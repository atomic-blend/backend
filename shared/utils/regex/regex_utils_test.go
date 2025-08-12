package regexutils

import (
	"testing"
)

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "String with spaces",
			input:    "  Hello World  ",
			expected: "hello world",
		},
		{
			name:     "Mixed case string",
			input:    "MiXeD CaSe",
			expected: "mixed case",
		},
		{
			name:     "String with special characters",
			input:    "hello.world*^$+?()",
			expected: "hello.world*^$+?()",
		},
		{
			name:     "Combined case: mixed case with spaces and special characters",
			input:    "  Hello.World* [Test]  ",
			expected: "hello.world* [test]",
		},
		{
			name:     "Only whitespace",
			input:    "   \t\n  ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeString(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeString() = %q, want %q", result, tt.expected)
			}
		})
	}
}
