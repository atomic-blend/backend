package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/atomic-blend/backend/auth/controllers/admin"
	"github.com/atomic-blend/backend/auth/controllers/auth"
	"github.com/atomic-blend/backend/auth/controllers/health"
	"github.com/atomic-blend/backend/auth/controllers/users"
	"github.com/atomic-blend/backend/auth/controllers/webhooks"
	"github.com/atomic-blend/backend/shared/models"
	"github.com/atomic-blend/backend/shared/utils/db"

	"github.com/gin-contrib/cors"

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

	// Initialize MongoDB connection
	client, err := db.ConnectMongo()
	if err != nil {
		log.Fatal().Err(err).Msg("❌ Error connecting to MongoDB")
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal().Err(err).Msg("❌ Error disconnecting from MongoDB")
		}
		log.Fatal().Msg("✅ Disconnected from MongoDB")
	}()

	// start grpc server
	go startGRPCServer()

	models.RegisterValidators()

	// Setup router with middleware
	router := gin.Default()

	// Configure CORS if environment variables are provided
	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	corsMethods := os.Getenv("CORS_ALLOWED_METHODS")
	corsHeaders := os.Getenv("CORS_ALLOWED_HEADERS")
	corsExposeHeaders := os.Getenv("CORS_EXPOSE_HEADERS")
	corsCredentials := os.Getenv("CORS_ALLOW_CREDENTIALS")
	corsMaxAge := os.Getenv("CORS_MAX_AGE")

	// Apply CORS only if we have sufficient configuration
	if corsOrigins != "" {
		var allowedOrigins []string

		// Build origins list
		if corsOrigins != "" {
			// Split by comma and trim spaces
			for _, origin := range strings.Split(corsOrigins, ",") {
				if trimmed := strings.TrimSpace(origin); trimmed != "" {
					allowedOrigins = append(allowedOrigins, trimmed)
				}
			}
		}

		// Default methods if not specified
		allowedMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
		if corsMethods != "" {
			allowedMethods = []string{}
			for _, method := range strings.Split(corsMethods, ",") {
				if trimmed := strings.TrimSpace(method); trimmed != "" {
					allowedMethods = append(allowedMethods, trimmed)
				}
			}
		}

		// Default headers if not specified
		allowedHeaders := []string{"Origin", "Content-Type", "Authorization"}
		if corsHeaders != "" {
			allowedHeaders = []string{}
			for _, header := range strings.Split(corsHeaders, ",") {
				if trimmed := strings.TrimSpace(header); trimmed != "" {
					allowedHeaders = append(allowedHeaders, trimmed)
				}
			}
		}

		// Default expose headers if not specified
		exposeHeaders := []string{"Content-Length"}
		if corsExposeHeaders != "" {
			exposeHeaders = []string{}
			for _, header := range strings.Split(corsExposeHeaders, ",") {
				if trimmed := strings.TrimSpace(header); trimmed != "" {
					exposeHeaders = append(exposeHeaders, trimmed)
				}
			}
		}

		// Default credentials to true if not specified
		allowCredentials := true
		if corsCredentials != "" {
			allowCredentials = strings.ToLower(corsCredentials) == "true"
		}

		// Default max age to 12 hours if not specified
		maxAge := 12 * time.Hour
		if corsMaxAge != "" {
			if duration, err := time.ParseDuration(corsMaxAge); err == nil {
				maxAge = duration
			}
		}

		router.Use(cors.New(cors.Config{
			AllowOrigins:     allowedOrigins,
			AllowMethods:     allowedMethods,
			AllowHeaders:     allowedHeaders,
			ExposeHeaders:    exposeHeaders,
			AllowCredentials: allowCredentials,
			MaxAge:           maxAge,
		}))

		log.Info().Strs("origins", allowedOrigins).Msg("CORS configured")
	} else {
		log.Info().Msg("No CORS configuration found, skipping CORS setup")
	}

	// Register all routes
	auth.SetupRoutes(router, db.Database)
	users.SetupRoutes(router, db.Database)
	admin.SetupRoutes(router, db.Database)
	health.SetupRoutes(router, db.Database)
	webhooks.SetupRoutes(router, db.Database)

	// Define port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Msgf("Server starting on port %s", port)
	router.Run(":" + port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
