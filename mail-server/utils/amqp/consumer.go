package amqp

import (
	"strings"

	"github.com/atomic-blend/backend/mail/utils/shortcuts"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

var (
	exchangeNamesRaw = getAMQPExchangeNames(false)
	queueName        = getAMQPQueueName(false)
	routingKeysRaw   = getAMQPRoutingKeys(false)
)

// MailMessages is the channel for the mail_queue AMQP messages
var MailMessages <-chan amqp.Delivery

// RetryMessages is the channel for the retry_queue AMQP messages
var RetryMessages <-chan amqp.Delivery

// InitConsumerAmqp initializes the AMQP consumer
func InitConsumerAmqp() {
	log.Debug().Msg("Initializing AMQP Consumer")
	var err error
	var q amqp.Queue

	shortcuts.CheckRequiredEnvVar("MAIL_SERVER_CONSUMER_AMQP_URL or MAIL_SERVER_AMQP_URL or AMQP_URL", getAMQPURL(false), "amqp://user:password@localhost:5672")
	shortcuts.CheckRequiredEnvVar("MAIL_SERVER_CONSUMER_AMQP_EXCHANGE_NAMES or MAIL_SERVER_AMQP_EXCHANGE_NAMES or AMQP_EXCHANGE_NAMES", exchangeNamesRaw, "")
	shortcuts.CheckRequiredEnvVar("MAIL_SERVER_CONSUMER_AMQP_QUEUE_NAME or MAIL_SERVER_AMQP_QUEUE_NAME or AMQP_QUEUE_NAME", queueName, "")
	shortcuts.CheckRequiredEnvVar("MAIL_SERVER_CONSUMER_AMQP_ROUTING_KEYS or MAIL_SERVER_AMQP_ROUTING_KEYS or AMQP_ROUTING_KEYS", routingKeysRaw, "")

	conn, err = amqp.Dial(getAMQPURL(false))
	shortcuts.FailOnError(err, "Failed to connect to RabbitMQ")

	exchangeNames := strings.Split(exchangeNamesRaw, ",")

	log.Debug().Msg("got Connection, getting Channel...")

	ch, err = conn.Channel()
	shortcuts.FailOnError(err, "Failed to open a channel")

	log.Debug().Msg("got Channel, setting QoS...")
	// set the QoS to 1 so round-robin between consumer is enabled
	ch.Qos(1, 0, true)

	log.Debug().Msg("Declaring Exchanges...")
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

	log.Debug().Msgf("declared Exchange, declaring Queue (%s)", queueName)

	q, err = ch.QueueDeclare(
		queueName, // name, leave empty to generate a unique name
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	shortcuts.FailOnError(err, "Error declaring the Queue")

	log.Debug().Msgf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		q.Name, q.Messages, q.Consumers, routingKeysRaw)

	routingKeys := strings.Split(routingKeysRaw, ",")
	log.Debug().Msgf("Binding Queue (%q) to Exchange with routing keys: %v", q.Name, routingKeys)

	for _, routingKey := range routingKeys {
		log.Info().Str("routingKey", routingKey).Msg("Binding to the Queue")
		splitted := strings.Split(routingKey, ":")
		err = ch.QueueBind(
			q.Name,      // name of the queue
			splitted[1], // bindingKey
			splitted[0], // sourceExchange
			false,       // noWait
			nil,         // arguments
		)
		shortcuts.FailOnError(err, "Error binding to the Queue")
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

	log.Debug().Msgf("Queue bound to Exchange, starting Consume (consumer tag %q)", "worker")

	MailMessages, err = ch.Consume(
		q.Name,    // queue
		queueName, // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	shortcuts.FailOnError(err, "Error consuming the Queue")

	if retryEnabled {
		log.Debug().Msg("Starting retry queue consumer")
		RetryMessages, err = ch.Consume(
			"retry_queue",  // queue
			"retry_worker", // consumer
			false,          // auto-ack
			false,          // exclusive
			false,          // no-local
			false,          // no-wait
			nil,            // args
		)
		shortcuts.FailOnError(err, "Error consuming the Retry Queue")
	} else {
		log.Debug().Msg("Retry queue is disabled, skipping retry queue consumer")
	}
}
