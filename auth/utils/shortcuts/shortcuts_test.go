package shortcuts

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func TestContainsDateTime(t *testing.T) {
	// Create test data
	now := time.Now()
	oneHourLater := now.Add(time.Hour)
	twoHoursLater := now.Add(2 * time.Hour)

	// Create slice with two times
	dateTimeSlice := []primitive.DateTime{
		primitive.NewDateTimeFromTime(now),
		primitive.NewDateTimeFromTime(twoHoursLater),
	}

	tests := []struct {
		name     string
		slice    []primitive.DateTime
		item     time.Time
		expected bool
	}{
		{
			name:     "should find existing time",
			slice:    dateTimeSlice,
			item:     now,
			expected: true,
		},
		{
			name:     "should find another existing time",
			slice:    dateTimeSlice,
			item:     twoHoursLater,
			expected: true,
		},
		{
			name:     "should not find non-existing time",
			slice:    dateTimeSlice,
			item:     oneHourLater,
			expected: false,
		},
		{
			name:     "should handle empty slice",
			slice:    []primitive.DateTime{},
			item:     now,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsDateTime(tt.slice, tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContainsInt(t *testing.T) {
	// Create test data
	intSlice := []int{1, 3, 5, 7, 9}

	tests := []struct {
		name     string
		slice    []int
		item     int
		expected bool
	}{
		{
			name:     "should find existing int at beginning",
			slice:    intSlice,
			item:     1,
			expected: true,
		},
		{
			name:     "should find existing int in middle",
			slice:    intSlice,
			item:     5,
			expected: true,
		},
		{
			name:     "should find existing int at end",
			slice:    intSlice,
			item:     9,
			expected: true,
		},
		{
			name:     "should not find non-existing int",
			slice:    intSlice,
			item:     2,
			expected: false,
		},
		{
			name:     "should handle empty slice",
			slice:    []int{},
			item:     1,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsInt(tt.slice, tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}
