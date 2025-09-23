package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	t.Run("should hash password successfully", func(t *testing.T) {
		password := "test_password123"
		hash, err := HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
	})

	t.Run("should generate different hashes for same password", func(t *testing.T) {
		password := "test_password123"
		hash1, err1 := HashPassword(password)
		hash2, err2 := HashPassword(password)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, hash1, hash2)
	})
}

func TestCheckPassword(t *testing.T) {
	t.Run("should validate correct password", func(t *testing.T) {
		password := "test_password123"
		hash, err := HashPassword(password)
		assert.NoError(t, err)

		isValid := CheckPassword(password, hash)
		assert.True(t, isValid)
	})

	t.Run("should reject incorrect password", func(t *testing.T) {
		password := "test_password123"
		wrongPassword := "wrong_password123"
		hash, err := HashPassword(password)
		assert.NoError(t, err)

		isValid := CheckPassword(wrongPassword, hash)
		assert.False(t, isValid)
	})

	t.Run("should reject invalid hash format", func(t *testing.T) {
		password := "test_password123"
		invalidHash := "invalid_hash_format"

		isValid := CheckPassword(password, invalidHash)
		assert.False(t, isValid)
	})
}
