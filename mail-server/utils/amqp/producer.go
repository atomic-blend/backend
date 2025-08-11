package amqp

import (
	"encoding/json"
	"os"
	"strconv"
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

var conn *amqp.Connection
var ch *amqp.Channel

// InitProducerAMQP initializes the AMQP producer
func InitProducerAMQP() {
	var err error

	shortcuts.CheckRequiredEnvVar("MAIL_SERVER_PRODUCER_AMQP_URL or MAIL_SERVER_AMQP_URL or AMQP_URL", amqpURL, "amqp://user:password@localhost:5672/")
	shortcuts.CheckRequiredEnvVar("MAIL_SERVER_PRODUCER_AMQP_QUEUE_NAME or MAIL_SERVER_AMQP_QUEUE_NAME or AMQP_QUEUE_NAME", getAMQPQueueName(true), "")
	shortcuts.CheckRequiredEnvVar("MAIL_SERVER_PRODUCER_AMQP_EXCHANGE_NAMES or MAIL_SERVER_PRODUCER_AMQP_EXCHANGE_NAMES or MAIL_SERVER_AMQP_EXCHANGE_NAMES", exchangeNames, "")

	// Close existing connection if it exists
	if conn != nil || ch != nil {
		log.Info().Msg("Closing existing AMQP connection")
		CloseConnection()
	}

	//split exchange names
	exchangeNamesList := strings.Split(exchangeNames, ",")

	log.Debug().Msg("Producer connecting to AMQP")
	conn, err = amqp.Dial(amqpURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to RabbitMQ")
		return
	}

	log.Debug().Msg("Opening a channel")
	ch, err = conn.Channel()
	if err != nil {
		log.Error().Err(err).Msg("Failed to open a channel")
		conn.Close()
		return
	}

	for _, exchangeName := range exchangeNamesList {
		log.Info().Str("exchange", exchangeName).Msg("Declaring exchange")
		err = ch.ExchangeDeclare(
			exchangeName, // name
			"topic",      // type
			true,         // durable
			false,        // auto-deleted
			false,        // internal
			false,        // noWait
			nil,          // arguments
		)
		if err != nil {
			log.Warn().Err(err).Str("exchange", exchangeName).Msg("Failed to declare exchange, continuing with others")
			continue
		}
		log.Info().Str("exchange", exchangeName).Msg("Exchange declared successfully")
	}

	log.Info().Msg("Declaring queue")
	queueName := getAMQPQueueName(true)
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Info().Err(err).Msg("Failed to declare a queue")
	}

	retryEnabled := getAMQPRetryEnabled(true)
	if retryEnabled {
		log.Info().Msg("Retry queue is enabled, declaring retry DLQ")

		// Ensure the mail exchange exists for retry binding
		mailExchange := "mail"
		err = ch.ExchangeDeclare(
			mailExchange, // name
			"topic",      // type
			true,         // durable
			false,        // auto-deleted
			false,        // internal
			false,        // noWait
			nil,          // arguments
		)
		if err != nil {
			log.Warn().Err(err).Str("exchange", mailExchange).Msg("Failed to declare mail exchange for retry binding")
		}

		// Declare the dead letter exchange for retries
		retryQueueName := getAMQPRetryQueueName(true)
		retryQueueArgs := amqp.Table{
			"x-dead-letter-exchange":    getAMQPRetryExchange(true),
			"x-dead-letter-routing-key": getAMQPRetryRoutingKey(true),
		}

		// Try to declare the retry queue
		_, err = ch.QueueDeclare(retryQueueName, true, false, false, false, retryQueueArgs)
		if err != nil {
			log.Warn().Err(err).Str("queue", retryQueueName).Msg("Failed to declare retry queue, may already exist with different settings")

			log.Info().Msg("Attempting to bind existing retry queue")
			err = ch.QueueBind(
				retryQueueName,               // name of the queue
				getAMQPRetryBindingKey(true), // bindingKey
				mailExchange,                 // sourceExchange
				false,                        // noWait
				nil,                          // arguments
			)
			if err != nil {
				log.Warn().Err(err).Str("queue", retryQueueName).Msg("Failed to bind existing retry queue")
			} else {
				log.Info().Str("queue", retryQueueName).Msg("Retry queue bound to exchange")
			}
		} else {
			log.Info().Str("queue", retryQueueName).Msg("Retry queue declared successfully")
		}

		log.Debug().Msg("DLQ declared, binding to queue")
		err = ch.QueueBind(
			retryQueueName,               // name of the queue
			getAMQPRetryBindingKey(true), // bindingKey
			mailExchange,                 // sourceExchange
			false,                        // noWait
			nil,                          // arguments
		)
		if err != nil {
			log.Warn().Err(err).Str("queue", retryQueueName).Str("exchange", mailExchange).Msg("Error binding retry queue to exchange, continuing without retry binding")
		} else {
			log.Debug().Msg("Retry queue bound to exchange")
		}
	} else {
		log.Info().Msg("Retry queue is disabled, skipping retry DLQ declaration")
	}

	// Final check to ensure connection is healthy
	if IsConnectionHealthy() {
		log.Info().Msg("✅\tAMQP connection established")
	} else {
		log.Error().Msg("❌\tAMQP connection failed to establish properly")
	}
}

// IsConnectionHealthy checks if the AMQP connection and channel are healthy
func IsConnectionHealthy() bool {
	return conn != nil && ch != nil && !conn.IsClosed()
}

// CloseConnection closes the AMQP connection and channel
func CloseConnection() {
	if ch != nil {
		ch.Close()
	}
	if conn != nil {
		conn.Close()
	}
	log.Info().Msg("AMQP connection closed")
}

// PublishMessage publishes a message to the AMQP broker
func PublishMessage(exchangeName string, topic string, message map[string]interface{}, headers *amqp.Table) {
	// Skip publishing in test environment
	if os.Getenv("GO_ENV") == "test" {
		log.Debug().Msg("Skipping AMQP message publishing (test environment)")
		return
	}

	// Check if connection and channel are healthy
	if !IsConnectionHealthy() {
		log.Warn().Msg("AMQP connection or channel is not available, attempting to reconnect")
		InitProducerAMQP()
		if !IsConnectionHealthy() {
			log.Error().Msg("Failed to establish AMQP connection, cannot publish message")
			return
		}
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
	publishing := amqp.Publishing{
		ContentType: "application/json",
		Body:        encodedPayload,
		Timestamp:   time.Now(),
	}

	if headers != nil && (*headers)["delay"] != nil {
		publishing.Expiration = strconv.Itoa((*headers)["delay"].(int))
	}

	// Add headers if provided
	if headers != nil {
		publishing.Headers = *headers
	}

	err = ch.Publish(
		exchangeName, // exchange
		topic,        // routing key
		false,        // mandatory
		false,        // immediate
		publishing)
	if err != nil {
		log.Error().Err(err).Msg("Failed to publish a message")
		// Try to reconnect on publish failure
		if conn != nil && conn.IsClosed() {
			log.Info().Msg("Connection is closed, will attempt to reconnect on next publish")
		}
	}
}
