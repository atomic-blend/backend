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

func HabitReminderNotificationCron() {
	log.Debug().Msg("Starting habit due notification cron job")
	ctx := context.TODO()
	userRepo := repositories.NewUserRepository(db.Database)
	habitRepo := repositories.NewHabitRepository(db.Database)

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

		habits, err := habitRepo.GetAll(ctx, user.ID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get habits")
			continue
		}

		deviceTokens := []string{}
		for _, device := range user.Devices {
			deviceTokens = append(deviceTokens, device.FcmToken)
		}

		for _, habit := range habits {
			log.Debug().Msgf("Processing habit: %s", habit.ID)
			now := time.Now()
			log.Debug().Msgf("Current time: %s", now.Format(time.RFC3339))

			switch *habit.Frequency {
			case models.FrequencyDaily:
				//TODO
				if habit.EndDate != nil && !now.Before(habit.EndDate.Time()) {
					log.Debug().Msgf("Habit end date is in the past: %s", habit.EndDate.Time().Format(time.RFC3339))
					continue
				}
				// if current day of week is not in days of week, skip
				if !shortcuts.ContainsInt(*habit.DaysOfWeek, int(now.Weekday()-1)) {
					log.Debug().Msgf("Current day of week is not in days of week: %d", int(now.Weekday()))
					continue
				}

				// if current time is not in reminders, skip
				var reminderToSend time.Time
				for _, reminder := range habit.Reminders {
					log.Debug().Msgf("Reminder: %s", reminder)
					reminderTime, err := time.Parse("15:04", reminder)
					if err != nil {
						log.Error().Err(err).Msgf("Failed to parse reminder time: %s", reminder)
						continue
					}
					if now.Hour() == reminderTime.Hour() && now.Minute() == reminderTime.Minute() {
						log.Debug().Msgf("Reminder time: %s", reminderTime.Format(time.RFC3339))
						reminderToSend = reminderTime
					}
				}
				if reminderToSend.IsZero() {
					log.Debug().Msgf("No reminder to send for habit: %s", habit.ID)
					continue
				}

				// send notification
				payload := payloads.NewHabitReminderPayload(
					*habit.Name,
					*habit.Citation,
					habit.Emoji,
				)
				log.Debug().Msgf("Payload: %v", payload)

				data := payload.GetData()

				log.Debug().Msgf("Data: %v", data)

				// send notification to the user
				fcmutils.SendMulticast(ctx, fcmClient, data, deviceTokens)
				log.Debug().Msgf("Sent notification for habit: %s", habit.ID)
			case models.FrequencyWeekly:
				//TODO
				if habit.EndDate != nil && !now.Before(habit.EndDate.Time()) {
					log.Debug().Msgf("Habit end date is in the past: %s", habit.EndDate.Time().Format(time.RFC3339))
					continue
				}
				// if current day of week is not in days of week, skip
				if !shortcuts.ContainsInt(*habit.DaysOfWeek, int(now.Weekday()-1)) {
					log.Debug().Msgf("Current day of week is not in days of week: %d", int(now.Weekday()))
					continue
				}

				// if current time is not in reminders, skip
				var reminderToSend time.Time
				for _, reminder := range habit.Reminders {
					reminderTime, err := time.Parse("15:04", reminder)
					if err != nil {
						log.Error().Err(err).Msgf("Failed to parse reminder time: %s", reminder)
						continue
					}
					if now.Hour() == reminderTime.Hour() && now.Minute() == reminderTime.Minute() {
						log.Debug().Msgf("Reminder time: %s", reminderTime.Format(time.RFC3339))
						reminderToSend = reminderTime
					}
				}
				if reminderToSend.IsZero() {
					log.Debug().Msgf("No reminder to send for habit: %s", habit.ID)
					continue
				}

				// send notification
				payload := payloads.NewHabitReminderPayload(
					*habit.Name,
					*habit.Citation,
					habit.Emoji,
				)

				data := payload.GetData()

				// send notification to the user
				fcmutils.SendMulticast(ctx, fcmClient, data, deviceTokens)
				log.Debug().Msgf("Sent notification for habit: %s", habit.ID)
			case models.FrequencyMonthly:
				//TODO
				if habit.EndDate != nil && !now.Before(habit.EndDate.Time()) {
					continue
				}
				// if current day of month is not in days of week, skip
				midnightNow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
				if !shortcuts.ContainsDateTime(*habit.DaysOfMonth, midnightNow) {
					continue
				}

				// if current time is not in reminders, skip
				var reminderToSend time.Time
				for _, reminder := range habit.Reminders {
					reminderTime, err := time.Parse("15:04", reminder)
					if err != nil {
						log.Error().Err(err).Msgf("Failed to parse reminder time: %s", reminder)
						continue
					}
					if now.Hour() == reminderTime.Hour() && now.Minute() == reminderTime.Minute() {
						reminderToSend = reminderTime
					}
				}
				if reminderToSend.IsZero() {
					log.Debug().Msgf("No reminder to send for habit: %s", habit.ID)
					continue
				}

				// send notification
				payload := payloads.NewHabitReminderPayload(
					*habit.Name,
					*habit.Citation,
					habit.Emoji,
				)

				data := payload.GetData()

				// send notification to the user
				fcmutils.SendMulticast(ctx, fcmClient, data, deviceTokens)
				log.Debug().Msgf("Sent notification for habit: %s", habit.ID)
			case models.FrequencyRepeating:
				//TODO
				if habit.EndDate != nil && !now.Before(habit.EndDate.Time()) {
					continue
				}

				lastEntryDate := habit.StartDate

				//search for latest entry date in entries
				for _, entry := range habit.Entries {
					if entry.EntryDate.Time().After(lastEntryDate.Time()) {
						lastEntryDate = &entry.EntryDate
					}
				}

				nextOccurrence := lastEntryDate.Time().Add(time.Duration(*habit.Duration) * time.Millisecond)
				if nextOccurrence.After(now) {
					log.Debug().Msgf("Next occurrence is in the future: %s", nextOccurrence.Format(time.RFC3339))
					continue
				}

				// if current time is not in reminders, skip
				var reminderToSend time.Time
				for _, reminder := range habit.Reminders {
					reminderTime, err := time.Parse("15:04", reminder)
					if err != nil {
						log.Error().Err(err).Msgf("Failed to parse reminder time: %s", reminder)
						continue
					}
					if now.Hour() == reminderTime.Hour() && now.Minute() == reminderTime.Minute() {
						reminderToSend = reminderTime
					}
				}

				if reminderToSend.IsZero() {

					log.Debug().Msgf("No reminder to send for habit: %s", habit.ID)
					continue
				}

				// send notification
				payload := payloads.NewHabitReminderPayload(
					*habit.Name,
					*habit.Citation,
					habit.Emoji,
				)

				data := payload.GetData()

				// send notification to the user
				fcmutils.SendMulticast(ctx, fcmClient, data, deviceTokens)
				log.Debug().Msgf("Sent notification for habit: %s", habit.ID)
			default:
				log.Error().Msgf("Unknown frequency: %s", *habit.Frequency)
				continue
			}
		}
	}

}
