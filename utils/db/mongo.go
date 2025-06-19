package db

import (
	"atomic_blend_api/utils/shortcuts"
	"context"
	"fmt"
	"os"
	"strings"
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

func buildMongoURI() string {
	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")
	database := os.Getenv("DATABASE_NAME")

	// Récupération des paramètres booléens
	ssl := os.Getenv("MONGO_SSL") == "true"
	tls := os.Getenv("MONGO_TLS") == "true"
	retryWrites := os.Getenv("MONGO_RETRY_WRITES") == "true"
	directConnection := os.Getenv("MONGO_DIRECT_CONNECTION") == "true"

	authSource := os.Getenv("MONGO_AUTH_SOURCE")
	
	// Construction de l'URI avec les paramètres
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
		username, password, host, port, database)

	// Ajout des paramètres de requête
	params := []string{}

	if !retryWrites {
		params = append(params, "retryWrites=false")
	}

	if authSource != "" {
		params = append(params, fmt.Sprintf("authSource=%s", authSource))
	}

	if directConnection {
		params = append(params, "directConnection=true")
	} else {
		params = append(params, "directConnection=false")
	}

	if ssl {
		params = append(params, "ssl=true")
	}

	if tls {
		params = append(params, "tls=true")
	}

	// Ajout des timeouts si définis
	if connectTimeout := os.Getenv("MONGO_CONNECT_TIMEOUT_MS"); connectTimeout != "" {
		params = append(params, fmt.Sprintf("connectTimeoutMS=%s", connectTimeout))
	}

	if serverSelectionTimeout := os.Getenv("MONGO_SERVER_SELECTION_TIMEOUT_MS"); serverSelectionTimeout != "" {
		params = append(params, fmt.Sprintf("serverSelectionTimeoutMS=%s", serverSelectionTimeout))
	}

	if len(params) > 0 {
		uri += "?" + strings.Join(params, "&")
	}

	print(uri)

	return uri
}

// ConnectMongo initializes and returns a MongoDB client
func ConnectMongo() (*mongo.Client, error) {
	env := os.Getenv("ENV")
	sslCaCert := os.Getenv("MONGO_SSL_CA_CERT_PATH")

	databaseName := os.Getenv("DATABASE_NAME")
	shortcuts.CheckRequiredEnvVar("DATABASE_NAME", databaseName, "")

	uri := buildMongoURI()

	if env == "test" {
		// setup in memory mongo for testing
		log.Debug().Msg("Setting up in memory mongo for testing")
		return nil, nil
	}

	clientOptions := options.Client().ApplyURI(uri)
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
