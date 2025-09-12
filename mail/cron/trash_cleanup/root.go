package trashcleanup

import (
	"context"

	"github.com/atomic-blend/backend/mail/repositories"
	"github.com/atomic-blend/backend/shared/utils/db"
	"github.com/rs/zerolog/log"
)

// CleanTrashCron initializes and starts the main notification cron jobs for the application.
func CleanTrashCron() {
	log.Debug().Msg("Starting trash cleanup cron jobs")

	mailRepo := repositories.NewMailRepository(db.Database)

	err := mailRepo.CleanupTrash(context.TODO())
	if err != nil {
		log.Error().Err(err).Msg("Failed to cleanup trash")
	}

}
