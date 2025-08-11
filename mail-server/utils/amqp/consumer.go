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

// Consumer-specific connection variables
var consumerConn *amqp.Connection
var consumerCh *amqp.Channel

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

	consumerConn, err = amqp.Dial(getAMQPURL(false))
	shortcuts.FailOnError(err, "Failed to connect to RabbitMQ")

	exchangeNames := strings.Split(exchangeNamesRaw, ",")

	log.Debug().Msg("got Connection, getting Channel...")

	consumerCh, err = consumerConn.Channel()
	shortcuts.FailOnError(err, "Failed to open a channel")

	log.Debug().Msg("got Channel, setting QoS...")
	// set the QoS to 1 so round-robin between consumer is enabled
	consumerCh.Qos(1, 0, true)

	log.Debug().Msg("Declaring Exchanges...")
	for _, exchangeName := range exchangeNames {
		log.Info().Str("exchange", exchangeName).Msg("Declaring exchange")
		err = consumerCh.ExchangeDeclare(
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

	q, err = consumerCh.QueueDeclare(
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

	log.Info().Msgf("📊 Queue: %s | Messages: %d | Consumers: %d", q.Name, q.Messages, q.Consumers)

	// If there are messages in the queue, log them
	if q.Messages > 0 {
		log.Info().Msgf("📬 Found %d messages waiting in queue!", q.Messages)
	} else {
		log.Info().Msg("📭 Queue is empty, waiting for new messages")
	}

	routingKeys := strings.Split(routingKeysRaw, ",")
	log.Debug().Msgf("Binding Queue (%q) to Exchange with routing keys: %v", q.Name, routingKeys)

	for _, routingKey := range routingKeys {
		log.Info().Str("routingKey", routingKey).Msg("Binding to the Queue")
		splitted := strings.Split(routingKey, ":")
		log.Info().Msgf("🔗 Binding queue '%s' to exchange '%s' with routing key '%s'", q.Name, splitted[0], splitted[1])

		err = consumerCh.QueueBind(
			q.Name,      // name of the queue
			splitted[1], // bindingKey
			splitted[0], // sourceExchange
			false,       // noWait
			nil,         // arguments
		)
		shortcuts.FailOnError(err, "Error binding to the Queue")
	}

	log.Debug().Msgf("Queue bound to Exchange, starting Consume (consumer tag %q)", "worker")

	MailMessages, err = consumerCh.Consume(
		q.Name,    // queue
		queueName, // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	shortcuts.FailOnError(err, "Error consuming the Queue")

	// setup the retry queue
	if getAMQPRetryEnabled(true) {
		retryQueueName := getAMQPRetryQueueName(true)
		retryQueueArgs := amqp.Table{
			"x-dead-letter-exchange":    getAMQPRetryExchange(true),
			"x-dead-letter-routing-key": getAMQPRetryRoutingKey(true),
		}

		retryQueue, err := consumerCh.QueueDeclare(
			retryQueueName, // name
			true,           // durable
			false,          // delete when unused
			false,          // exclusive
			false,          // no-wait
			retryQueueArgs, // arguments
		)
		shortcuts.FailOnError(err, "Error declaring the Retry Queue")

		//bind the retry queue to the exchange
		err = consumerCh.QueueBind(
			retryQueue.Name,              // name of the queue
			getAMQPRetryRoutingKey(true), // bindingKey
			getAMQPRetryExchange(true),   // sourceExchange
			false,                        // noWait
			nil,                          // arguments
		)
		shortcuts.FailOnError(err, "Error binding the Retry Queue")
	}

	log.Info().Msgf("📨 Started consuming from queue: %s", q.Name)

	log.Info().Msgf("🎯 Consumer is now listening on queue '%s' for messages with routing key 'sent'", q.Name)

	if getAMQPRetryEnabled(true) {
		log.Debug().Msg("Starting retry queue consumer")
		RetryMessages, err = consumerCh.Consume(
			"retry_queue",  // queue
			"retry_worker", // consumer
			false,          // auto-ack
			false,          // exclusive
			false,          // no-local
			false,          // no-wait
			nil,            // args
		)
		shortcuts.FailOnError(err, "Error consuming the Retry Queue")
		log.Info().Msg("🔄 Started consuming from retry queue: retry_queue")
	} else {
		log.Debug().Msg("Retry queue is disabled, skipping retry queue consumer")
	}

	log.Info().Msg("✅ Consumer AMQP connection established")
	log.Info().Msgf("📋 Final status: Queue '%s' bound to exchange 'mail' with routing key 'sent'", queueName)
	log.Debug().Interface("connection_status", GetConsumerConnectionStatus()).Msg("Consumer connection details")
	log.Info().Msg("🚀 Consumer is ready to process messages!")
}

// CloseConsumerConnection closes the consumer AMQP connection
func CloseConsumerConnection() {
	if consumerCh != nil {
		consumerCh.Close()
		consumerCh = nil
	}
	if consumerConn != nil {
		consumerConn.Close()
		consumerConn = nil
	}
	log.Info().Msg("Consumer AMQP connection closed")
}

// IsConsumerConnectionHealthy checks if the consumer AMQP connection is healthy
func IsConsumerConnectionHealthy() bool {
	return consumerConn != nil && consumerCh != nil && !consumerConn.IsClosed()
}

// GetConsumerConnectionStatus returns the status of the consumer connection for debugging
func GetConsumerConnectionStatus() map[string]interface{} {
	status := map[string]interface{}{
		"consumer_conn_exists": consumerConn != nil,
		"consumer_ch_exists":   consumerCh != nil,
	}

	if consumerConn != nil {
		status["consumer_conn_closed"] = consumerConn.IsClosed()
	}

	return status
}

// AreConsumerChannelsReady checks if the consumer channels are ready to receive messages
func AreConsumerChannelsReady() bool {
	// MailMessages should always be available
	if MailMessages == nil {
		return false
	}

	// RetryMessages is only required if retry is enabled
	if getAMQPRetryEnabled(false) && RetryMessages == nil {
		return false
	}

	return true
}
