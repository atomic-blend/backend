package ageencryption

import (
	"bytes"
	"encoding/base64"
	"io"

	"filippo.io/age"
	"github.com/rs/zerolog/log"
)

// EncryptString encrypts a string using the age encryption library
func EncryptString(publicKey string, plaintext string) (string, error) {
	recipient, err := age.ParseX25519Recipient(publicKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse public key")
		return "", err
	}

	out := &bytes.Buffer{}

	w, err := age.Encrypt(out, recipient)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create encrypted file")
		return "", err
	}
	if _, err := io.WriteString(w, plaintext); err != nil {
		log.Error().Err(err).Msg("Failed to write to encrypted file")
		return "", err
	}
	if err := w.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close encrypted file")
		return "", err
	}

	// use base64 decode + strings.NewReader to get the encrypted content when decrypting
	encryptedContent := base64.StdEncoding.EncodeToString(out.Bytes())

	return encryptedContent, nil
}

// EncryptBytes encrypts a byte array using the age encryption library
func EncryptBytes(publicKey string, plaintext []byte) ([]byte, error) {
	recipient, err := age.ParseX25519Recipient(publicKey)
	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}

	w, err := age.Encrypt(out, recipient)
	if err != nil {
		return nil, err
	}

	if _, err := w.Write(plaintext); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
