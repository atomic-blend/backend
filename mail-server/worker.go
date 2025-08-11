package main

import (
	"time"

	amqpworker "github.com/atomic-blend/backend/mail-server/amqp-worker"
	"github.com/atomic-blend/backend/mail-server/utils/amqp"
	"github.com/rs/zerolog/log"
)

func processMessages() {
	log.Info().Msg("🚀 Starting AMQP message processing worker")

	// Keep the function running indefinitely
	for {
		// Wait a moment for channels to be ready
		time.Sleep(1 * time.Second)

		// Check if channels are ready
		if !amqp.AreConsumerChannelsReady() {
			log.Error().Msg("❌ Consumer channels are not ready, retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		log.Info().Msg("✅ Consumer channels are ready")

		// Process mail messages
		go func() {
			log.Info().Msg("📨 Starting mail message consumer")
			for m := range amqp.MailMessages {
				log.Info().Str("exchange", m.Exchange).Str("routing_key", m.RoutingKey).Msg("📧 Processing mail message")
				amqpworker.RouteMessage(&m)
			}
			log.Warn().Msg("📨 Mail message consumer stopped")
		}()

		log.Info().Msg("✅ AMQP message processing worker started successfully")

		// Keep this goroutine alive forever
		select {}
	}
}
