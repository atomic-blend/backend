package amqp

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

var conn *amqp.Connection
var ch *amqp.Channel

// InitProducerAmqp initializes the AMQP producer
func InitProducerAmqp() {
	var err error

	shortcuts.CheckRequiredEnvVar("MAIL_SERVER_PRODUCER_AMQP_URL or MAIL_SERVER_AMQP_URL or AMQP_URL", amqpURL, "amqp://user:password@localhost:5672/")
	shortcuts.CheckRequiredEnvVar("MAIL_SERVER_PRODUCER_AMQP_QUEUE_NAME or MAIL_SERVER_AMQP_QUEUE_NAME or AMQP_QUEUE_NAME", getAMQPQueueName(true), "")
	shortcuts.CheckRequiredEnvVar("MAIL_SERVER_PRODUCER_AMQP_EXCHANGE_NAMES or MAIL_SERVER_AMQP_EXCHANGE_NAMES or AMQP_EXCHANGE_NAMES", exchangeNames, "")

	//split exchange names
	exchangeNames := strings.Split(exchangeNames, ",")

	log.Debug().Msg("Producer connecting to AMQP")
	conn, err = amqp.Dial(amqpURL)
	shortcuts.FailOnError(err, "Failed to connect to RabbitMQ")

	log.Debug().Msg("Opening a channel")
	ch, err = conn.Channel()
	shortcuts.FailOnError(err, "Failed to open a channel")

	for _, exchangeName := range exchangeNames {
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
		shortcuts.FailOnError(err, "Failed to declare the Exchange")
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
		// Declare the dead letter exchange for retries
		ch.QueueDeclare("retry_queue", true, false, false, false, amqp.Table{
			"x-dead-letter-exchange":    "mail",
			"x-dead-letter-routing-key": "sent",
		})

		log.Debug().Msg("DLQ declared, binding to queue")
		err = ch.QueueBind(
			"retry_queue", // name of the queue
			"send_retry",  // bindingKey
			"mail",        // sourceExchange
			false,         // noWait
			nil,           // arguments
		)
		shortcuts.FailOnError(err, "Error binding retry queue to exchange")
		log.Debug().Msg("Retry queue bound to exchange")
	} else {
		log.Info().Msg("Retry queue is disabled, skipping retry DLQ declaration")
	}

	log.Info().Msg("✅\tAMQP connection established")
}

// PublishMessage publishes a message to the AMQP broker
func PublishMessage(exchangeName string, topic string, message map[string]interface{}, headers *amqp.Table) {
	// Skip publishing in test environment
	if os.Getenv("GO_ENV") == "test" || ch == nil {
		log.Debug().Msg("Skipping AMQP message publishing (test environment or no connection)")
		return
	}

	log.Debug().Msg("Publishing message to AMQP")
	log.Debug().Msgf("Exchange: %s, Topic: %s, Message: %v", exchangeName, topic, message)
	log.Debug().Msg("Encoding message")
	encodedPayload, err := json.Marshal(message)
	shortcuts.FailOnError(err, "Failed to encode message")

	log.Debug().Msg("Publishing message")
	err = ch.Publish(
		exchangeName, // exchange
		topic,        // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Transient,
			ContentType:  "application/json",
			Body:         encodedPayload,
			Timestamp:    time.Now(),
		})
	shortcuts.FailOnError(err, "Failed to publish a message")
}
