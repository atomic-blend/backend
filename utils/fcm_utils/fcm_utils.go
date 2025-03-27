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
			Tokens: deviceTokens,
			Android: &messaging.AndroidConfig{
				Priority: "high",
			},
			APNS: &messaging.APNSConfig{
				Payload: &messaging.APNSPayload{
					Aps: &messaging.Aps{
						ContentAvailable: true, // 🔥 This is the magic field
					},
				},
				Headers: map[string]string{
					"apns-push-type": "background",
					"apns-priority": "5", // Must be `5` when `contentAvailable` is set to true.
				},
			},
		},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to send message to FCM")
	}
}
