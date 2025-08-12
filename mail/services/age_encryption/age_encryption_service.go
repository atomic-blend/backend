package ageencryptionservice

import (
	interfaces "github.com/atomic-blend/backend/mail/services/age_encryption/interfaces"
	ageencryption "github.com/atomic-blend/backend/mail/utils/age_encryption"
)

// AgeEncryptionServiceWrapper wraps the existing age encryption functionality
type AgeEncryptionServiceWrapper struct{}

// NewAgeEncryptionService creates a new age encryption service wrapper
func NewAgeEncryptionService() interfaces.AgeEncryptionServiceInterface {
	return &AgeEncryptionServiceWrapper{}
}

// EncryptString encrypts a string using the age encryption library
func (a *AgeEncryptionServiceWrapper) EncryptString(publicKey string, plaintext string) (string, error) {
	return ageencryption.EncryptString(publicKey, plaintext)
}

// EncryptBytes encrypts a byte array using the age encryption library
func (a *AgeEncryptionServiceWrapper) EncryptBytes(publicKey string, plaintext []byte) ([]byte, error) {
	return ageencryption.EncryptBytes(publicKey, plaintext)
}
