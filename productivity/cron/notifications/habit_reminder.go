package notifications

import (
	"context"
	"os"
	"time"

	"github.com/atomic-blend/backend/productivity/cron/notifications/payloads"
	"github.com/atomic-blend/backend/productivity/models"
	"github.com/atomic-blend/backend/productivity/repositories"
	"github.com/atomic-blend/backend/productivity/utils/db"
	fcmutils "github.com/atomic-blend/backend/productivity/utils/fcm_utils"
	"github.com/atomic-blend/backend/productivity/utils/shortcuts"

	fcm "github.com/appleboy/go-fcm"
	"github.com/rs/zerolog/log"
)

// HabitReminderNotificationCron is a cron job that sends notifications to users for their habits
func HabitReminderNotificationCron() {
	log.Debug().Msg("Starting habit due notification cron job")
	ctx := context.TODO()

	userRepo := repositories.NewUserRepository(db.Database)
	habitRepo := repositories.NewHabitRepository(db.Database)

	log.Debug().Msg("Initializing the FCM client")
	firebaseProjectID := os.Getenv("FIREBASE_PROJECT_ID")
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

	// Get all habits from the database
	habits, err := habitRepo.GetAll(ctx, nil) // nil to get all habits for all users
	if err != nil {
		log.Error().Err(err).Msg("Failed to get habits")
		return
	}

	now := time.Now()
	log.Debug().Msgf("Current time: %s", now.Format(time.RFC3339))

	// Process each habit
	for _, habit := range habits {
		log.Debug().Msgf("Processing habit: %s", habit.ID)

		// Check if this habit should send a notification now
		shouldSendNotification, reminderToSend := shouldSendHabitNotification(habit, now)
		if !shouldSendNotification {
			log.Debug().Msgf("No notification needed for habit: %s", habit.ID)
			continue
		}

		log.Debug().Msgf("Reminder time matched: %s", reminderToSend.Format(time.RFC3339))

		// Get the user for this habit
		user, err := userRepo.GetByID(ctx, habit.UserID.Hex())
		if err != nil {
			log.Error().Err(err).Msgf("Failed to get user for habit: %s", habit.ID)
			continue
		}
		if user == nil {
			log.Error().Msgf("User not found for habit: %s", habit.ID)
			continue
		}

		// Extract device tokens
		deviceTokens := []string{}
		for _, device := range user.Devices {
			deviceTokens = append(deviceTokens, device.FcmToken)
		}

		if len(deviceTokens) == 0 {
			log.Debug().Msgf("No device tokens found for user: %s", user.ID.Hex())
			continue
		}

		// Create and send notification
		payload := payloads.NewHabitReminderPayload(
			*habit.Name,
			*habit.Citation,
			habit.Emoji,
		)
		log.Debug().Msgf("Payload: %v", payload)

		data := payload.GetData()
		log.Debug().Msgf("Data: %v", data)

		// Send notification to the user
		fcmutils.SendMulticast(ctx, fcmClient, data, deviceTokens)
		log.Debug().Msgf("Sent notification for habit: %s", habit.ID)
	}
}

// shouldSendHabitNotification determines if a habit should send a notification at the current time
// Returns true and the reminder time if a notification should be sent, false otherwise
func shouldSendHabitNotification(habit *models.Habit, now time.Time) (bool, time.Time) {
	// Check if habit has ended
	if habit.EndDate != nil && !now.Before(habit.EndDate.Time()) {
		log.Debug().Msgf("Habit end date is in the past: %s", habit.EndDate.Time().Format(time.RFC3339))
		return false, time.Time{}
	}

	switch *habit.Frequency {
	case models.FrequencyDaily, models.FrequencyWeekly:
		return shouldSendDailyWeeklyNotification(habit, now)
	case models.FrequencyMonthly:
		return shouldSendMonthlyNotification(habit, now)
	case models.FrequencyRepeating:
		return shouldSendRepeatingNotification(habit, now)
	default:
		log.Error().Msgf("Unknown frequency: %s", *habit.Frequency)
		return false, time.Time{}
	}
}

// shouldSendDailyWeeklyNotification handles daily and weekly frequency logic
func shouldSendDailyWeeklyNotification(habit *models.Habit, now time.Time) (bool, time.Time) {
	// Check if current day of week is in days of week
	if !shortcuts.ContainsInt(*habit.DaysOfWeek, int(now.Weekday()-1)) {
		log.Debug().Msgf("Current day of week is not in days of week: %d", int(now.Weekday()))
		return false, time.Time{}
	}

	// Check if current time matches any reminder
	return findMatchingReminder(habit.Reminders, now)
}

// shouldSendMonthlyNotification handles monthly frequency logic
func shouldSendMonthlyNotification(habit *models.Habit, now time.Time) (bool, time.Time) {
	// Check if current day of month is in days of month
	midnightNow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if !shortcuts.ContainsDateTime(*habit.DaysOfMonth, midnightNow) {
		return false, time.Time{}
	}

	// Check if current time matches any reminder
	return findMatchingReminder(habit.Reminders, now)
}

// shouldSendRepeatingNotification handles repeating frequency logic
func shouldSendRepeatingNotification(habit *models.Habit, now time.Time) (bool, time.Time) {
	lastEntryDate := habit.StartDate

	// Search for latest entry date in entries
	for _, entry := range habit.Entries {
		if entry.EntryDate.Time().After(lastEntryDate.Time()) {
			lastEntryDate = &entry.EntryDate
		}
	}

	nextOccurrence := lastEntryDate.Time().Add(time.Duration(*habit.Duration) * time.Millisecond)
	if nextOccurrence.After(now) {
		log.Debug().Msgf("Next occurrence is in the future: %s", nextOccurrence.Format(time.RFC3339))
		return false, time.Time{}
	}

	// Check if current time matches any reminder
	return findMatchingReminder(habit.Reminders, now)
}

// findMatchingReminder checks if the current time matches any of the habit's reminders
func findMatchingReminder(reminders []string, now time.Time) (bool, time.Time) {
	for _, reminder := range reminders {
		log.Debug().Msgf("Checking reminder: %s", reminder)
		reminderTime, err := time.Parse("15:04", reminder)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to parse reminder time: %s", reminder)
			continue
		}
		if now.Hour() == reminderTime.Hour() && now.Minute() == reminderTime.Minute() {
			log.Debug().Msgf("Reminder time matched: %s", reminderTime.Format(time.RFC3339))
			return true, reminderTime
		}
	}
	log.Debug().Msg("No matching reminder found")
	return false, time.Time{}
}
