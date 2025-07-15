package notifications

import "github.com/rs/zerolog/log"

// MainNotificationCron initializes and starts the main notification cron jobs for the application.
func MainNotificationCron() {
	log.Debug().Msg("Starting notification cron jobs")
	TaskDueNotificationCron()
	HabitReminderNotificationCron()
}
