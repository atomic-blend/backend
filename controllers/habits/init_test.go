package habits

import (
	"atomic_blend_api/models"
	"testing"
)

func TestMain(m *testing.M) {
	// Register validators before running tests
	models.RegisterValidators()

	// Run all tests
	m.Run()
}
