package main

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/controllers/admin"
	"atomic_blend_api/controllers/health"
	"atomic_blend_api/controllers/tasks"
	"atomic_blend_api/controllers/users"
	"atomic_blend_api/cron"
	"atomic_blend_api/utils/db"
	"context"
	"os"

	"github.com/jasonlvhit/gocron"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Router is exported for use in other packages
// nolint
var Router *gin.Engine

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	_ = godotenv.Load()

	mongoURI := os.Getenv("MONGO_URI")

	mongoUsername := os.Getenv("MONGO_USERNAME")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")
	databaseName := os.Getenv("DATABASE_NAME")
	ssl := os.Getenv("MONGO_SSL")
	tls := os.Getenv("MONGO_TLS")
	retryWrites := os.Getenv("MONGO_RETRY_WRITES")

	if mongoUsername != "" && mongoPassword != "" && mongoHost != "" {
		mongoURI = "mongodb://" + mongoUsername + ":" + mongoPassword + "@" + mongoHost
		if mongoPort != "" {
			mongoURI = mongoURI + ":" + mongoPort
		}
		if databaseName != "" {
			mongoURI = mongoURI + "/" + databaseName
		}
	} else if mongoURI == "" {
		mongoURI = "mongodb://mongo_user:password@mongodb:27017"
	}

	if ssl == "true" {
		log.Debug().Msg("Setting SSL to true")
		mongoURI += "?ssl=true"
	}
	if tls == "true" {
		if ssl != "true" {
			log.Debug().Msg("Setting TLS to true")
			mongoURI += "?tls=true"
		} else {
			log.Debug().Msg("Setting TLS to true with SSL")
			mongoURI += "&tls=true"
		}
	}
	if retryWrites == "true" {
		if ssl != "true" && tls != "true" {
			log.Debug().Msg("Setting retryWrites to true")
			mongoURI += "?retryWrites=true"
		} else {
			log.Debug().Msg("Setting retryWrites to true with SSL/TLS")
			mongoURI += "&retryWrites=true"
		}
	}

	log.Debug().Msgf("MONGO_USERNAME: %s", mongoUsername)
	log.Debug().Msgf("MONGO_PASSWORD: %s", mongoPassword)
	log.Debug().Msgf("MONGO_HOST: %s", mongoHost)
	log.Debug().Msgf("MONGO_PORT: %s", mongoPort)
	log.Debug().Msgf("DATABASE_NAME: %s", databaseName)
	log.Debug().Msgf("MONGO_SSL: %s", ssl)
	log.Debug().Msgf("MONGO_TLS: %s", tls)
	log.Debug().Msgf("MONGO_RETRY_WRITES: %s", retryWrites)
	log.Debug().Msgf("ENV: %s", os.Getenv("ENV"))

	// Initialize MongoDB connection
	client, err := db.ConnectMongo(&mongoURI)
	if err != nil {
		log.Fatal().Err(err).Msg("❌ Error connecting to MongoDB")
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal().Err(err).Msg("❌ Error disconnecting from MongoDB")
		}
		log.Fatal().Msg("✅ Disconnected from MongoDB")
	}()

	// start cron
	go func() {
		err := gocron.Every(60).Seconds().Do(cron.MainCron)
		if err != nil {
			log.Error().Err(err).Msg("Error defining cron job")
		}
		<-gocron.Start()
	}()

	// Setup router with middleware
	router := gin.Default()

	// Register all routes
	auth.SetupRoutes(router, db.Database)
	users.SetupRoutes(router, db.Database)
	admin.SetupRoutes(router, db.Database)
	tasks.SetupRoutes(router, db.Database)
	health.SetupRoutes(router, db.Database)

	// Define port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Msgf("Server starting on port %s", port)
	router.Run(":" + port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
