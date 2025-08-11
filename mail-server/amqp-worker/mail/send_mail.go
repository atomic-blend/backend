package mail

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math"
	"net"
	"os"
	"sort"
	"strings"

	"github.com/atomic-blend/backend/mail/models"
	"github.com/emersion/go-message"
	"github.com/emersion/go-msgauth/dkim"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

const (
	//TODO: make these configurable by env vars
	MaxRetries      = 5
	MaxDelayMillis  = 172800000 // 2 days in ms
	BaseDelayMillis = 10000     // 10 seconds initial backoff
)

// computeDelay returns delay in ms using exponential backoff with x-day cap
func computeDelay(retryCount int) int {
	delay := BaseDelayMillis * int(math.Pow(2, float64(retryCount-1)))
	if delay > MaxDelayMillis {
		return MaxDelayMillis
	}
	return delay
}

func handleTemporaryFailure(body []byte, err error, retryCount int, recipientsToRetry []string) {
	log.Info().Msgf("Temporary failure for message: %s, error: %v, retry count: %d", body, err, retryCount)

	// Compute the delay before retrying
	delay := computeDelay(retryCount)
	log.Info().Msgf("Retrying in %d milliseconds", delay)
	//TODO: make the gRPC call to store the retry count, delay in the DB and the reason for failure

	// TODO: publish the message into the retry queue with the delay (Dead letter Queue to the original routing key)
	// store the list of recipients that needs retry in the message headers
}

// handlePermanentFailure handles messages that have permanently failed
func handlePermanentFailure(body []byte, err error) {
	log.Info().Msgf("Permanent failure for message: %s, error: %v", body, err)
	//TODO: make the gRPC call to store the failure reason in the DB + status to failed
}

func handleSuccess(body []byte) {
	log.Info().Msgf("Message processed successfully: %s", body)
	//TODO: make the gRPC call to store the success in the DB
}

// loadDKIMPrivateKey loads the DKIM private key from the configured path
func loadDKIMPrivateKey() (crypto.Signer, error) {
	keyPath := "/app/dkim_private_key.pem"
	if envPath := os.Getenv("DKIM_PRIVATE_KEY_PATH"); envPath != "" {
		keyPath = envPath
	}

	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read DKIM private key: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from DKIM private key")
	}

	var privateKey crypto.Signer
	switch block.Type {
	case "RSA PRIVATE KEY":
		rsaKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
		}
		privateKey = rsaKey
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKCS8 private key: %w", err)
		}
		signer, ok := key.(crypto.Signer)
		if !ok {
			return nil, fmt.Errorf("parsed key does not implement crypto.Signer")
		}
		privateKey = signer
	default:
		return nil, fmt.Errorf("unsupported private key type: %s", block.Type)
	}

	return privateKey, nil
}

// convertMessageToString converts a message entity to a string for DKIM signing
func convertMessageToString(msg *message.Entity) (string, error) {
	var buf bytes.Buffer
	if err := msg.WriteTo(&buf); err != nil {
		return "", fmt.Errorf("failed to write message to buffer: %w", err)
	}
	return buf.String(), nil
}

// signEmailWithDKIM signs the email with DKIM if a private key is available
func signEmailWithDKIM(rawMail models.RawMail) (string, error) {
	// Convert mail to proper message format once
	msg, err := rawMail.ToMessageEntity()
	if err != nil {
		return "", fmt.Errorf("failed_to_create_message")
	}

	// Check if we can sign the email
	privateKey, err := loadDKIMPrivateKey()
	if err != nil {
		return "", fmt.Errorf("dkim_private_key_load_failed")
	}

	// Get the From domain for DKIM signing
	fromHeader, ok := rawMail.Headers["From"].(string)
	if !ok {
		return "", fmt.Errorf("from_header_missing")
	}

	fromDomain := extractDomain(fromHeader)
	if fromDomain == "" {
		return "", fmt.Errorf("from_domain_extraction_failed")
	}

	// Convert message to string for DKIM signing
	mailString, err := convertMessageToString(msg)
	if err != nil {
		return "", fmt.Errorf("message_to_string_conversion_failed")
	}

	r := strings.NewReader(mailString)

	// Configure DKIM signing options
	selector := "default"
	if envSelector := os.Getenv("DKIM_SELECTOR"); envSelector != "" {
		selector = envSelector
	}

	options := &dkim.SignOptions{
		Domain:   fromDomain,
		Selector: selector,
		Signer:   privateKey,
	}

	var signedBuffer bytes.Buffer
	if err := dkim.Sign(&signedBuffer, r, options); err != nil {
		return "", fmt.Errorf("dkim_signing_failed")
	}

	log.Info().Str("domain", fromDomain).Msg("Email signed successfully with DKIM")
	return signedBuffer.String(), nil
}

func sendEmail(mail models.RawMail) ([]string, error) {
	log.Info().Interface("To", mail.Headers["To"]).Interface("From", mail.Headers["From"]).Msg("Sending email")

	// Sign the email with DKIM first
	signedEmail, err := signEmailWithDKIM(mail)
	if err != nil {
		log.Error().Err(err).Msg("Failed to process email for sending")
		return []string{}, err
	}

	recipientsToRetry := []string{}

	for _, recipient := range mail.Headers["To"].([]string) {
		// Resolve the mail server via MX lookup
		domain := extractDomain(recipient)
		if domain == "" {
			log.Error().Str("recipient", recipient).Msg("Failed to extract domain from recipient")
			recipientsToRetry = append(recipientsToRetry, recipient)
			continue
		}

		mxRecords, err := net.LookupMX(domain)
		if err != nil {
			log.Error().Err(err).Str("domain", domain).Msg("Failed to lookup MX records")
			recipientsToRetry = append(recipientsToRetry, recipient)
			continue
		}

		if len(mxRecords) == 0 {
			log.Error().Str("domain", domain).Msg("No MX records found for domain")
			recipientsToRetry = append(recipientsToRetry, recipient)
			continue
		}

		// Sort the MX records by preference value
		sort.Slice(mxRecords, func(i, j int) bool {
			return mxRecords[i].Pref < mxRecords[j].Pref
		})

		log.Info().Str("recipient", recipient).Str("domain", domain).Str("mx_host", mxRecords[0].Host).Msg("Resolved MX record")

		sendSuccess := false
		for _, mxRecord := range mxRecords {
			log.Info().Str("recipient", recipient).Str("domain", domain).Str("mx_host", mxRecord.Host).Msg("Attempting to send via MX record")

			// TODO: send email via SMTP using mxRecord.Host, if failed, retry with the next MX record
			// For now, just log the attempt
			log.Info().Str("mx_host", mxRecord.Host).Str("signed_email_length", fmt.Sprintf("%d", len(signedEmail))).Msg("Would send signed email via SMTP")

			// TODO: Implement actual SMTP sending here
			// sendSuccess = sendViaSMTP(mxRecord.Host, signedEmail, recipient)
			// if sendSuccess {
			//     break
			// }
		}

		if !sendSuccess {
			recipientsToRetry = append(recipientsToRetry, recipient)
		}
	}

	return recipientsToRetry, nil
}

// extractDomain extracts the domain part from an email address
func extractDomain(email string) string {
	atIndex := -1
	for i, char := range email {
		if char == '@' {
			atIndex = i
			break
		}
	}

	if atIndex == -1 || atIndex == len(email)-1 {
		return ""
	}

	return email[atIndex+1:]
}

func processSendMailMessage(message *amqp.Delivery, rawMail models.RawMail) error {
	// lookup the message to check if it's a retry or not
	isRetry := false
	retryCount, ok := message.Headers["x-retry-count"].(int)
	if !ok {
		log.Error().Msg("Error getting retry count from message headers")
		retryCount = 0
		isRetry = true
	} else {
		log.Debug().Msgf("Retry count from message headers: %d", retryCount)
	}

	if !isRetry && retryCount > MaxRetries {
		handlePermanentFailure(message.Body, nil) // No error, just a retry limit reached
		return nil
	}

	recipientsToRetry, err := sendEmail(rawMail)
	if err != nil && retryCount < MaxRetries {
		handleTemporaryFailure(message.Body, err, retryCount, recipientsToRetry)
		return nil
	} else if err != nil {
		handlePermanentFailure(message.Body, err)
		return nil
	}

	log.Info().Msg("Email sent successfully, sending success to mail service")
	handleSuccess(message.Body)

	log.Info().Msg("Acknowledging message")

	// If email sent successfully, acknowledge the message
	if err := message.Ack(false); err != nil {
		log.Error().Err(err).Msg("Failed to acknowledge message")
	}

	log.Info().Msg("Email sent successfully")

	return nil
}
