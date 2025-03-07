package shortcuts

import (
	"github.com/rs/zerolog/log"
)

func CheckRequiredEnvVar(varName string, configValue string, defaultValue string) {
	if configValue == "" {
		if defaultValue != "" {
			configValue = defaultValue
		} else {
			panic(varName + " is required")
			log.Fatal().Msgf("%s is required", varName)
		}
	}
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Error().Err(err).Msg(msg)
		panic(err)
	}
}

func LogOnError(err error, msg string) bool {
	if err != nil {
		log.Error().Err(err).Msg(msg)
		return true
	}
	return false
}
