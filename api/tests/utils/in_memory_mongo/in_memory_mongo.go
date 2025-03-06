package in_memory_mongo

import (
	"context"
	"os"
	"time"

	logger "log"

	"github.com/rs/zerolog/log"
	"github.com/atomic-blend/memongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateInMemoryMongoDB creates and starts an in-memory MongoDB instance
func CreateInMemoryMongoDB() (*memongo.Server, error) {
	// Check MongoDB version or set default
	mongoVersion := os.Getenv("MONGO_VERSION")
	if mongoVersion == "" {
		mongoVersion = "8.0.4" // Default version if not specified
	}

	// Create an in-memory MongoDB instance
	options := &memongo.Options{
		MongoVersion:   mongoVersion,
		StartupTimeout: 15 * time.Second,
		Logger:         logger.New(os.Stdout, "memongo: ", logger.LstdFlags),
		// LogLevel:       memongolog.LogLevelDebug,
	}
	mongoServer, err := memongo.StartWithOptions(options)
	if err != nil {
		log.Error().Err(err).Msg("Error starting in-memory MongoDB")
		return nil, err
	}

	return mongoServer, nil
}

// ConnectToInMemoryDB connects to the in-memory MongoDB instance
func ConnectToInMemoryDB(mongoURI string) (*mongo.Client, error) {
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to in-memory MongoDB")
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error pinging in-memory MongoDB")
		return nil, err
	}

	return client, nil
}