package main

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/controllers/admin"
	"atomic_blend_api/controllers/health"
	"atomic_blend_api/controllers/tasks"
	"atomic_blend_api/controllers/users"
	"atomic_blend_api/utils/db"
	"context"
	"os"

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
	if mongoURI == "" {
		mongoURI = "mongodb://mongo_user:password@mongodb:27017"
	}

	log.Info().Msgf("Connecting to MongoDB at %s", mongoURI)

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

	// Get database instance

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
