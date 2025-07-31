package mail

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

func RouteMessage(routingKey string, body []byte) {
	switch routingKey {
	case "received":
		// body is usually an map[string]interface{}
		// we need to parse it to a map[string]interface{}
		var data map[string]interface{}
		err := json.Unmarshal(body, &data)
		if err != nil {
			log.Error().Err(err).Msg("Error unmarshalling body")
			return
		}
		mimeContent, ok := data["content"].(string)
		if !ok {
			log.Error().Msg("Content is not a string")
			return
		}

		receiveMail(mimeContent)
	case "sent":
		//routeSentMessage()
	}
}
