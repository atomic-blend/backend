package main

import (
	"fmt"
	"time"

	amqpworker "github.com/atomic-blend/backend/mail-server/amqp-worker"
	"github.com/atomic-blend/backend/mail-server/utils/amqp"
	"github.com/rs/zerolog/log"
)

func processMessages() {
	log.Info().Msg("🚀 Starting AMQP message processing worker")

	// Keep the function running indefinitely
	for {
		// Add panic recovery
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Error().Interface("panic", r).Msg("❌ Worker panicked, restarting in 10 seconds...")
					time.Sleep(10 * time.Second)
				}
			}()

			// Wait a moment for channels to be ready
			time.Sleep(1 * time.Second)

			// Check if channels are ready
			if !amqp.AreConsumerChannelsReady() {
				log.Error().Msg("❌ Consumer channels are not ready, attempting to reinitialize...")

				// Try to reinitialize the consumer connection
				if err := reinitializeConsumer(); err != nil {
					log.Error().Err(err).Msg("❌ Failed to reinitialize consumer, retrying in 10 seconds...")
					time.Sleep(10 * time.Second)
					return
				}

				log.Info().Msg("✅ Consumer reinitialized successfully")
			}

			log.Info().Msg("✅ Consumer channels are ready")

			// Process mail messages in a loop that can restart if needed
			processMailMessages()

			log.Warn().Msg("📨 Mail message processing stopped, restarting in 5 seconds...")
			time.Sleep(5 * time.Second)
		}()
	}
}

func reinitializeConsumer() error {
	log.Info().Msg("🔄 Reinitializing AMQP consumer connection...")

	// Close existing connection if any
	amqp.CloseConsumerConnection()

	// Wait a moment for cleanup
	time.Sleep(2 * time.Second)

	// Reinitialize the consumer
	amqp.InitConsumerAMQP()

	// Wait a moment for the connection to stabilize
	time.Sleep(2 * time.Second)

	// Verify the connection is healthy
	if !amqp.AreConsumerChannelsReady() {
		return fmt.Errorf("consumer channels are still not ready after reinitialization")
	}

	return nil
}

func processMailMessages() {
	log.Info().Msg("📨 Starting mail message consumer")

	// Get the current mail messages channel
	mailMessages := amqp.MailMessages
	if mailMessages == nil {
		log.Error().Msg("❌ Mail messages channel is nil, cannot process messages")
		return
	}

	// Process messages in a loop
	for {
		// Check channel health before processing
		if !amqp.AreConsumerChannelsReady() {
			log.Warn().Msg("📨 Consumer channels are no longer healthy, restarting consumer")
			return
		}

		select {
		case message, ok := <-mailMessages:
			if !ok {
				// Channel is closed, exit this function to restart
				log.Warn().Msg("📨 Mail messages channel closed, restarting consumer")
				return
			}

			log.Info().Str("exchange", message.Exchange).Str("routing_key", message.RoutingKey).Msg("📧 Processing mail message")

			// Process the message - let the existing logic handle acknowledgment
			amqpworker.RouteMessage(&message)

		case <-time.After(30 * time.Second):
			// Check if the channel is still healthy every 30 seconds
			if !amqp.AreConsumerChannelsReady() {
				log.Warn().Msg("📨 Consumer channels are no longer healthy, restarting consumer")
				return
			}
		}
	}
}
