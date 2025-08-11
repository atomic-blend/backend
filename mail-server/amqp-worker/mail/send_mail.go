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
	//TODO: make the gRPC call to store the retry count and delay in the DB

	// TODO: publish the message into the retry queue with the delay (Dead letter Queue to the original routing key)
}

// handlePermanentFailure handles messages that have permanently failed
func handlePermanentFailure(body []byte, err error) {
	log.Info().Msgf("Permanent failure for message: %s, error: %v", body, err)
	//TODO: make the gRPC call to store the failure reason in the DB + status to failed
}

func sendEmail(to, subject, body string) error {
	log.Printf("Sending email to %s: %s", to, subject)
	// TODO: implement actual email sending logic with MX resolve + SMTP sending
	return nil
}

func processSendMailMessage(message *amqp.Delivery, rawMail models.RawMail) error {
	// TODO: [DONE] declare queue per worker

	// TODO: lookup the message to check if it's a retry or not
	

	// TODO: implement the first send logic
	// TODO: implement the retry logic
	// TODO: store the latest reason for failure into DB via a gRPC call
	// TODO: make gRPC calls when it's a definitive success
	// TODO: make gRPC calls when it's a definitive failure
	return nil
}
