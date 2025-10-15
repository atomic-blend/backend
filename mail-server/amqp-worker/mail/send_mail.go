package mail

import (
	"context"
	"encoding/json"
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

func handleTemporaryFailure(m *amqppackage.Delivery, body []byte, failedReason error, retryCount int, recipientsToRetry []string) {
	log.Info().Msgf("Temporary failure for message: %s, error: %v, retry count: %d", body, failedReason, retryCount)

	// Compute the delay before retrying
	delay := computeDelay(retryCount)
	log.Info().Msgf("Retrying in %d milliseconds", delay)

	mailClient, err := getMailClient()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create mail client")
		return
	}

	// publish the message into the retry queue with the delay (Dead letter Queue to the original routing key)
	var message map[string]interface{}
	err = json.Unmarshal(body, &message)
	if err != nil {
		log.Error().Err(err).Msg("Error unmarshalling message")
		return
	}

	sendEmailID, ok := message["send_mail_id"].(string)
	if !ok {
		log.Error().Msg("send_email_id not found in message")
		return
	}

	req := mailclient.CreateRetryStatusRequest(sendEmailID, failedReason.Error(), int32(retryCount))
	_, err = mailClient.UpdateMailStatus(context.Background(), req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update mail status for retry")
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
	//TODO: send a email to the user with the failure reason
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

	req := mailclient.CreateSuccessStatusRequest(sendEmailID)
	_, err = mailClient.UpdateMailStatus(context.Background(), req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update mail status for success")
		return
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

	if !isRetry && retryCount > MaxRetries {
		handlePermanentFailure(message.Body, nil) // No error, just a retry limit reached
		return nil
	}

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

func getMailClient() (*mailclient.MailClient, error) {
	mailClient, err := mailclient.NewMailClient()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create mail client")
		return nil, err
	}
	return mailClient, nil
}
