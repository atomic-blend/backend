// Package amqputils contains the AMQP utils
package amqputils

import (
	"os"
)

// getEnvWithFallback returns environment variable value with hierarchical fallback
// If isProducer is true, uses "PRODUCER" prefix, otherwise uses "CONSUMER" prefix
// Fallback order: {workerName}_{PRODUCER/CONSUMER}_AMQP_{suffix} -> {workerName}_AMQP_{suffix} -> AMQP_{suffix}
func getEnvWithFallback(workerName, suffix string, isProducer bool) string {
	// Skip environment variable checks in test mode
	if os.Getenv("GO_ENV") == "test" {
		return "test_value"
	}

	prefix := "CONSUMER"
	if isProducer {
		prefix = "PRODUCER"
	}

	// Try specific producer/consumer variable first
	if value := os.Getenv(workerName + "_AMQP_" + prefix + "_" + suffix); value != "" {
		return value
	}
	// Try generic worker variable
	if value := os.Getenv(workerName + "_AMQP_" + suffix); value != "" {
		return value
	}

	// Fall back to global variable
	return os.Getenv("AMQP_" + suffix)
}

// getAMQPURL returns the AMQP URL from environment variables
// If isProducer is true, it tries {workerName}_PRODUCER_AMQP_URL first
// If isProducer is false, it tries {workerName}_CONSUMER_AMQP_URL first
// Then falls back to {workerName}_AMQP_URL, then AMQP_URL
func getAMQPURL(workerName string, isProducer bool) string {
	return getEnvWithFallback(workerName, "URL", isProducer)
}

// getAMQPExchangeNames returns the AMQP exchange names from environment variables
// If isProducer is true, it tries {workerName}_PRODUCER_AMQP_EXCHANGE_NAMES first
// If isProducer is false, it tries {workerName}_CONSUMER_AMQP_EXCHANGE_NAMES first
// Then falls back to {workerName}_AMQP_EXCHANGE_NAMES, then AMQP_EXCHANGE_NAMES
func getAMQPExchangeNames(workerName string, isProducer bool) string {
	return getEnvWithFallback(workerName, "EXCHANGE_NAMES", isProducer)
}

// getAMQPQueueName returns the AMQP queue name from environment variables
// If isProducer is true, it tries {workerName}_PRODUCER_AMQP_QUEUE_NAME first
// If isProducer is false, it tries {workerName}_CONSUMER_AMQP_QUEUE_NAME first
// Then falls back to {workerName}_AMQP_QUEUE_NAME, then AMQP_QUEUE_NAME
func getAMQPQueueName(workerName string, isProducer bool) string {
	return getEnvWithFallback(workerName, "QUEUE_NAME", isProducer)
}

// getAMQPRoutingKeys returns the AMQP routing keys from environment variables
// If isProducer is true, it tries {workerName}_PRODUCER_AMQP_ROUTING_KEYS first
// If isProducer is false, it tries {workerName}_CONSUMER_AMQP_ROUTING_KEYS first
// Then falls back to {workerName}_AMQP_ROUTING_KEYS, then AMQP_ROUTING_KEYS
func getAMQPRoutingKeys(workerName string, isProducer bool) string {
	return getEnvWithFallback(workerName, "ROUTING_KEYS", isProducer)
}

// getAMQPRetryEnabled returns whether AMQP retry is enabled from environment variables
func getAMQPRetryEnabled(workerName string, isProducer bool) bool {
	return getEnvWithFallback(workerName, "RETRY_ENABLED", isProducer) == "true"
}

// getAMQPRetryExchange returns the AMQP retry exchange name from environment variables
func getAMQPRetryExchange(workerName string, isProducer bool) string {
	return getEnvWithFallback(workerName, "RETRY_EXCHANGE", isProducer)
}

// getAMQPRetryRoutingKey returns the AMQP retry routing key from environment variables
func getAMQPRetryRoutingKey(workerName string, isProducer bool) string {
	return getEnvWithFallback(workerName, "RETRY_ROUTING_KEY", isProducer)
}

func getAMQPRetryQueueName(workerName string, isProducer bool) string {
	return getEnvWithFallback(workerName, "RETRY_QUEUE_NAME", isProducer)
}

func getAMQPRetryBindingKey(workerName string, isProducer bool) string {
	return getEnvWithFallback(workerName, "RETRY_BINDING_KEY", isProducer)
}
