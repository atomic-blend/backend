package cron

import (
	"atomic_blend_api/cron/notifications"
	"github.com/rs/zerolog/log"
)

// MainCron initializes and starts the main cron jobs for the application.
func MainCron() {
	log.Debug().Msg("Starting cron jobs")
	notifications.MainNotificationCron()
	// Add more cron jobs here as needed
}
