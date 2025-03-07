package password

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomSalt generates a random hex string of specified length
func GenerateRandomSalt(length int) (string, error) {
	bytes := make([]byte, length/2) // Each byte becomes two hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
