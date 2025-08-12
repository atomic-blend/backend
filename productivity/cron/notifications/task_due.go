package notifications

import (
	"context"
	"os"
	"time"

	"connectrpc.com/connect"
	authv1 "github.com/atomic-blend/backend/grpc/gen/auth/v1"
	userv1 "github.com/atomic-blend/backend/grpc/gen/user/v1"
	"github.com/atomic-blend/backend/productivity/cron/notifications/payloads"
	"github.com/atomic-blend/backend/productivity/grpc/clients"
	"github.com/atomic-blend/backend/productivity/models"
	"github.com/atomic-blend/backend/productivity/repositories"
	"github.com/atomic-blend/backend/shared/utils/db"
	fcmutils "github.com/atomic-blend/backend/shared/utils/fcm_utils"
	"github.com/atomic-blend/backend/shared/utils/shortcuts"

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

	userService, err := clients.NewUserClient()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user client")
		return
	}

	log.Debug().Msg("Initializing the FCM client")
	shortcuts.CheckRequiredEnvVar("FIREBASE_PROJECT_ID", firebaseProjectID, "FIREBASE_PROJECT_ID is required for FCM")

	log.Debug().Msgf("Firebase project ID: %s", firebaseProjectID)

	fcmClient, err := fcm.NewClient(
		ctx,
		fcm.WithProjectID(firebaseProjectID),
		fcm.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create FCM client")
		return
	}

	// Get all tasks from the database
	tasks, err := taskRepo.GetAll(ctx, nil) // nil to get all tasks for all users
	if err != nil {
		log.Error().Err(err).Msg("Failed to get tasks")
		return
	}

	now := time.Now()
	log.Debug().Msgf("Current time: %s", now.Format(time.RFC3339))

	// Process each task
	for _, task := range tasks {
		log.Debug().Msgf("Processing task: %s", task.ID)

		// Check if this task should send a notification now
		notificationType := shouldSendTaskNotification(task)
		if notificationType == "" {
			log.Debug().Msgf("No notification needed for task: %s", task.ID)
			continue
		}

		log.Debug().Msgf("Notification type: %s for task: %s", notificationType, task.ID)

		// Get the user for this task
		userID := task.User.Hex()

		// get user devices using gRPC client
		req := &connect.Request[userv1.GetUserDevicesRequest]{
			Msg: &userv1.GetUserDevicesRequest{
				User: &authv1.User{
					Id: userID,
				},
			},
		}
		resp, err := userService.GetUserDevices(ctx, req)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to get user devices for task: %s", task.ID)
			continue
		}

		deviceTokens := []string{}
		for _, device := range resp.Msg.Devices {
			if device.FcmToken != "" {
				deviceTokens = append(deviceTokens, device.FcmToken)
			}
		}

		if len(deviceTokens) == 0 {
			log.Debug().Msgf("No device tokens found for user: %s", userID)
			continue
		}

		log.Debug().Msgf("Found %d device tokens for user: %s", len(deviceTokens), userID)

		// Create and send appropriate notification based on type
		var payload interface{ GetData() map[string]string }

		switch notificationType {
		case "due":
			payload = payloads.NewTaskDuePayload(task.Title)
		case "starting":
			payload = payloads.NewTaskStartingPayload(task.Title)
		case "reminder":
			// Find the specific reminder that triggered
			for _, reminder := range task.Reminders {
				if isDateNow(reminder, task.Completed) {
					payload = payloads.NewTaskReminderPayload(
						task.Title,
						reminder.Time().UTC().Format(time.RFC3339),
					)
					break
				}
			}
		}

		if payload == nil {
			log.Error().Msgf("Failed to create payload for task: %s", task.ID)
			continue
		}

		data := payload.GetData()
		fcmutils.SendMulticast(ctx, fcmClient, data, deviceTokens)
		log.Debug().Msgf("Sent %s notification for task: %s", notificationType, task.ID)
	}
}

// shouldSendTaskNotification determines if a task should send a notification at the current time
// Returns notification type ("due", "starting", "reminder") or empty string if no notification
func shouldSendTaskNotification(task *models.TaskEntity) string {
	// Check if task is due
	if isDateNow(task.EndDate, task.Completed) {
		log.Debug().Msgf("Task is due: %s", task.EndDate.Time().Format(time.RFC3339))
		return "due"
	}

	// Check if task is starting
	if isDateNow(task.StartDate, task.Completed) {
		log.Debug().Msgf("Task is starting: %s", task.StartDate.Time().Format(time.RFC3339))
		return "starting"
	}

	// Check for reminders
	for _, reminder := range task.Reminders {
		if isDateNow(reminder, task.Completed) {
			log.Debug().Msgf("Task reminder: %s", reminder.Time().Format(time.RFC3339))
			return "reminder"
		}
	}

	return ""
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
