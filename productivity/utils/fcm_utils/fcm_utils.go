package fcmutils

import (
	"context"

	"firebase.google.com/go/v4/messaging"
	"github.com/appleboy/go-fcm"
	"github.com/rs/zerolog/log"
)

// SendMulticast sends a multicast message to multiple device tokens.
func SendMulticast(ctx context.Context, client *fcm.Client, data map[string]string, deviceTokens []string) {
	res, err := client.SendMulticast(
		ctx,
		&messaging.MulticastMessage{
			Data:   data,
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
					"apns-priority":  "5", // Must be `5` when `contentAvailable` is set to true.
				},
			},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send message to FCM")
	}
	log.Debug().Msgf("Sent message to %d devices, %d errors", res.SuccessCount, res.FailureCount)
	for i, respo := range res.Responses {
		if respo.Error != nil {
			log.Error().Err(respo.Error).Msgf("Failed to send message to device %s", deviceTokens[i])
		} else {
			log.Debug().Msgf("Successfully sent message to device %s", deviceTokens[i])
		}
	}
}
