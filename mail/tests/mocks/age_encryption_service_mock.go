package mocks

import (
	ageencryptionservice "github.com/atomic-blend/backend/mail/services/age_encryption"
	"github.com/stretchr/testify/mock"
)

// MockAgeEncryptionService provides a mock implementation of age encryption service
type MockAgeEncryptionService struct {
	mock.Mock
}

// Ensure MockAgeEncryptionService implements the interface
var _ ageencryptionservice.AgeEncryptionServiceInterface = (*MockAgeEncryptionService)(nil)

// EncryptString encrypts a string using the age encryption library
func (m *MockAgeEncryptionService) EncryptString(publicKey string, plaintext string) (string, error) {
	args := m.Called(publicKey, plaintext)
	return args.String(0), args.Error(1)
}

// EncryptBytes encrypts a byte array using the age encryption library
func (m *MockAgeEncryptionService) EncryptBytes(publicKey string, plaintext []byte) ([]byte, error) {
	args := m.Called(publicKey, plaintext)
	return args.Get(0).([]byte), args.Error(1)
}
