package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/atomic-blend/backend/mail-server/utils/amqp"
	mailsender "github.com/atomic-blend/backend/mail-server/utils/mail-sender"
	"github.com/atomic-blend/backend/mail/models"
	mailclient "github.com/atomic-blend/backend/shared/grpc/mail"
	"github.com/rs/zerolog/log"
	amqppackage "github.com/streadway/amqp"
)

var (
	// MaxRetries is the maximum number of retries for a mail (configurable via MAX_RETRIES env var)
	MaxRetries = getEnvAsInt("MAX_RETRIES", 5)
	// MaxDelayMillis is the maximum delay for a mail (configurable via MAX_DELAY_MILLIS env var)
	MaxDelayMillis = getEnvAsInt("MAX_DELAY_MILLIS", 172800000) // 2 days in ms
	// BaseDelayMillis is the base delay for a mail (configurable via BASE_DELAY_MILLIS env var)
	BaseDelayMillis = getEnvAsInt("BASE_DELAY_MILLIS", 10000) // 10 seconds initial backoff
)

// getEnvAsInt retrieves an environment variable as an integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Warn().Str("key", key).Str("value", value).Int("default", defaultValue).Msg("Invalid environment variable value, using default")
	}
	return defaultValue
}

// computeDelay returns delay in ms using exponential backoff with x-day cap
func computeDelay(retryCount int) int {
	delay := BaseDelayMillis * int(math.Pow(2, float64(retryCount-1)))
	if delay > MaxDelayMillis {
		return MaxDelayMillis
	}
	return delay
}

func handleTemporaryFailure(m *amqppackage.Delivery, body []byte, failedReason string, retryCount int, recipientsToRetry []string) {
	log.Info().Msgf("Temporary failure for message: %s, error: %v, retry count: %d", body, failedReason, retryCount)

	// Compute the delay before retrying
	delay := computeDelay(retryCount)
	log.Info().Msgf("Retrying in %d milliseconds", delay)

	// publish the message into the retry queue with the delay (Dead letter Queue to the original routing key)
	var message map[string]interface{}
	err := json.Unmarshal(body, &message)
	if err != nil {
		log.Error().Err(err).Msg("Error unmarshalling message")
		return
	}

	// Only make gRPC calls if send_mail_id is present
	sendEmailID, hasSendMailID := message["send_mail_id"].(string)
	if hasSendMailID {
		mailClient, err := getMailClient()
		if err != nil {
			log.Error().Err(err).Msg("Failed to create mail client")
			return
		}

		req := mailclient.CreateRetryStatusRequest(sendEmailID, failedReason, int32(retryCount))
		_, err = mailClient.UpdateMailStatus(context.Background(), req)
		if err != nil {
			log.Error().Err(err).Msg("Failed to update mail status for retry")
			return
		}
	} else {
		log.Info().Msg("No send_mail_id found in message, skipping gRPC status update")
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
func handlePermanentFailure(message *amqppackage.Delivery, failedReason string, retryCount int, sendEmailID string, toAddress string, recipientsFailed []string) {
	log.Info().Msgf("Permanent failure for message: %s, error: %v", message.Body, failedReason)

	// Only make gRPC calls if send_mail_id is present
	if sendEmailID != "" {
		mailClient, err := getMailClient()
		if err != nil {
			log.Error().Err(err).Msg("Failed to create mail client")
			return
		}

		log.Info().Msgf("Updating mail status for failure: %s", failedReason)

		req := mailclient.CreateFailureStatusRequest(sendEmailID, failedReason, int32(retryCount))
		log.Debug().Msgf("Request: %v", req)
		_, err = mailClient.UpdateMailStatus(context.Background(), req)
		if err != nil {
			log.Error().Err(err).Msg("Failed to update mail status for failure")
			return
		}
	} else {
		log.Info().Msg("No send_mail_id provided, skipping gRPC status update")
	}

	// Only send reject mail if we have a valid toAddress
	if toAddress != "" {
		// send a email to the user with the failure reason
		rejectMail := &models.RawMail{
			Headers: map[string]any{
				"To":      []string{toAddress},
				"From":    "mailer-daemon@atomic-blend.com",
				"Subject": "MAILER DAEMON - MAIL REJECTED",
			},
			// mail is rejected because every retry failed
			TextContent: fmt.Sprintf("Your email from %s was rejected by the recipients: %s. Please contact support if you believe this is an error. The failure reason is: %s", toAddress, strings.Join(recipientsFailed, ", "), failedReason),
			HTMLContent: fmt.Sprintf("<p>Your email from %s was rejected by the recipients: %s. Please contact support if you believe this is an error. The failure reason is: %s</p>", toAddress, strings.Join(recipientsFailed, ", "), failedReason),
		}

		_, err := mailsender.SendEmail(*rejectMail, nil)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send reject mail")
		}
	} else {
		log.Info().Msg("No toAddress provided, skipping reject mail")
	}

	message.Ack(false)

}

func handleSuccess(message *amqppackage.Delivery) {
	log.Info().Msgf("Message processed successfully: %s", message.Body)

	var messageWrapper map[string]interface{}
	err := json.Unmarshal(message.Body, &messageWrapper)
	if err != nil {
		log.Error().Err(err).Msg("Error unmarshalling AMQP payload wrapper")
		return
	}

	// Only make gRPC calls if send_mail_id is present
	sendEmailID, hasSendMailID := messageWrapper["send_mail_id"].(string)
	if hasSendMailID {
		mailClient, err := getMailClient()
		if err != nil {
			log.Error().Err(err).Msg("Failed to create mail client")
			return
		}

		req := mailclient.CreateSuccessStatusRequest(sendEmailID)
		_, err = mailClient.UpdateMailStatus(context.Background(), req)
		if err != nil {
			log.Error().Err(err).Msg("Failed to update mail status for success")
			return
		}
	} else {
		log.Info().Msg("No send_mail_id found in message, skipping gRPC status update")
	}

	message.Ack(false)
}

func processSendMailMessage(message *amqppackage.Delivery, rawMail models.RawMail) error {
	// lookup the message to check if it's a retry or not
	isRetry := false
	retryCount := 0

	if message.Headers["retry-count"] != nil {
		isRetry = true
		retryCountInt, ok := message.Headers["retry-count"].(int32)
		if !ok {
			retryCount = 0
		} else {
			retryCount = int(retryCountInt)
		}
	}

	// publish the message into the retry queue with the delay (Dead letter Queue to the original routing key)
	var parsedMessage map[string]interface{}
	err := json.Unmarshal(message.Body, &parsedMessage)
	if err != nil {
		log.Error().Err(err).Msg("Error unmarshalling message")
		return nil
	}

	// Extract send_mail_id if present (optional for external emails)
	sendEmailID, hasSendMailID := parsedMessage["send_mail_id"].(string)
	if !hasSendMailID {
		log.Info().Msg("No send_mail_id found in message - treating as external email")
		sendEmailID = "" // Set empty string to indicate no tracking
	}

	// if !isRetry && retryCount > MaxRetries {
	// 	handlePermanentFailure(message.Body, errors.New("retry_limit_reached"), sendEmailID, retryCount) // No error, just a retry limit reached
	// 	return nil
	// }

	recipientsToSend := []any{}
	if !isRetry && message.Headers["recipients"] != nil {
		// split the recipients by comma
		parsedRetry := strings.SplitSeq(message.Headers["recipients"].(string), ",")
		for recipient := range parsedRetry {
			recipientsToSend = append(recipientsToSend, recipient)
		}
	}

	recipientsToRetry, err := mailsender.SendEmail(rawMail, recipientsToSend)
	if err != nil && retryCount < MaxRetries {
		retryCount++
		handleTemporaryFailure(message, message.Body, err.Error(), retryCount, recipientsToRetry)
		return nil
	} else if err != nil {
		// Extract From address safely
		fromAddress := ""
		if fromHeader, exists := rawMail.Headers["From"]; exists {
			if fromStr, ok := fromHeader.(string); ok {
				fromAddress = fromStr
			}
		}
		handlePermanentFailure(message, "retry_limit_reached", retryCount, sendEmailID, fromAddress, recipientsToRetry)
		return nil
	}

	log.Info().Msg("Email sent successfully, sending success to mail service")
	handleSuccess(message)

	return nil
}

func getMailClient() (*mailclient.MailClient, error) {
	mailClient, err := mailclient.NewMailClient()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create mail client")
		return nil, err
	}
	return mailClient, nil
}
