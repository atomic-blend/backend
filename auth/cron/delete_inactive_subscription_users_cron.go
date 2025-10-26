package cron

import "github.com/rs/zerolog/log"

func DeleteInactiveSubscriptionUsersCron() {
	log.Info().Msg("Running DeleteInactiveSubscriptionUsersCron")
	// TODO: get the user with subscriptionId == nil || status == cancelled
	// for more than X days (from env, since cancellation date if cancelled, or createdAt date if no subscriptionId)
}
