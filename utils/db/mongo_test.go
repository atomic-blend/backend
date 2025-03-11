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
}
