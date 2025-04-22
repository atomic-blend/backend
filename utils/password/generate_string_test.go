package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomSalt(t *testing.T) {
	// Test different lengths
	lengths := []int{8, 16, 32, 64}

	for _, length := range lengths {
		salt, err := GenerateRandomString(length)
		assert.NoError(t, err)
		assert.Len(t, salt, length)

		// Verify it's valid hex
		assert.Regexp(t, "^[0-9a-f]+$", salt)
	}

	// Test uniqueness
	salt1, _ := GenerateRandomString(32)
	salt2, _ := GenerateRandomString(32)
	assert.NotEqual(t, salt1, salt2)
}
