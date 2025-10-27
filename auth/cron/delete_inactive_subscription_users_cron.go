package cron

import (
	"context"
	"errors"

	userdeleter "github.com/atomic-blend/backend/auth/utils/user_deleter"
	productivityclient "github.com/atomic-blend/backend/shared/grpc/productivity"
	"github.com/atomic-blend/backend/shared/repositories/user"
	"github.com/atomic-blend/backend/shared/utils/db"
	"github.com/rs/zerolog/log"
)

func DeleteInactiveSubscriptionUsersCron() {
	log.Info().Msg("Running DeleteInactiveSubscriptionUsersCron")
	// TODO: get the user with subscriptionId == nil || status == cancelled
	// for more than X days (from env, since cancellation date if cancelled, or createdAt date if no subscriptionId)

	// Initialize database and repositories
	userRepo, productivityclientClient, err := initialize()
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize repositories")
		return
	}

	DeleteInactiveSubscriptionUsers(userRepo, productivityclientClient)

	log.Info().Msg("Completed DeleteInactiveSubscriptionUsersCron")
}

func DeleteInactiveSubscriptionUsers(userRepo user.Interface, productivityClient productivityclient.Interface) {
	// get users with inactive subscriptions that exceed the grace period
	gracePeriodDays := 7
	users, err := userRepo.FindInactiveSubscriptionUsers(context.TODO(), gracePeriodDays)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find inactive subscription users")
		return
	}

	for _, u := range users {
		// delete user and data
		log.Info().Str("userID", u.ID.Hex()).Msg("delete user with inactive subscription")
		err := userdeleter.DeletePersonalDataAndUser(*u.ID, productivityClient, userRepo)
		if err != nil {
			log.Error().Err(err).Str("userID", u.ID.Hex()).Msg("Failed to delete personal data and user")
			continue
		}
		log.Info().Str("userID", u.ID.Hex()).Msg("Successfully deleted user with inactive subscription")
	}
}

func initialize() (*user.Repository, *productivityclient.ProductivityClient, error) {
	// Initialize database and repositories
	userRepo := user.NewUserRepository(db.Database)
	if userRepo == nil {
		log.Error().Msg("Failed to create user repository")
		return nil, nil, errors.New("failed to create user repository")
	}

	productivityclientClient, err := productivityclient.NewProductivityClient()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create productivity gRPC client")
		return nil, nil, err
	}

	return userRepo, productivityclientClient, nil
}
