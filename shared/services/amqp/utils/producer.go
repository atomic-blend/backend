// Package amqputils contains the AMQP producer utils
package amqputils

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/atomic-blend/backend/shared/utils/shortcuts"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

// Producer-specific connection variables
var producerConn *amqp.Connection
var producerCh *amqp.Channel

// InitProducerAMQP initializes the AMQP producer
func InitProducerAMQP(workerName string) {
	var err error

	// Set values from environment variables
	amqpURL := getAMQPURL(workerName, true)
	exchangeNames := getAMQPExchangeNames(workerName, true)

	// Skip initialization in test environment
	if os.Getenv("GO_ENV") == "test" {
		log.Debug().Msg("Skipping AMQP producer initialization (test environment)")
		return
	}

	shortcuts.CheckRequiredEnvVar(workerName+"_PRODUCER_AMQP_URL or "+workerName+"_AMQP_URL or AMQP_URL", amqpURL, "amqp://user:password@localhost:5672/")
	shortcuts.CheckRequiredEnvVar(workerName+"_PRODUCER_AMQP_QUEUE_NAME or "+workerName+"_AMQP_QUEUE_NAME or AMQP_QUEUE_NAME", getAMQPQueueName(workerName, true), "")
	shortcuts.CheckRequiredEnvVar(workerName+"_PRODUCER_AMQP_EXCHANGE_NAMES or "+workerName+"_AMQP_EXCHANGE_NAMES or AMQP_EXCHANGE_NAMES", exchangeNames, "")

	//split exchange names
	exchangeNamesList := strings.Split(exchangeNames, ",")

	log.Debug().Msg("Producer connecting to AMQP")
	producerConn, err = amqp.Dial(amqpURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to RabbitMQ")
		return
	}

	log.Debug().Msg("Opening a channel")
	producerCh, err = producerConn.Channel()
	if err != nil {
		log.Error().Err(err).Msg("Failed to open a channel")
		return
	}

	log.Info().Msg("Declaring queue")
	queueName := getAMQPQueueName(workerName, true)
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
		if err != nil {
			log.Error().Err(err).Msg("Failed to declare the Exchange")
			return
		}
	}
	log.Info().Msg("✅\tAMQP connection established")
}

// PublishMessage publishes a message to the AMQP broker
func PublishMessage(exchangeName string, topic string, message map[string]interface{}, headers *amqp.Table) {
	// Skip publishing in test environment
	if os.Getenv("GO_ENV") == "test" || producerCh == nil {
		log.Debug().Msg("Skipping AMQP message publishing (test environment or no connection)")
		return
	}

	log.Debug().Msg("Publishing message to AMQP")
	log.Debug().Msgf("Exchange: %s, Topic: %s, Message: %v", exchangeName, topic, message)
	log.Debug().Msg("Encoding message")
	encodedPayload, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to encode message")
		return
	}

	log.Debug().Msg("Publishing message")
	// Handle nil headers
	var amqpHeaders amqp.Table
	if headers != nil {
		amqpHeaders = *headers
	}

	err = producerCh.Publish(
		exchangeName, // exchange
		topic,        // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        encodedPayload,
			Timestamp:   time.Now(),
			Headers:     amqpHeaders,
		})
	if err != nil {
		log.Error().Err(err).Msg("Failed to publish a message")
		return
	}
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
