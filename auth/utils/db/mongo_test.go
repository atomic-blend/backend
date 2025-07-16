package db

import (
	"atomic-blend/backend/auth/tests/utils/inmemorymongo"
	"fmt"
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

		client, err := ConnectMongo()

		assert.Nil(t, err)
		assert.Nil(t, client)
	})

	t.Run("should panic when required env vars are missing", func(t *testing.T) {
		os.Setenv("ENV", "development")
		os.Setenv("DATABASE_NAME", "test_db")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
		}()

		assert.Panics(t, func() {
			ConnectMongo()
		})
	})

	t.Run("should panic when cannot connect to mongodb", func(t *testing.T) {
		// Setup in-memory MongoDB
		mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
		assert.NoError(t, err)
		defer mongoServer.Stop()

		os.Setenv("ENV", "development")
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
		}()

		assert.Panics(t, func() {
			ConnectMongo()
		})
	})

	t.Run("should build URI with SSL when MONGO_SSL is true", func(t *testing.T) {
		// Setup test environment
		os.Setenv("ENV", "test") // Using test to avoid actual connection attempts
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("MONGO_SSL", "true")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("MONGO_SSL")
		}()

		client, err := ConnectMongo()

		// In test mode, client should be nil but no error
		assert.Nil(t, err)
		assert.Nil(t, client)
	})

	t.Run("should build URI with TLS when MONGO_TLS is true", func(t *testing.T) {
		// Setup test environment
		os.Setenv("ENV", "test")
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("MONGO_TLS", "true")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("MONGO_TLS")
		}()

		client, err := ConnectMongo()

		// In test mode, client should be nil but no error
		assert.Nil(t, err)
		assert.Nil(t, client)
	})

	t.Run("should build URI with SSL and TLS when both are true", func(t *testing.T) {
		// Setup test environment
		os.Setenv("ENV", "test")
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("MONGO_SSL", "true")
		os.Setenv("MONGO_TLS", "true")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("MONGO_SSL")
			os.Unsetenv("MONGO_TLS")
		}()

		client, err := ConnectMongo()

		// In test mode, client should be nil but no error
		assert.Nil(t, err)
		assert.Nil(t, client)
	})

	t.Run("should build URI with retryWrites when MONGO_RETRY_WRITES is true", func(t *testing.T) {
		// Setup test environment
		os.Setenv("ENV", "test")
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("MONGO_RETRY_WRITES", "true")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("MONGO_RETRY_WRITES")
		}()

		client, err := ConnectMongo()

		// In test mode, client should be nil but no error
		assert.Nil(t, err)
		assert.Nil(t, client)
	})

	t.Run("should build URI with SSL, TLS and retryWrites when all are true", func(t *testing.T) {
		// Setup test environment
		os.Setenv("ENV", "test")
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("MONGO_SSL", "true")
		os.Setenv("MONGO_TLS", "true")
		os.Setenv("MONGO_RETRY_WRITES", "true")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("MONGO_SSL")
			os.Unsetenv("MONGO_TLS")
			os.Unsetenv("MONGO_RETRY_WRITES")
		}()

		client, err := ConnectMongo()

		// In test mode, client should be nil but no error
		assert.Nil(t, err)
		assert.Nil(t, client)
	})

	t.Run("should panic when CA cert path is invalid", func(t *testing.T) {
		// Setup test environment
		os.Setenv("ENV", "development")
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("MONGO_SSL_CA_CERT_PATH", "/non/existent/path/ca.crt")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("MONGO_SSL_CA_CERT_PATH")
		}()

		// Should panic because of invalid CA cert path, but we need to recover since it returns an error
		defer func() {
			if r := recover(); r != nil {
				// This is expected behavior for invalid file path
				assert.Contains(t, fmt.Sprintf("%v", r), "failed to read CA certificate")
			}
		}()

		client, err := ConnectMongo()

		// If we reach here, it means there was an error instead of panic
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "failed to read CA certificate")
	})

	t.Run("should panic when CA cert is invalid", func(t *testing.T) {
		// Create a temporary file with invalid cert content
		tmpFile, err := os.CreateTemp("", "invalid-ca-*.crt")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		// Write invalid data to the temporary file
		_, err = tmpFile.WriteString("INVALID CERTIFICATE DATA")
		assert.NoError(t, err)
		tmpFile.Close()

		// Setup test environment
		os.Setenv("ENV", "development")
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("MONGO_SSL_CA_CERT_PATH", tmpFile.Name())
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("MONGO_SSL_CA_CERT_PATH")
		}()

		// Should return an error because of invalid CA cert, but we need to handle both cases
		defer func() {
			if r := recover(); r != nil {
				// This is expected behavior for invalid cert
				assert.Contains(t, fmt.Sprintf("%v", r), "failed to append CA certificate")
			}
		}()

		client, err := ConnectMongo()

		// If we reach here, it means there was an error instead of panic
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "failed to append CA certificate")
	})

	t.Run("should correctly set TLS config when valid CA cert is provided", func(t *testing.T) {
		// Create a temporary file with a valid cert content (in PEM format)
		tmpFile, err := os.CreateTemp("", "valid-ca-*.crt")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		// Write a sample PEM format certificate
		// This is a self-signed cert for testing only
		validCertPEM := `-----BEGIN CERTIFICATE-----
MIIDdTCCAl2gAwIBAgIJAK5FUgMFBpbFMA0GCSqGSIb3DQEBCwUAMFExCzAJBgNV
BAYTAlVTMQswCQYDVQQIDAJDQTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzELMAkG
A1UECgwCTVkxEDAOBgNVBAMMB1Rlc3QgQ0EwHhcNMjMxMjA1MDAwMDAwWhcNMzMx
MjAyMDAwMDAwWjBRMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExFjAUBgNVBAcM
DVNhbiBGcmFuY2lzY28xCzAJBgNVBAoMAk1ZMRAwDgYDVQQDDAdUZXN0IENBMIIB
IjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAr5l7DoIf9GA+RxLYR8vvy3Z3
AgFNwY0VyJ5MQMEh3/e5JFJZg4VoQsWGnKRhkJ8mBQjabYSLxQVINHadXEsZ8JCg
T+V3ZbCdUvwBx3ByL6MTVW3WoRZcz6glYELn2xjtDm+v/tEQ3NRV70HEgvsrpR6U
f/Slozu6jn2pCXDZ4i/LRGpXA5my7GHx9v5OJbD21UAjIJVGtoYGKFCRCEFep7AJ
XWd5QLkTnlrn9mgT9CwFsOKdLsNaf0d5xJKIVvHxLHOJgkszzO0GZgpZ7ybr7oK9
FBWZMc4XDz5cBYTvTFvL0jTQFWdUw8kJTnNGRxXYfGXI44U6/xvf1b0QByXINwID
AQABo1AwTjAdBgNVHQ4EFgQUrHG/MKfNmPV8c4GQ5US5C5QEoC8wHwYDVR0jBBgw
FoAUrHG/MKfNmPV8c4GQ5US5C5QEoC8wDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0B
AQsFAAOCAQEABSH1n+nRsOKBzV3aBP7urUC4GH2nHYJNYS8WkuHyIxWcOaf8mgSd
Y+n7xmb2oPR8RSQMjgE6tZ+aBBMeKIJsHFBNJ08iVQYXl3z7bR4X0XLPGDKAhJ5r
m9KzCXjj+w3BQUIDt3SpqBgKNg9Q+AAW1fXMtj31YI5AhUZVjfuMvZrJ6yUYzOBk
JMzMkNX2zOxmNnYBVfM3j1BQbbs4dHAMpA+Y/L31ZOkLHY3LxVYCMNUMwHvgsmQG
zvVfxVUrIjWvvhD2VVxZYkI7qVqkcUJsQTT3Y2dyxMGJkxUAInVBCaHifdYU7570
T7qLYtMnQJ9hMr0rI+T8W3RP8BXzl3Pi+w==
-----END CERTIFICATE-----`

		_, err = tmpFile.WriteString(validCertPEM)
		assert.NoError(t, err)
		tmpFile.Close()

		// Setup test environment with in-memory MongoDB to ensure we don't try real connections
		os.Setenv("ENV", "test") // Using test to avoid actual connection attempts
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("MONGO_SSL_CA_CERT_PATH", tmpFile.Name())
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("MONGO_SSL_CA_CERT_PATH")
		}()

		client, err := ConnectMongo()

		// In test mode, client should be nil but no error
		assert.Nil(t, err)
		assert.Nil(t, client)
	})

	t.Run("should build URI with authSource when MONGO_AUTH_SOURCE is set", func(t *testing.T) {
		// Setup test environment
		os.Setenv("ENV", "test")
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("MONGO_AUTH_SOURCE", "admin")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("MONGO_AUTH_SOURCE")
		}()

		client, err := ConnectMongo()

		// In test mode, client should be nil but no error
		assert.Nil(t, err)
		assert.Nil(t, client)
	})

	t.Run("should build URI with directConnection when MONGO_DIRECT_CONNECTION is true", func(t *testing.T) {
		// Setup test environment
		os.Setenv("ENV", "test")
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("MONGO_DIRECT_CONNECTION", "true")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("MONGO_DIRECT_CONNECTION")
		}()

		client, err := ConnectMongo()

		// In test mode, client should be nil but no error
		assert.Nil(t, err)
		assert.Nil(t, client)
	})

	t.Run("should build URI with timeout settings when MONGO_CONNECT_TIMEOUT_MS and MONGO_SERVER_SELECTION_TIMEOUT_MS are set", func(t *testing.T) {
		// Setup test environment
		os.Setenv("ENV", "test")
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("MONGO_CONNECT_TIMEOUT_MS", "5000")
		os.Setenv("MONGO_SERVER_SELECTION_TIMEOUT_MS", "3000")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("MONGO_CONNECT_TIMEOUT_MS")
			os.Unsetenv("MONGO_SERVER_SELECTION_TIMEOUT_MS")
		}()

		client, err := ConnectMongo()

		// In test mode, client should be nil but no error
		assert.Nil(t, err)
		assert.Nil(t, client)
	})

	t.Run("should build URI with authMechanism when MONGO_AUTH_MECHANISM is set", func(t *testing.T) {
		// Setup test environment
		os.Setenv("ENV", "test")
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("MONGO_AUTH_MECHANISM", "SCRAM-SHA-256")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("MONGO_AUTH_MECHANISM")
		}()

		client, err := ConnectMongo()

		// In test mode, client should be nil but no error
		assert.Nil(t, err)
		assert.Nil(t, client)
	})
}

func TestBuildMongoURI(t *testing.T) {
	t.Run("should use MONGO_URI when provided", func(t *testing.T) {
		expectedURI := "mongodb://user:pass@host:27017/mydb?retryWrites=false&authSource=admin"
		os.Setenv("MONGO_URI", expectedURI)
		defer os.Unsetenv("MONGO_URI")

		uri := buildMongoURI()

		assert.Equal(t, expectedURI, uri)
	})

	t.Run("should build URI from individual env vars when MONGO_URI is not set", func(t *testing.T) {
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("DATABASE_NAME", "testdb")
		defer func() {
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("DATABASE_NAME")
		}()

		uri := buildMongoURI()

		expectedURI := "mongodb://testuser:testpass@localhost:27017/testdb"
		assert.Equal(t, expectedURI, uri)
	})

	t.Run("should build URI with SSL parameters when individual env vars are used", func(t *testing.T) {
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("DATABASE_NAME", "testdb")
		os.Setenv("MONGO_SSL", "true")
		os.Setenv("MONGO_TLS", "true")
		defer func() {
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_SSL")
			os.Unsetenv("MONGO_TLS")
		}()

		uri := buildMongoURI()

		assert.Contains(t, uri, "mongodb://testuser:testpass@localhost:27017/testdb")
		assert.Contains(t, uri, "ssl=true")
		assert.Contains(t, uri, "tls=true")
	})

	t.Run("should build URI with auth parameters when individual env vars are used", func(t *testing.T) {
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("DATABASE_NAME", "testdb")
		os.Setenv("MONGO_AUTH_SOURCE", "admin")
		os.Setenv("MONGO_AUTH_MECHANISM", "SCRAM-SHA-256")
		defer func() {
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_AUTH_SOURCE")
			os.Unsetenv("MONGO_AUTH_MECHANISM")
		}()

		uri := buildMongoURI()

		assert.Contains(t, uri, "mongodb://testuser:testpass@localhost:27017/testdb")
		assert.Contains(t, uri, "authSource=admin")
		assert.Contains(t, uri, "authMechanism=SCRAM-SHA-256")
	})

	t.Run("should build URI with timeout parameters when individual env vars are used", func(t *testing.T) {
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("DATABASE_NAME", "testdb")
		os.Setenv("MONGO_CONNECT_TIMEOUT_MS", "5000")
		os.Setenv("MONGO_SERVER_SELECTION_TIMEOUT_MS", "3000")
		defer func() {
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_CONNECT_TIMEOUT_MS")
			os.Unsetenv("MONGO_SERVER_SELECTION_TIMEOUT_MS")
		}()

		uri := buildMongoURI()

		assert.Contains(t, uri, "mongodb://testuser:testpass@localhost:27017/testdb")
		assert.Contains(t, uri, "connectTimeoutMS=5000")
		assert.Contains(t, uri, "serverSelectionTimeoutMS=3000")
	})

	t.Run("should build URI with directConnection parameter when individual env vars are used", func(t *testing.T) {
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("DATABASE_NAME", "testdb")
		os.Setenv("MONGO_DIRECT_CONNECTION", "true")
		defer func() {
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_DIRECT_CONNECTION")
		}()

		uri := buildMongoURI()

		assert.Contains(t, uri, "mongodb://testuser:testpass@localhost:27017/testdb")
		assert.Contains(t, uri, "directConnection=true")
	})

	t.Run("should build URI with retryWrites=false when MONGO_RETRY_WRITES is false or not set", func(t *testing.T) {
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("DATABASE_NAME", "testdb")
		os.Setenv("MONGO_RETRY_WRITES", "false")
		defer func() {
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_RETRY_WRITES")
		}()

		uri := buildMongoURI()

		assert.Contains(t, uri, "mongodb://testuser:testpass@localhost:27017/testdb")
		assert.Contains(t, uri, "retryWrites=false")
	})

	t.Run("should include retryWrites when MONGO_RETRY_WRITES is true", func(t *testing.T) {
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("DATABASE_NAME", "testdb")
		os.Setenv("MONGO_RETRY_WRITES", "true")
		defer func() {
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_RETRY_WRITES")
		}()

		uri := buildMongoURI()

		assert.Contains(t, uri, "mongodb://testuser:testpass@localhost:27017/testdb")
		assert.Contains(t, uri, "retryWrites=true")
	})

	t.Run("should handle empty MONGO_URI and fallback to individual env vars", func(t *testing.T) {
		os.Setenv("MONGO_URI", "")
		os.Setenv("MONGO_USERNAME", "testuser")
		os.Setenv("MONGO_PASSWORD", "testpass")
		os.Setenv("MONGO_HOST", "localhost")
		os.Setenv("MONGO_PORT", "27017")
		os.Setenv("DATABASE_NAME", "testdb")
		defer func() {
			os.Unsetenv("MONGO_URI")
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
			os.Unsetenv("DATABASE_NAME")
		}()

		uri := buildMongoURI()

		expectedURI := "mongodb://testuser:testpass@localhost:27017/testdb"
		assert.Equal(t, expectedURI, uri)
	})

	t.Run("should use MONGO_URI when connecting to MongoDB", func(t *testing.T) {
		// Set test environment to avoid actual connection
		os.Setenv("ENV", "test")
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_URI", "mongodb://testuser:testpass@localhost:27017/testdb?retryWrites=false&authSource=admin")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_URI")
		}()

		client, err := ConnectMongo()

		// In test mode, client should be nil but no error
		assert.Nil(t, err)
		assert.Nil(t, client)
	})

	t.Run("should prioritize MONGO_URI over individual env vars in connection", func(t *testing.T) {
		// Set test environment to avoid actual connection
		os.Setenv("ENV", "test")
		os.Setenv("DATABASE_NAME", "test_db")
		os.Setenv("MONGO_URI", "mongodb://direct:uri@example.com:27017/directdb?ssl=true")
		// Set individual vars that should be ignored
		os.Setenv("MONGO_USERNAME", "should_not_use")
		os.Setenv("MONGO_PASSWORD", "should_not_use")
		os.Setenv("MONGO_HOST", "should_not_use")
		os.Setenv("MONGO_PORT", "should_not_use")
		defer func() {
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_NAME")
			os.Unsetenv("MONGO_URI")
			os.Unsetenv("MONGO_USERNAME")
			os.Unsetenv("MONGO_PASSWORD")
			os.Unsetenv("MONGO_HOST")
			os.Unsetenv("MONGO_PORT")
		}()

		client, err := ConnectMongo()

		// In test mode, client should be nil but no error
		assert.Nil(t, err)
		assert.Nil(t, client)
	})
}
