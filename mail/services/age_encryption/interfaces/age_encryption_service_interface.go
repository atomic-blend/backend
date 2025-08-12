package interfaces

// AgeEncryptionServiceInterface defines the interface for age encryption operations
type AgeEncryptionServiceInterface interface {
	EncryptString(publicKey string, plaintext string) (string, error)
	EncryptBytes(publicKey string, plaintext []byte) ([]byte, error)
}
