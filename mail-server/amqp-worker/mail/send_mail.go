package mail

import (
	"bytes"
	"context"
	"crypto"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math"
	"net"
	"os"
	"sort"
	"strings"

	"github.com/atomic-blend/backend/mail-server/grpc/client"
	"github.com/atomic-blend/backend/mail-server/utils/amqp"
	"github.com/atomic-blend/backend/mail/models"
	"github.com/emersion/go-message"
	"github.com/emersion/go-msgauth/dkim"
	"github.com/emersion/go-smtp"
	"github.com/rs/zerolog/log"
	amqppackage "github.com/streadway/amqp"
)

const (
	//TODO: make these configurable by env vars
	MaxRetries      = 5
	MaxDelayMillis  = 172800000 // 2 days in ms
	BaseDelayMillis = 10000     // 10 seconds initial backoff
)

var mailClient *client.MailClient

// computeDelay returns delay in ms using exponential backoff with x-day cap
func computeDelay(retryCount int) int {
	delay := BaseDelayMillis * int(math.Pow(2, float64(retryCount-1)))
	if delay > MaxDelayMillis {
		return MaxDelayMillis
	}
	return delay
}

func handleTemporaryFailure(m *amqppackage.Delivery, body []byte, failedReason error, retryCount int, recipientsToRetry []string) {
	log.Info().Msgf("Temporary failure for message: %s, error: %v, retry count: %d", body, failedReason, retryCount)

	// Compute the delay before retrying
	delay := computeDelay(retryCount)
	log.Info().Msgf("Retrying in %d milliseconds", delay)

	//TODO: make the gRPC call to store the retry count, delay in the DB and the reason for failure
	mailClient, err := getMailClient()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create mail client")
		return
	}

	sendEmailID, ok := m.Headers["send_email_id"].(string)
	if !ok {
		log.Error().Msg("send_email_id not found in message headers")
		return
	}

	req := client.CreateRetryStatusRequest(sendEmailID, failedReason.Error(), int32(retryCount))
	_, err = mailClient.UpdateMailStatus(context.Background(), req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update mail status for retry")
		return
	}

	// publish the message into the retry queue with the delay (Dead letter Queue to the original routing key)
	message := make(map[string]interface{})
	err = json.Unmarshal(body, &message)
	if err != nil {
		log.Error().Err(err).Msg("Error unmarshalling message")
		return
	}

	amqp.PublishMessage("mail", "send_retry", message, &amqppackage.Table{
		"retry-count": retryCount,
		"delay":       delay,
		"recipients":  strings.Join(recipientsToRetry, ","),
	})

	// ack the original message
	m.Ack(false)
}

// handlePermanentFailure handles messages that have permanently failed
func handlePermanentFailure(body []byte, err error) {
	log.Info().Msgf("Permanent failure for message: %s, error: %v", body, err)
	//TODO: make the gRPC call to store the failure reason in the DB + status to failed
}

func handleSuccess(message *amqppackage.Delivery) {
	log.Info().Msgf("Message processed successfully: %s", message.Body)
	//TODO: make the gRPC call to store the success in the DB
	mailClient, err := getMailClient()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create mail client")
		return
	}

	var messageWrapper map[string]interface{}
	err = json.Unmarshal(message.Body, &messageWrapper)
	if err != nil {
		log.Error().Err(err).Msg("Error unmarshalling AMQP payload wrapper")
		return
	}

	sendEmailID, ok := messageWrapper["send_mail_id"].(string)
	if !ok {
		log.Error().Msg("send_email_id not found in message headers")
		return
	}

	req := client.CreateSuccessStatusRequest(sendEmailID)
	_, err = mailClient.UpdateMailStatus(context.Background(), req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update mail status for success")
		return
	}

	message.Ack(false)
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
	fromHeader, ok := rawMail.Headers["from"].(string)
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

	// TODO: if the email is a retry, use the list of failed recipients from the message headers
	for _, recipientRaw := range mail.Headers["to"].([]interface{}) {
		recipient, ok := recipientRaw.(string)
		if !ok {
			log.Error().Str("recipient", recipientRaw.(string)).Msg("Failed to convert recipient to string")
			recipientsToRetry = append(recipientsToRetry, recipientRaw.(string))
			continue
		}

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

			from, ok := mail.Headers["from"].(string)
			if !ok {
				log.Error().Str("recipient", recipient).Msg("From header not found")
				recipientsToRetry = append(recipientsToRetry, recipient)
				break
			}

			// Send email via SMTP using mxRecord.Host
			err := sendViaSMTP(mxRecord.Host, from, recipient, signedEmail)
			if err != nil {
				log.Warn().Err(err).Str("mx_host", mxRecord.Host).Str("recipient", recipient).Msg("Failed to send via MX record, trying next one")
				continue
			}

			// Success! Mark as sent and break out of the loop
			sendSuccess = true
			log.Info().Str("mx_host", mxRecord.Host).Str("recipient", recipient).Msg("Email sent successfully via SMTP")
			break
		}

		if !sendSuccess {
			recipientsToRetry = append(recipientsToRetry, recipient)
		}
	}
	if len(recipientsToRetry) > 0 {
		return recipientsToRetry, fmt.Errorf("failed_to_send_to_all_recipients")
	}

	return nil, nil
}

// sendViaSMTP sends an email via SMTP to the specified host
func sendViaSMTP(host string, from string, to string, emailContent string) error {
	// Connect to the remote SMTP server
	c, err := smtp.Dial(host + ":25")
	if err != nil {
		return fmt.Errorf("smtp_connection_failed: %w", err)
	}
	defer c.Quit()

	// Set the sender
	if err := c.Mail(from, nil); err != nil {
		return fmt.Errorf("smtp_mail_command_failed: %w", err)
	}

	// Set the recipient
	if err := c.Rcpt(to, nil); err != nil {
		return fmt.Errorf("smtp_rcpt_command_failed: %w", err)
	}

	// Send the email body
	wc, err := c.Data()
	if err != nil {
		return fmt.Errorf("smtp_data_command_failed: %w", err)
	}

	// Write the email content
	_, err = fmt.Fprintf(wc, "%s", emailContent)
	if err != nil {
		wc.Close()
		return fmt.Errorf("smtp_write_failed: %w", err)
	}

	// Close the data writer
	err = wc.Close()
	if err != nil {
		return fmt.Errorf("smtp_data_close_failed: %w", err)
	}

	log.Info().Str("host", host).Str("from", from).Str("to", to).Msg("Email sent successfully via SMTP")
	return nil
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

func processSendMailMessage(message *amqppackage.Delivery, rawMail models.RawMail) error {
	// lookup the message to check if it's a retry or not
	isRetry := false
	retryCount := 0

	if message.Headers["retry-count"] != nil {
		retryCountInt, ok := message.Headers["retry-count"].(int32)
		if !ok {
			retryCount = 0
		} else {
			retryCount = int(retryCountInt)
		}
	}

	if !isRetry && retryCount > MaxRetries {
		handlePermanentFailure(message.Body, nil) // No error, just a retry limit reached
		return nil
	}

	recipientsToRetry, err := sendEmail(rawMail)
	if err != nil && retryCount < MaxRetries {
		retryCount++
		handleTemporaryFailure(message, message.Body, err, retryCount, recipientsToRetry)
		return nil
	} else if err != nil {
		handlePermanentFailure(message.Body, err)
		return nil
	}

	log.Info().Msg("Email sent successfully, sending success to mail service")
	handleSuccess(message)

	return nil
}

func getMailClient() (*client.MailClient, error) {
	mailClient, err := client.NewMailClient()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create mail client")
		return nil, err
	}
	return mailClient, nil
}
