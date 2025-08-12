// Package amqputils contains the AMQP producer utils
package amqputils

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/atomic-blend/backend/mail/utils/shortcuts"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

var (
	// amqpURL is the URL of the AMQP broker
	amqpURL = getAMQPURL(true)
	// exchangeNames is the list of exchange names to use
	exchangeNames = getAMQPExchangeNames(true)
)

// Producer-specific connection variables
var producerConn *amqp.Connection
var producerCh *amqp.Channel

// InitProducerAMQP initializes the AMQP producer
func InitProducerAMQP() {
	var err error

	// Skip initialization in test environment
	if os.Getenv("GO_ENV") == "test" {
		log.Debug().Msg("Skipping AMQP producer initialization (test environment)")
		return
	}

	shortcuts.CheckRequiredEnvVar("MAIL_PRODUCER_AMQP_URL or MAIL_AMQP_URL or AMQP_URL", amqpURL, "amqp://user:password@localhost:5672/")
	shortcuts.CheckRequiredEnvVar("MAIL_PRODUCER_AMQP_QUEUE_NAME or MAIL_AMQP_QUEUE_NAME or AMQP_QUEUE_NAME", getAMQPQueueName(true), "")
	shortcuts.CheckRequiredEnvVar("MAIL_PRODUCER_AMQP_EXCHANGE_NAMES or MAIL_AMQP_EXCHANGE_NAMES or AMQP_EXCHANGE_NAMES", exchangeNames, "")

	//split exchange names
	exchangeNamesList := strings.Split(exchangeNames, ",")

	log.Debug().Msg("Producer connecting to AMQP")
	producerConn, err = amqp.Dial(amqpURL)
	shortcuts.FailOnError(err, "Failed to connect to RabbitMQ")

	log.Debug().Msg("Opening a channel")
	producerCh, err = producerConn.Channel()
	shortcuts.FailOnError(err, "Failed to open a channel")

	log.Info().Msg("Declaring queue")
	queueName := getAMQPQueueName(true)
	_, err = producerCh.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Info().Err(err).Msg("Failed to declare a queue")
	}

	for _, exchangeName := range exchangeNamesList {
		log.Info().Str("exchange", exchangeName).Msg("Declaring exchange")
		err = producerCh.ExchangeDeclare(
			exchangeName, // name
			"topic",      // type
			true,         // durable
			false,        // auto-deleted
			false,        // internal
			false,        // noWait
			nil,          // arguments
		)
		shortcuts.FailOnError(err, "Failed to declare the Exchange")
	}
	log.Info().Msg("✅\tAMQP connection established")
}

// PublishMessage publishes a message to the AMQP broker
func PublishMessage(exchangeName string, topic string, message map[string]interface{}) {
	// Skip publishing in test environment
	if os.Getenv("GO_ENV") == "test" || producerCh == nil {
		log.Debug().Msg("Skipping AMQP message publishing (test environment or no connection)")
		return
	}

	log.Debug().Msg("Publishing message to AMQP")
	log.Debug().Msgf("Exchange: %s, Topic: %s, Message: %v", exchangeName, topic, message)
	log.Debug().Msg("Encoding message")
	encodedPayload, err := json.Marshal(message)
	shortcuts.FailOnError(err, "Failed to encode message")

	log.Debug().Msg("Publishing message")
	err = producerCh.Publish(
		exchangeName, // exchange
		topic,        // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        encodedPayload,
			Timestamp:   time.Now(),
		})
	shortcuts.FailOnError(err, "Failed to publish a message")
}

// CloseProducerConnection closes the producer AMQP connection
func CloseProducerConnection() {
	if producerCh != nil {
		producerCh.Close()
		producerCh = nil
	}
	if producerConn != nil {
		producerConn.Close()
		producerConn = nil
	}
	log.Info().Msg("Producer AMQP connection closed")
}

// IsProducerConnectionHealthy checks if the producer AMQP connection is healthy
func IsProducerConnectionHealthy() bool {
	return producerConn != nil && producerCh != nil && !producerConn.IsClosed()
}
