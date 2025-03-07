package in_memory_mongo

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateInMemoryMongoDB(t *testing.T) {
	t.Run("Default Version", func(t *testing.T) {
		// Clear any existing version
		os.Unsetenv("MONGO_VERSION")

		server, err := CreateInMemoryMongoDB()
		require.NoError(t, err)
		assert.NotNil(t, server)

		defer server.Stop()

		// Verify the server is running by getting its URI
		uri := server.URI()
		assert.NotEmpty(t, uri)
	})

	t.Run("Custom Version", func(t *testing.T) {
		// Set custom version
		customVersion := "6.0.0"
		os.Setenv("MONGO_VERSION", customVersion)
		defer os.Unsetenv("MONGO_VERSION")

		server, err := CreateInMemoryMongoDB()
		require.NoError(t, err)
		assert.NotNil(t, server)

		defer server.Stop()

		uri := server.URI()
		assert.NotEmpty(t, uri)
	})
}

func TestConnectToInMemoryDB(t *testing.T) {
	t.Run("Successful Connection", func(t *testing.T) {
		// First create a server
		server, err := CreateInMemoryMongoDB()
		require.NoError(t, err)
		defer server.Stop()

		// Try to connect
		client, err := ConnectToInMemoryDB(server.URI())
		require.NoError(t, err)
		assert.NotNil(t, client)

		// Verify connection works by pinging
		err = client.Ping(context.Background(), nil)
		assert.NoError(t, err)

		client.Disconnect(context.Background())
	})

	t.Run("Invalid URI", func(t *testing.T) {
		client, err := ConnectToInMemoryDB("invalid-uri")
		assert.Error(t, err)
		assert.Nil(t, client)
	})
}
