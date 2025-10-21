package cron

import (
	"bytes"
	"context"
	"html/template"
	"os"
	"strconv"

	waitinglist "github.com/atomic-blend/backend/auth/models/waiting_list"
	"github.com/atomic-blend/backend/auth/repositories"
	"github.com/atomic-blend/backend/shared/repositories/user"
	amqpservice "github.com/atomic-blend/backend/shared/services/amqp"
	"github.com/atomic-blend/backend/shared/utils/db"
	"github.com/atomic-blend/backend/shared/utils/password"
	"github.com/rs/zerolog/log"
	"slices"
)

func WaitingListCron() {
	log.Debug().Msg("Starting waiting list cron job")

	// get the domain from the environment variables
	domain := os.Getenv("PUBLIC_ADDRESS")
	if domain == "" {
		domain = "mail.atomic-blend.com"
	} else {
		domain = "mail." + domain
	}
	log.Debug().Msgf("Domain: %s", domain)

	// get the max number of users
	maxUsersStr := os.Getenv("AUTH_MAX_NB_USER")
	if maxUsersStr == "" {
		maxUsersStr = "10"
	}
	maxUsers, err := strconv.ParseInt(maxUsersStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert max number of users to int")
		return
	}

	userRepo := user.NewUserRepository(db.Database)

	// get the current number of users
	nbOfUsers, err := userRepo.Count(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current number of users")
		return
	}

	if nbOfUsers >= maxUsers {
		log.Debug().Msg("Max number of users reached, no more users can be added")
		return
	}

	remainingSpots := maxUsers - nbOfUsers
	if remainingSpots <= 0 {
		log.Debug().Msg("No remaining spots, no more users can be added")
		return
	}

	log.Debug().Msgf("Remaining spots: %d", remainingSpots)

	// get the waiting list
	waitingListRepo := repositories.NewWaitingListRepository(db.Database)
	waitingList, err := waitingListRepo.GetOldest(context.TODO(), remainingSpots)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get waiting list")
		return
	}

	// remove items that already have a code
	waitingList = slices.DeleteFunc(waitingList, func(item *waitinglist.WaitingList) bool {
		return item.Code != nil && *item.Code != ""
	})

	for _, waitingListItem := range waitingList {
		log.Debug().Msgf("Waiting list item: %s", waitingListItem.Email)

		// generate a code (32 char base64 encoded string)
		code, err := password.GenerateRandomString(32)
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate code")
			return
		}

		// update the waiting list item with the code
		waitingListItem.Code = &code
		_, err = waitingListRepo.Update(context.TODO(), waitingListItem.ID.Hex(), waitingListItem)
		if err != nil {
			log.Error().Err(err).Msg("Failed to update waiting list item")
			return
		}

		// send an email to the user with the code
		htmlTemplate, err := template.ParseFiles("./email_templates/waiting_list_success/waiting_list_success.html")
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse HTML template")
			return
		}

		textTemplate, err := template.ParseFiles("./email_templates/waiting_list_success/waiting_list_success.txt")
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse text template")
			return
		}

		// template the plain text with gotemplate
		var htmlContent bytes.Buffer
		err = htmlTemplate.Execute(&htmlContent, map[string]string{
			"email":         waitingListItem.Email,
			"domain":        domain,
			"securityToken": *waitingListItem.SecurityToken,
		})

		if err != nil {
			log.Error().Err(err).Msg("Failed to execute HTML template")
			return
		}

		var textContent bytes.Buffer
		err = textTemplate.Execute(&textContent, map[string]string{
			"email":         waitingListItem.Email,
			"domain":        domain,
			"securityToken": *waitingListItem.SecurityToken,
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to execute text template")
			return
		}

		// Create RawMail structure for AMQP
		rawMail := map[string]interface{}{
			"headers": map[string]interface{}{
				"To":      []string{waitingListItem.Email},
				"From":    "noreply@atomic-blend.com",
				"Subject": "You just joined the waiting list!",
			},
			"htmlContent":    htmlContent.String(),
			"textContent":    textContent.String(),
			"rejected":       false,
			"rewriteSubject": false,
			"graylisted":     false,
		}

		// Publish message to AMQP queue
		amqpService := amqpservice.NewAMQPService("AUTH")
		amqpService.PublishMessage("mail", "sent", map[string]interface{}{
			"waiting_list_email": true,
			"content":            rawMail,
		}, nil)

		log.Info().Msg("Waiting list email queued for sending")
	}

}
