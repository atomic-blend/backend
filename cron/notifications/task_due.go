package notifications

import (
	"atomic_blend_api/models"
	"atomic_blend_api/repositories"
	"atomic_blend_api/utils/db"
	"context"
	"github.com/rs/zerolog/log"
)

func TaskDueNotificationCron() {
	log.Debug().Msg("Starting task due notification cron job")
	ctx := context.TODO()
	taskRepo := repositories.NewTaskRepository(db.Database)
	userRepo := repositories.NewUserRepository(db.Database)

	// get a cursor to all the users
	userCursor, err := userRepo.GetAllIterable(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user cursor")
		return
	}
	defer userCursor.Close(ctx)

	// iterate through the cursor with user access
	for userCursor.Next(ctx) {
		var user models.UserEntity
		if err := userCursor.Decode(&user); err != nil {
			log.Error().Err(err).Msg("Failed to decode user")
			continue
		}
		log.Debug().Msgf("Processing user: %s", user.ID)
		// get the tasks for the user
		tasks, err := taskRepo.GetAll(ctx, user.ID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get tasks for user")
			continue
		}

		for _, task := range tasks {
			//TODO: implement task due notification logic
			log.Debug().Msgf("Processing task: %s", task.ID)
		}
	}
}
