package amqp

import (
	"os"
)

const workerName = "MAIL_SERVER"

// getEnvWithFallback returns environment variable value with hierarchical fallback
// If isProducer is true, uses "PRODUCER" prefix, otherwise uses "CONSUMER" prefix
// Fallback order: MAIL_SERVER_{PRODUCER/CONSUMER}_AMQP_{suffix} -> MAIL_SERVER_AMQP_{suffix} -> AMQP_{suffix}
func getEnvWithFallback(suffix string, isProducer bool) string {
	prefix := "CONSUMER"
	if isProducer {
		prefix = "PRODUCER"
	}

	// Try specific producer/consumer variable firs
	if value := os.Getenv(workerName + "_AMQP_" + prefix + "_" + suffix); value != "" {
		return value
	}
	// Try generic mail server variable
	if value := os.Getenv(workerName + "_AMQP_" + suffix); value != "" {
		return value
	}

	// Fall back to global variable
	return os.Getenv("AMQP_" + suffix)
}

// getAMQPURL returns the AMQP URL from environment variables
// If isProducer is true, it tries MAIL_SERVER_PRODUCER_AMQP_URL first
// If isProducer is false, it tries MAIL_SERVER_CONSUMER_AMQP_URL first
// Then falls back to MAIL_SERVER_AMQP_URL, then AMQP_URL
func getAMQPURL(isProducer bool) string {
	return getEnvWithFallback("URL", isProducer)
}

// getAMQPExchangeNames returns the AMQP exchange names from environment variables
// If isProducer is true, it tries MAIL_SERVER_PRODUCER_AMQP_EXCHANGE_NAMES first
// If isProducer is false, it tries MAIL_SERVER_CONSUMER_AMQP_EXCHANGE_NAMES first
// Then falls back to MAIL_SERVER_AMQP_EXCHANGE_NAMES, then AMQP_EXCHANGE_NAMES
func getAMQPExchangeNames(isProducer bool) string {
	return getEnvWithFallback("EXCHANGE_NAMES", isProducer)
}

// getAMQPQueueName returns the AMQP queue name from environment variables
// If isProducer is true, it tries MAIL_SERVER_PRODUCER_AMQP_QUEUE_NAME first
// If isProducer is false, it tries MAIL_SERVER_CONSUMER_AMQP_QUEUE_NAME first
// Then falls back to MAIL_SERVER_AMQP_QUEUE_NAME, then AMQP_QUEUE_NAME
func getAMQPQueueName(isProducer bool) string {
	return getEnvWithFallback("QUEUE_NAME", isProducer)
}

// getAMQPRoutingKeys returns the AMQP routing keys from environment variables
// If isProducer is true, it tries MAIL_SERVER_PRODUCER_AMQP_ROUTING_KEYS first
// If isProducer is false, it tries MAIL_SERVER_CONSUMER_AMQP_ROUTING_KEYS first
// Then falls back to MAIL_SERVER_AMQP_ROUTING_KEYS, then AMQP_ROUTING_KEYS
func getAMQPRoutingKeys(isProducer bool) string {
	return getEnvWithFallback("ROUTING_KEYS", isProducer)
}

// getAMQPRetryEnabled returns whether AMQP retry is enabled from environment variables
func getAMQPRetryEnabled(isProducer bool) bool {
	return getEnvWithFallback("RETRY_ENABLED", isProducer) == "true"
}

// getAMQPRetryExchange returns the AMQP retry exchange name from environment variables
func getAMQPRetryExchange(isProducer bool) string {
	return getEnvWithFallback("RETRY_EXCHANGE", isProducer)
}

// getAMQPRetryRoutingKey returns the AMQP retry routing key from environment variables
func getAMQPRetryRoutingKey(isProducer bool) string {
	return getEnvWithFallback("RETRY_ROUTING_KEY", isProducer)
}

func getAMQPRetryQueueName(isProducer bool) string {
	return getEnvWithFallback("RETRY_QUEUE_NAME", isProducer)
}

func getAMQPRetryBindingKey(isProducer bool) string {
	return getEnvWithFallback("RETRY_BINDING_KEY", isProducer)
}
