package shortcuts

import (
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CheckRequiredEnvVar checks if a required environment variable is set
// If the variable is not set, it will log a fatal error
// @param varName string
// @param configValue string
// @param defaultValue string
// @return void
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

// FailOnError panics if err is not nil
// @param err error
// @param msg string
// @return void
func FailOnError(err error, msg string) {
	if err != nil {
		log.Error().Err(err).Msg(msg)
		panic(err)
	}
}

// LogOnError logs an error message if err is not nil
// Returns true if err is not nil
// @param err error
// @param msg string
// @return bool
func LogOnError(err error, msg string) bool {
	if err != nil {
		log.Error().Err(err).Msg(msg)
		return true
	}
	return false
}

func ContainsDateTime(slice []primitive.DateTime, item time.Time) bool {
	for _, a := range slice {
		if a.Time().Equal(item) {
			return true
		}
	}
	return false
}

func ContainsInt(slice []int, item int) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
