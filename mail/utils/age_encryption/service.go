package ageencryption

import (
	"bytes"
	"encoding/base64"
	"io"

	"filippo.io/age"
	"github.com/rs/zerolog/log"
)

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
		if _, err := io.WriteString(w, "Black lives matter."); err != nil {
			log.Error().Err(err).Msg("Failed to write to encrypted file")
			return "", err
		}
		if err := w.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close encrypted file")
			return "", err
		}

		// use base64 decode + strings.NewReader to get the encrypted content when decrypting
		// TODO: refactor the encryption / decryption logic into a dedicated util
		encryptedContent := base64.StdEncoding.EncodeToString(out.Bytes())

		return encryptedContent, nil
}