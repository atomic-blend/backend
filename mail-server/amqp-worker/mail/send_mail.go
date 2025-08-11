package mail

import (
	"math"

	"github.com/atomic-blend/backend/mail/models"
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

func handleTemporaryFailure(body []byte, err error, retryCount int) {
	log.Info().Msgf("Temporary failure for message: %s, error: %v, retry count: %d", body, err, retryCount)

	// Compute the delay before retrying
	delay := computeDelay(retryCount)
	log.Info().Msgf("Retrying in %d milliseconds", delay)
	//TODO: make the gRPC call to store the retry count, delay in the DB and the reason for failure

	// TODO: publish the message into the retry queue with the delay (Dead letter Queue to the original routing key)
}

// handlePermanentFailure handles messages that have permanently failed
func handlePermanentFailure(body []byte, err error) {
	log.Info().Msgf("Permanent failure for message: %s, error: %v", body, err)
	//TODO: make the gRPC call to store the failure reason in the DB + status to failed
}

func sendEmail(mail models.RawMail) error {
	log.Info().Interface("To", mail.Headers["To"]).Interface("From", mail.Headers["From"]).Msg("Sending email")
	// TODO: implement actual email sending logic with MX resolve + SMTP sending
	return nil
}

func processSendMailMessage(message *amqp.Delivery, rawMail models.RawMail) error {
	// TODO: [DONE] declare queue per worker

	// TODO: [DONE] lookup the message to check if it's a retry or not
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

	err := sendEmail(rawMail)
	if err != nil && retryCount < MaxRetries {
		handleTemporaryFailure(message.Body, err, retryCount)
		return nil
	} else if err != nil {
		handlePermanentFailure(message.Body, err)
		return nil
	}

	log.Info().Msg("Email sent successfully, sending success to mail service")
	//TODO: make the gRPC call to store the success in the DB


	log.Info().Msg("Acknowledging message")

	// If email sent successfully, acknowledge the message
	if err := message.Ack(false); err != nil {
		log.Error().Err(err).Msg("Failed to acknowledge message")
	}

	log.Info().Msg("Email sent successfully")

	return nil
}
