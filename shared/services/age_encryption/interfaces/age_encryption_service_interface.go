// Package ageinterfaces contains the interfaces for the age encryption service
package ageinterfaces

// AgeEncryptionServiceInterface defines the interface for age encryption operations
type AgeEncryptionServiceInterface interface {
	EncryptString(publicKey string, plaintext string) (string, error)
	EncryptBytes(publicKey string, plaintext []byte) ([]byte, error)
}
