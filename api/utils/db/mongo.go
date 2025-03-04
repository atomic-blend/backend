package db

import (
	"atomic_blend_api/utils/shortcuts"
	"context"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var Database *mongo.Database

var (
	databaseName = os.Getenv("DATABASE_NAME")
)


// ConnectMongo initializes and returns a MongoDB client
func ConnectMongo(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)

	shortcuts.CheckRequiredEnvVar("DATABASE_NAME", databaseName, "")


	// Set timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Debug().Msg("Connecting to MongoDB")
	client, err := mongo.Connect(ctx, clientOptions)
	shortcuts.FailOnError(err, "Error connecting to MongoDB")

	// Ping the database to confirm connection
	log.Debug().Msg("Pinging MongoDB")
	err = client.Ping(ctx, nil)
	shortcuts.FailOnError(err, "Error pinging MongoDB")

	log.Debug().Msg("✅ Successfully connected to MongoDB")
	MongoClient = client
	Database = client.Database(databaseName)
	return client, nil
}
