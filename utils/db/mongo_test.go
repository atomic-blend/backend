package db

import (
	"atomic_blend_api/tests/utils/inmemorymongo"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectMongo(t *testing.T) {
	t.Run("should return nil when in test environment", func(t *testing.T) {
		// Set test environment
		os.Setenv("ENV", "test")
		os.Setenv("DATABASE_NAME", "test_db")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
		}()

		uri := "mongodb://localhost:27017"
		client, err := ConnectMongo(&uri)

		assert.Nil(t, err)
		assert.Nil(t, client)
	})

	t.Run("should return error when uri is nil", func(t *testing.T) {
		os.Setenv("ENV", "development")
		os.Setenv("DATABASE_NAME", "test_db")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
		}()

		client, err := ConnectMongo(nil)

		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Equal(t, "MONGO_URI is not set", err.Error())
	})

	t.Run("should successfully connect to mongodb", func(t *testing.T) {
		// Setup in-memory MongoDB
		mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
		assert.NoError(t, err)
		defer mongoServer.Stop()

		os.Setenv("ENV", "development")
		os.Setenv("DATABASE_NAME", "test_db")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
		}()

		uri := mongoServer.URIWithRandomDB()
		client, err := ConnectMongo(&uri)

		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.NotNil(t, Database)
		assert.Equal(t, "test_db", Database.Name())

		// Cleanup
		client.Disconnect(context.TODO())
	})

	t.Run("should add ssl=true to uri when MONGO_SSL is true", func(t *testing.T) {
		// Setup test environment
		os.Setenv("ENV", "test") // Using test to avoid actual connection attempts
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_SSL", "true")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_SSL")
		}()

		uri := "mongodb://localhost:27017"
		originalURI := uri
		_, _ = ConnectMongo(&uri)

		// Test that SSL parameter was added
		assert.NotEqual(t, originalURI, uri)
		assert.Contains(t, uri, "?ssl=true")
	})

	t.Run("should add tls=true to uri when MONGO_TLS is true", func(t *testing.T) {
		// Setup test environment
		os.Setenv("ENV", "test")
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_TLS", "true")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_TLS")
		}()

		uri := "mongodb://localhost:27017"
		originalURI := uri
		_, _ = ConnectMongo(&uri)

		// Test that TLS parameter was added
		assert.NotEqual(t, originalURI, uri)
		assert.Contains(t, uri, "?tls=true")
	})

	t.Run("should add ssl=true and tls=true to uri when both are true", func(t *testing.T) {
		// Setup test environment
		os.Setenv("ENV", "test")
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_SSL", "true")
		os.Setenv("MONGO_TLS", "true")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_SSL")
			os.Unsetenv("MONGO_TLS")
		}()

		uri := "mongodb://localhost:27017"
		originalURI := uri
		_, _ = ConnectMongo(&uri)

		// Test that both parameters were added
		assert.NotEqual(t, originalURI, uri)
		assert.Contains(t, uri, "ssl=true")
		assert.Contains(t, uri, "tls=true")
	})

	t.Run("should add retryWrites=true to uri when MONGO_RETRY_WRITES is true", func(t *testing.T) {
		// Setup test environment
		os.Setenv("ENV", "test")
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_RETRY_WRITES", "true")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_RETRY_WRITES")
		}()

		uri := "mongodb://localhost:27017"
		originalURI := uri
		_, _ = ConnectMongo(&uri)

		// Test that retryWrites parameter was added
		assert.NotEqual(t, originalURI, uri)
		assert.Contains(t, uri, "?retryWrites=true")
	})

	t.Run("should add ssl, tls and retryWrites to uri when all are true", func(t *testing.T) {
		// Setup test environment
		os.Setenv("ENV", "test")
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_SSL", "true")
		os.Setenv("MONGO_TLS", "true")
		os.Setenv("MONGO_RETRY_WRITES", "true")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_SSL")
			os.Unsetenv("MONGO_TLS")
			os.Unsetenv("MONGO_RETRY_WRITES")
		}()

		uri := "mongodb://localhost:27017"
		originalURI := uri
		_, _ = ConnectMongo(&uri)

		// Test that all parameters were added
		assert.NotEqual(t, originalURI, uri)
		assert.Contains(t, uri, "ssl=true")
		assert.Contains(t, uri, "tls=true")
		assert.Contains(t, uri, "retryWrites=true")
	})
}
