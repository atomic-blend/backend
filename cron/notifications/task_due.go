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
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		fcm.WithProjectID(firebaseProjectID),
		fcm.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")),
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

		// get the tasks for the user
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
			log.Debug().Msgf("Current time: %s", now.Format(time.RFC3339))
			if isDateNow(task.EndDate, task.Completed) {
				log.Debug().Msgf("Task is due: %s", task.EndDate.Time().Format(time.RFC3339))
				payload := payloads.NewTaskDuePayload(
					task.Title,
				)

				data := payload.GetData()
				fcmutils.SendMulticast(ctx, fcmClient, data, deviceTokens)
				continue
			}

			if isDateNow(task.StartDate, task.Completed) {
				log.Debug().Msgf("Task is starting: %s", task.StartDate.Time().Format(time.RFC3339))
				payload := payloads.NewTaskStartingPayload(
					task.Title,
				)

				data := payload.GetData()
				fcmutils.SendMulticast(ctx, fcmClient, data, deviceTokens)
				continue
			}

			for _, reminder := range task.Reminders {
				if isDateNow(reminder, task.Completed) {
					log.Debug().Msgf("Task reminder: %s", reminder.Time().Format(time.RFC3339))
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

func isDateNow(date *primitive.DateTime, completed *bool) bool {
	now := time.Now()
	if date == nil {
		return false
	}

	dateTime := date.Time()
	return dateTime.Year() == now.Year() &&
		dateTime.Month() == now.Month() &&
		dateTime.Day() == now.Day() &&
		dateTime.Hour() == now.Hour() &&
		dateTime.Minute() == now.Minute() &&
		completed != nil && !*completed
}
