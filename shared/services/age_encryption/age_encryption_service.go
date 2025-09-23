// Package ageencryptionservice contains the age encryption service
package ageencryptionservice

import (
	interfaces "github.com/atomic-blend/backend/shared/services/age_encryption/interfaces"
)

// Wrapper wraps the existing age encryption functionality
type Wrapper struct{}

// NewAgeEncryptionService creates a new age encryption service wrapper
func NewAgeEncryptionService() interfaces.AgeEncryptionServiceInterface {
	return &Wrapper{}
}

// EncryptString encrypts a string using the age encryption library
func (a *Wrapper) EncryptString(publicKey string, plaintext string) (string, error) {
	return EncryptString(publicKey, plaintext)
}

// EncryptBytes encrypts a byte array using the age encryption library
func (a *Wrapper) EncryptBytes(publicKey string, plaintext []byte) ([]byte, error) {
	return EncryptBytes(publicKey, plaintext)
}
