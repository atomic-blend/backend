package main

import (
	"time"

	amqpworker "github.com/atomic-blend/backend/mail-server/amqp-worker"
	"github.com/atomic-blend/backend/mail-server/utils/amqp"
	"github.com/rs/zerolog/log"
)

func processMessages() {
	log.Info().Msg("🚀 Starting AMQP message processing worker")

	// Wait a moment for channels to be ready
	time.Sleep(1 * time.Second)

	// Check if channels are ready
	if !amqp.AreConsumerChannelsReady() {
		log.Error().Msg("❌ Consumer channels are not ready, cannot start worker")
		return
	}

	log.Info().Msg("✅ Consumer channels are ready")

	// Process mail messages
	go func() {
		log.Info().Msg("📨 Starting mail message consumer")
		messageCount := 0
		for m := range amqp.MailMessages {
			messageCount++
			log.Info().Str("exchange", m.Exchange).Str("routing_key", m.RoutingKey).Msg("📧 Processing mail message")
			amqpworker.RouteMessage(&m)
		}
		log.Warn().Msg("📨 Mail message consumer stopped")
	}()

	// Process retry messages
	go func() {
		log.Info().Msg("🔄 Starting retry message consumer")
		messageCount := 0
		for m := range amqp.RetryMessages {
			messageCount++
			log.Info().Str("exchange", m.Exchange).Str("routing_key", m.RoutingKey).Msg("🔄 Processing retry message")
			amqpworker.RouteMessage(&m)
		}
		log.Warn().Msg("🔄 Retry message consumer stopped")
	}()

	log.Info().Msg("✅ AMQP message processing worker started successfully")
}
