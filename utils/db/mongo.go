package db

import (
	"atomic_blend_api/utils/shortcuts"
	"context"
	"fmt"
	"os"
	"time"

	tls "crypto/tls"
	"crypto/x509"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient is the MongoDB client
var MongoClient *mongo.Client

// Database is the MongoDB database instance
var Database *mongo.Database

// ConnectMongo initializes and returns a MongoDB client
func ConnectMongo(uri *string) (*mongo.Client, error) {
	env := os.Getenv("ENV")
	databaseName := os.Getenv("DATABASE_NAME")
	ssl := os.Getenv("MONGO_SSL")
	tlsConfig := os.Getenv("MONGO_TLS")
	sslCaCert := os.Getenv("MONGO_SSL_CA_CERT_PATH")
	retryWrites := os.Getenv("MONGO_RETRY_WRITES")

	// optionally set ssl and tls and retryWrites if they are set to true
	if ssl == "true" {
		log.Debug().Msg("Setting SSL to true")
		*uri += "?ssl=true"
	}
	if tlsConfig == "true" {
		if ssl != "true" {
			log.Debug().Msg("Setting TLS to true")
			*uri += "?tls=true"
		} else {
			log.Debug().Msg("Setting TLS to true with SSL")
			*uri += "&tls=true"
		}
	}
	if retryWrites == "true" {
		log.Debug().Msg("Setting retryWrites to param")
		*uri += "?retryWrites="+retryWrites
	}
	

	if env == "test" {
		// setup in memory mongo for testing
		log.Debug().Msg("Setting up in memory mongo for testing")
		return nil, nil
	}
	if uri == nil {
		err := fmt.Errorf("MONGO_URI is not set")
		return nil, err
	}

	shortcuts.CheckRequiredEnvVar("DATABASE_NAME", databaseName, "")

	clientOptions := options.Client().ApplyURI(*uri)
	if sslCaCert != "" {
		// Load CA cert if provided
		caCert, err := os.ReadFile(sslCaCert)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}
		
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to append CA certificate")
		}
		
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    certPool,
		}
		clientOptions.SetTLSConfig(tlsConfig)
	}

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
