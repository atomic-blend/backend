package notifications

import (
	"atomic_blend_api/cron/notifications/payloads"
	"atomic_blend_api/models"
	"atomic_blend_api/repositories"
	"atomic_blend_api/utils/db"
	fcmutils "atomic_blend_api/utils/fcm_utils"
	"atomic_blend_api/utils/shortcuts"
	"context"
	"os"
	"time"

	fcm "github.com/appleboy/go-fcm"
	"github.com/rs/zerolog/log"
)

var (
	// firebaseProjectID is the Firebase project ID used for FCM
	firebaseProjectID = os.Getenv("FIREBASE_PROJECT_ID")
)

// TaskDueNotificationCron initializes and starts the task due notification cron job.
func TaskDueNotificationCron() {
	log.Debug().Msg("Starting task due notification cron job")
	ctx := context.TODO()
	taskRepo := repositories.NewTaskRepository(db.Database)
	userRepo := repositories.NewUserRepository(db.Database)

	log.Debug().Msg("Initializing the FCM client")
	shortcuts.CheckRequiredEnvVar("FIREBASE_PROJECT_ID", firebaseProjectID, "FIREBASE_PROJECT_ID is required for FCM")
	log.Debug().Msgf("Firebase project ID: %s", firebaseProjectID)
	fcmClient, err := fcm.NewClient(
		ctx,
		fcm.WithProjectID(
			firebaseProjectID,
		),
		// initial with service account
		// fcm.WithServiceAccount("my-client-id@my-project-id.iam.gserviceaccount.com"),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create FCM client")
		return
	}

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
		//TODO: replace with a cursor here

		tasks, err := taskRepo.GetAll(ctx, user.ID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get tasks for user")
			continue
		}

		deviceTokens := []string{}
		for _, device := range user.Devices {
			deviceTokens = append(deviceTokens, device.FcmToken)
		}

		for _, task := range tasks {
			//TODO: implement task due notification logic
			log.Debug().Msgf("Processing task: %s", task.ID)

			// send notification to the user when:
			// - task is due and the due date is equal to the current date
			// - task have a reminder set and the reminder is hour and minute equal to the current time
			now := time.Now()
			if task.StartDate == nil && task.EndDate.Time().Hour() == now.Hour() && task.EndDate.Time().Minute() == now.Minute() {
				payload := payloads.NewTaskDuePayload(
					task.Title,
				)

				data := payload.GetData()
				fcmutils.SendMulticast(ctx, fcmClient, data, deviceTokens)
				continue
			}

			if task.StartDate != nil && task.StartDate.Time().Hour() == now.Hour() && task.StartDate.Time().Minute() == now.Minute() {
				payload := payloads.NewTaskStartingPayload(
					task.Title,
				)

				data := payload.GetData()
				fcmutils.SendMulticast(ctx, fcmClient, data, deviceTokens)
				continue
			}

			for _, reminder := range task.Reminders {
				if reminder.Time().Hour() == now.Hour() && reminder.Time().Minute() == now.Minute() {
					payload := payloads.NewTaskReminderPayload(
						task.Title,
						reminder.Time().UTC().Format(time.RFC3339),
					)

					data := payload.GetData()
					fcmutils.SendMulticast(ctx, fcmClient, data, deviceTokens)
					continue
				}
			}
		}
	}
}
