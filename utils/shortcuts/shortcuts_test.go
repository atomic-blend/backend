package shortcuts

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckRequiredEnvVar(t *testing.T) {
	t.Run("should use provided config value", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("The function panicked when it should not: %v", r)
			}
		}()

		CheckRequiredEnvVar("TEST_VAR", "config_value", "default_value")
	})

	t.Run("should use default value when config is empty", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("The function panicked when it should not: %v", r)
			}
		}()

		CheckRequiredEnvVar("TEST_VAR", "", "default_value")
	})

	t.Run("should panic when both config and default are empty", func(t *testing.T) {
		panicked := false
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
			assert.True(t, panicked, "Expected function to panic, but it did not")
		}()

		CheckRequiredEnvVar("TEST_VAR", "", "")
	})
}

func TestFailOnError(t *testing.T) {
	t.Run("should not panic on nil error", func(t *testing.T) {
		defer func() {
			assert.Nil(t, recover())
		}()

		FailOnError(nil, "test message")
	})

	t.Run("should panic on error", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.NotNil(t, r)
		}()

		err := errors.New("test error")
		FailOnError(err, "test message")
	})
}

func TestLogOnError(t *testing.T) {
	t.Run("should return false on nil error", func(t *testing.T) {
		result := LogOnError(nil, "test message")
		assert.False(t, result)
	})

	t.Run("should return true on error", func(t *testing.T) {
		err := errors.New("test error")
		result := LogOnError(err, "test message")
		assert.True(t, result)
	})
}
