package mailsender

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"

	"github.com/atomic-blend/backend/mail/models"
	"github.com/emersion/go-message"
	"github.com/emersion/go-msgauth/dkim"
	"github.com/emersion/go-smtp"
	"github.com/rs/zerolog/log"
)

func SendEmail(mail models.RawMail, recipients []any) ([]string, error) {
	log.Info().Interface("To", mail.Headers["To"]).Interface("From", mail.Headers["From"]).Msg("Sending email")

	// Sign the email with DKIM first
	signedEmail, err := signEmailWithDKIM(mail)
	if err != nil {
		log.Error().Err(err).Msg("Failed to process email for sending")
		return []string{}, err
	}

	recipientsToRetry := []string{}

	var recipientsToSend []any
	if len(recipients) > 0 {
		recipientsToSend = recipients
	} else {
		recipientsToSend = mail.Headers["To"].([]any)
	}

	// TODO: if the email is a retry, use the list of failed recipients from the message headers
	for _, recipientRaw := range recipientsToSend {
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

		from, ok := mail.Headers["From"].(string)
		if !ok {
			log.Error().Str("recipient", recipient).Msg("From header not found")
			recipientsToRetry = append(recipientsToRetry, recipient)
			continue
		}

		sendSuccess := false
		// Try all MX records with TLS first on port 25
		for _, mxRecord := range mxRecords {
			log.Info().Str("mx_host", mxRecord.Host).Int("port", 25).Bool("useTLS", true).Str("recipient", recipient).Msg("Attempting to send via SMTP with TLS on port 25")

			err := sendViaSMTP(mxRecord.Host, 25, from, recipient, signedEmail, true)
			if err != nil {
				log.Warn().Err(err).Str("mx_host", mxRecord.Host).Int("port", 25).Bool("useTLS", true).Str("recipient", recipient).Msg("Failed to send via SMTP with TLS on port 25, trying next MX record")
				continue
			}

			// Success! Mark as sent and break out of the loop
			sendSuccess = true
			log.Info().Str("mx_host", mxRecord.Host).Int("port", 25).Bool("useTLS", true).Str("recipient", recipient).Msg("Email sent successfully via SMTP with TLS on port 25")
			break
		}

		// If TLS didn't work, try all MX records without TLS on port 25
		if !sendSuccess {
			log.Info().Str("recipient", recipient).Msg("TLS failed for all MX records, trying without TLS on port 25")
			for _, mxRecord := range mxRecords {
				log.Info().Str("mx_host", mxRecord.Host).Int("port", 25).Bool("useTLS", false).Str("recipient", recipient).Msg("Attempting to send via SMTP without TLS on port 25")

				err := sendViaSMTP(mxRecord.Host, 25, from, recipient, signedEmail, false)
				if err != nil {
					log.Warn().Err(err).Str("mx_host", mxRecord.Host).Int("port", 25).Bool("useTLS", false).Str("recipient", recipient).Msg("Failed to send via SMTP without TLS on port 25, trying next MX record")
					continue
				}

				// Success! Mark as sent and break out of the loop
				sendSuccess = true
				log.Info().Str("mx_host", mxRecord.Host).Int("port", 25).Bool("useTLS", false).Str("recipient", recipient).Msg("Email sent successfully via SMTP without TLS on port 25")
				break
			}
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

// sendViaSMTP sends an email via SMTP to the specified host and port
func sendViaSMTP(host string, port int, from string, to string, emailContent string, useTLS bool) error {
	var c *smtp.Client
	var err error

	if useTLS {
		c, err = smtp.DialStartTLS(fmt.Sprintf("%s:%d", host, port), nil)
		if err != nil {
			return fmt.Errorf("smtp_connection_failed: %w", err)
		}
	} else {
		c, err = smtp.Dial(fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			return fmt.Errorf("smtp_connection_failed: %w", err)
		}
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

	log.Info().Str("host", host).Int("port", port).Str("from", from).Str("to", to).Msg("Email sent successfully via SMTP")
	return nil
}

// convertMessageToString converts a message entity to a string for DKIM signing
func convertMessageToString(msg *message.Entity) (string, error) {
	var buf bytes.Buffer
	if err := msg.WriteTo(&buf); err != nil {
		return "", fmt.Errorf("failed to write message to buffer: %w", err)
	}
	return buf.String(), nil
}
