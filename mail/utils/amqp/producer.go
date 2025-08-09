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
	amqpURL = os.Getenv("AMQP_URL")
	// exchangeNames is the list of exchange names to use
	exchangeNames = os.Getenv("AMQP_EXCHANGE_NAMES")
)

var conn *amqp.Connection
var ch *amqp.Channel

// InitProducerAmqp initializes the AMQP producer
func InitProducerAmqp() {
	var err error

	shortcuts.CheckRequiredEnvVar("AMQP_URL", amqpURL, "amqp://user:password@localhost:5672/")
	shortcuts.CheckRequiredEnvVar("AMQP_QUEUE_NAME", queueName, "")
	shortcuts.CheckRequiredEnvVar("AMQP_EXCHANGE_NAMES", exchangeNames, "")

	//split exchange names
	exchangeNames := strings.Split(exchangeNames, ",")

	log.Debug().Msg("Connecting to AMQP")
	conn, err = amqp.Dial(amqpURL)
	shortcuts.FailOnError(err, "Failed to connect to RabbitMQ")

	log.Debug().Msg("Opening a channel")
	ch, err = conn.Channel()
	shortcuts.FailOnError(err, "Failed to open a channel")

	log.Info().Msg("Declaring queue")
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Info().Err(err).Msg("Failed to declare a queue")
	}

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
	log.Info().Msg("✅\tAMQP connection established")
}

// PublishMessage publishes a message to the AMQP broker
func PublishMessage(exchangeName string, topic string, message map[string]interface{}) {
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
