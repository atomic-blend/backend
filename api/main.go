package main

import (
	"atomic_blend_api/utils/db"
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Router *gin.Engine

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://mongo_user:password@mongodb:27017"
	}
	log.Info().Msgf("Connecting to MongoDB at %s", mongoURI)
	// Initialize MongoDB connection
	client, err := db.ConnectMongo(mongoURI)
	if err != nil {
		log.Fatal().Err(err).Msg("❌ Error connecting to MongoDB")
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal().Err(err).Msg("❌ Error disconnecting from MongoDB")
		}
		log.Fatal().Msg("✅ Disconnected from MongoDB")
	}()

	Router := gin.Default()

	Router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
