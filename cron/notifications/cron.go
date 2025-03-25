package notifications

import "github.com/rs/zerolog/log"

func MainNotificationCron() {
	log.Debug().Msg("Starting notification cron jobs")
	TaskDueNotificationCron()
}
