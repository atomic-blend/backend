package fcmutils

import (
	"context"

	"firebase.google.com/go/v4/messaging"
	"github.com/appleboy/go-fcm"
	"github.com/rs/zerolog/log"
)

func SendMulticast(client *fcm.Client, ctx context.Context, data map[string]string, deviceTokens []string) {
	_, err := client.SendMulticast(
		ctx,
		&messaging.MulticastMessage{
			Data: data,
			Notification: &messaging.Notification{
				Title: data["title"],
			},
			Tokens: deviceTokens,
		},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to send message to FCM")
	}
}
